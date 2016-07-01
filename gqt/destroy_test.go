package gqt_test

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/cloudfoundry-incubator/garden"
	"github.com/cloudfoundry-incubator/guardian/gqt/runner"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("Destroying a Container", func() {
	var (
		client *runner.RunningGarden
	)

	BeforeEach(func() {
		client = startGarden()
	})

	AfterEach(func() {
		Expect(client.DestroyAndStop()).To(Succeed())
	})

	It("should not leak goroutines", func() {
		handle := fmt.Sprintf("goroutine-leak-test-%d", GinkgoParallelNode())

		numGoroutinesBefore, err := client.NumGoroutines()
		Expect(err).NotTo(HaveOccurred())

		_, err = client.Create(garden.ContainerSpec{
			Handle: handle,
		})
		Expect(err).NotTo(HaveOccurred())

		client.Destroy(handle)

		Eventually(func() int {
			numGoroutinesAfter, err := client.NumGoroutines()
			Expect(err).NotTo(HaveOccurred())
			return numGoroutinesAfter
		}).Should(Equal(numGoroutinesBefore))
	})

	It("should destroy the container's rootfs", func() {
		container, err := client.Create(garden.ContainerSpec{})
		Expect(err).NotTo(HaveOccurred())

		info, err := container.Info()
		Expect(err).NotTo(HaveOccurred())
		containerRootfs := info.ContainerPath

		Expect(client.Destroy(container.Handle())).To(Succeed())

		Expect(containerRootfs).NotTo(BeAnExistingFile())
	})

	It("should destroy the container's depot directory", func() {
		container, err := client.Create(garden.ContainerSpec{})
		Expect(err).NotTo(HaveOccurred())

		Expect(client.Destroy(container.Handle())).To(Succeed())

		Expect(filepath.Join(client.DepotDir, container.Handle())).NotTo(BeAnExistingFile())
	})

	It("should kill the container's init process", func() {
		container, err := client.Create(garden.ContainerSpec{})
		Expect(err).NotTo(HaveOccurred())

		initProcPid := initProcessPID(container.Handle())

		_, err = container.Run(garden.ProcessSpec{
			Path: "/bin/sh",
			Args: []string{
				"-c", "read x",
			},
		}, ginkgoIO)
		Expect(err).NotTo(HaveOccurred())

		Expect(client.Destroy(container.Handle())).To(Succeed())

		var killExitCode = func() int {
			sess, err := gexec.Start(exec.Command("kill", "-0", fmt.Sprintf("%d", initProcPid)), GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())
			sess.Wait(1 * time.Second)
			return sess.ExitCode()
		}

		Eventually(killExitCode).Should(Equal(1))
	})

	Describe("networking resources", func() {
		var (
			container         garden.Container
			networkSpec       string
			contIfaceName     string
			networkBridgeName string
		)

		BeforeEach(func() {
			var err error

			networkSpec = fmt.Sprintf("177.100.%d.0/24", GinkgoParallelNode())
			container, err = client.Create(garden.ContainerSpec{
				Network: networkSpec,
			})
			Expect(err).NotTo(HaveOccurred())
			contIfaceName = ethInterfaceName(container)

			networkBridgeName, err = container.Property("kawasaki.bridge-interface")
			Expect(err).NotTo(HaveOccurred())
		})

		var itCleansUpPerContainerNetworkingResources = func() {
			It("should remove iptable entries", func() {
				out, err := exec.Command("iptables", "-w", "-S", "-t", "filter").CombinedOutput()
				Expect(err).NotTo(HaveOccurred())
				Expect(string(out)).NotTo(MatchRegexp("w-%d-instance.* 177.100.%d.0/24", GinkgoParallelNode(), GinkgoParallelNode()))
			})

			It("should remove virtual ethernet cards", func() {
				ifconfigExits := func() int {
					session, err := gexec.Start(exec.Command("ifconfig", contIfaceName), GinkgoWriter, GinkgoWriter)
					Expect(err).NotTo(HaveOccurred())

					return session.Wait().ExitCode()
				}
				Eventually(ifconfigExits).ShouldNot(Equal(0))
			})
		}

		var itRemovesTheNetworkBridge = func() {
			It("should remove the network bridge", func() {
				session, err := gexec.Start(
					exec.Command("ip", "link", "show", networkBridgeName),
					GinkgoWriter, GinkgoWriter,
				)
				Expect(err).NotTo(HaveOccurred())

				session.Wait()
				Expect(session.ExitCode()).NotTo(Equal(0))
			})
		}

		Context("when destroy is called", func() {
			JustBeforeEach(func() {
				Expect(client.Destroy(container.Handle())).To(Succeed())
			})

			itCleansUpPerContainerNetworkingResources()
			itRemovesTheNetworkBridge()

			Context("and there was more than one containers in the same subnet", func() {
				var otherContainer garden.Container

				BeforeEach(func() {
					var err error

					otherContainer, err = client.Create(garden.ContainerSpec{
						Network: networkSpec,
					})
					Expect(err).NotTo(HaveOccurred())
				})

				JustBeforeEach(func() {
					Expect(client.Destroy(otherContainer.Handle())).To(Succeed())
				})

				itRemovesTheNetworkBridge()
			})
		})
	})
})

func ethInterfaceName(container garden.Container) string {
	buffer := gbytes.NewBuffer()
	proc, err := container.Run(
		garden.ProcessSpec{
			Path: "sh",
			Args: []string{"-c", "ifconfig | grep 'Ethernet' | cut -f 1 -d ' '"},
			User: "root",
		},
		garden.ProcessIO{
			Stdout: buffer,
			Stderr: GinkgoWriter,
		},
	)
	Expect(err).NotTo(HaveOccurred())
	Expect(proc.Wait()).To(Equal(0))

	contIfaceName := string(buffer.Contents()) // g3-abc-1

	return contIfaceName[:len(contIfaceName)-2] + "0" // g3-abc-0
}
