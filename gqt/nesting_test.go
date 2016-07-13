package gqt_test

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"

	"io/ioutil"

	"github.com/cloudfoundry-incubator/garden"
	gclient "github.com/cloudfoundry-incubator/garden/client"
	gconn "github.com/cloudfoundry-incubator/garden/client/connection"
	"github.com/cloudfoundry-incubator/guardian/gqt/runner"
)

var _ = FDescribe("When nested", func() {
	nestedRootfsPath := os.Getenv("GARDEN_NESTABLE_TEST_ROOTFS")
	var client *runner.RunningGarden
	BeforeEach(func() {
		client = startGarden()
	})

	startNestedGarden := func() (garden.Container, string) {
		absoluteGardenPath, err := filepath.Abs(gardenBin)
		Expect(err).ToNot(HaveOccurred())

		// TODO: WHICH RUNC
		absoluteRuncPath, err := filepath.Abs("/usr/local/bin/runc")
		Expect(err).ToNot(HaveOccurred())

		absoluteIODaemonPath, err := filepath.Abs(iodaemonBin)
		Expect(err).ToNot(HaveOccurred())

		absoluteDadooPath, err := filepath.Abs(dadooBin)
		Expect(err).ToNot(HaveOccurred())

		absoluteInitPath, err := filepath.Abs(initBin)
		Expect(err).ToNot(HaveOccurred())

		absoluteTarPath, err := filepath.Abs(runner.TarBin)
		Expect(err).ToNot(HaveOccurred())

		absoluteNstarPath, err := filepath.Abs(nstarBin)
		Expect(err).ToNot(HaveOccurred())

		container, err := client.Create(garden.ContainerSpec{
			RootFSPath: nestedRootfsPath,
			// only privileged containers support nesting
			Privileged: true,
			BindMounts: []garden.BindMount{
				{
					SrcPath: filepath.Dir(absoluteGardenPath),
					DstPath: "/root/bin/",
					Mode:    garden.BindMountModeRW,
				},
				{
					SrcPath: filepath.Dir(absoluteRuncPath),
					DstPath: "/root/bin/runc/",
					Mode:    garden.BindMountModeRW,
				},
				{
					SrcPath: filepath.Dir(absoluteIODaemonPath),
					DstPath: "/root/bin/iodaemon/",
					Mode:    garden.BindMountModeRW,
				},
				{
					SrcPath: filepath.Dir(absoluteDadooPath),
					DstPath: "/root/bin/dadoo/",
					Mode:    garden.BindMountModeRW,
				},
				{
					SrcPath: filepath.Dir(absoluteInitPath),
					DstPath: "/root/bin/init/",
					Mode:    garden.BindMountModeRW,
				},
				{
					SrcPath: filepath.Dir(absoluteTarPath),
					DstPath: "/root/bin/tar/",
					Mode:    garden.BindMountModeRW,
				},
				{
					SrcPath: filepath.Dir(absoluteNstarPath),
					DstPath: "/root/bin/nstar/",
					Mode:    garden.BindMountModeRW,
				},
				{
					SrcPath: runner.RootFSPath,
					DstPath: "/root/rootfs",
					Mode:    garden.BindMountModeRW,
				},
			},
		})
		Expect(err).ToNot(HaveOccurred())

		nestedServerOutput := gbytes.NewBuffer()

		_, err = container.Run(garden.ProcessSpec{
			Path: "sh",
			User: "root",
			Dir:  "/root",
			Args: []string{
				"-c",
				fmt.Sprintf(`
				set -e

				tmpdir=/tmp/dir
				rm -fr $tmpdir
				mkdir $tmpdir
				mount -t tmpfs none $tmpdir
				echo "{}" > /root/network.props

				mkdir $tmpdir/depot
				mkdir $tmpdir/snapshots
				mkdir $tmpdir/state
				mkdir $tmpdir/graph

				/root/bin/guardian \
					--default-rootfs /root/rootfs \
					--depot $tmpdir/depot \
					--graph $tmpdir/graph \
					--tag n \
					--bind-socket tcp \
					--bind-ip 0.0.0.0 \
					--bind-port 7778 \
					--network-pool 10.254.6.0/22 \
					--runc-bin /root/bin/runc/runc \
					--init-bin /root/bin/init/init \
					--iodaemon-bin /root/bin//iodaemon/iodaemon \
					--dadoo-bin /root/bin/dadoo/dadoo \
					--nstar-bin /root/bin/nstar/nstar \
					--port-pool-properties-path /root/network.props \
					--tar-bin /root/bin/tar/tar \
					--port-pool-start 30000
				`),
			},
		}, garden.ProcessIO{
			Stdout: io.MultiWriter(nestedServerOutput, gexec.NewPrefixedWriter("\x1b[32m[o]\x1b[34m[nested-garden-runc]\x1b[0m ", GinkgoWriter)),
			Stderr: gexec.NewPrefixedWriter("\x1b[91m[e]\x1b[34m[nested-garden-runc]\x1b[0m ", GinkgoWriter),
		})

		info, err := container.Info()
		Expect(err).ToNot(HaveOccurred())

		nestedGardenAddress := fmt.Sprintf("%s:7778", info.ContainerIP)
		Eventually(nestedServerOutput, "60s").Should(gbytes.Say("guardian.started"))

		return container, nestedGardenAddress
	}

	FIt("can start a nested garden server and run a container inside it", func() {
		container, nestedGardenAddress := startNestedGarden()
		defer func() {
			Expect(client.Destroy(container.Handle())).To(Succeed())
		}()

		nestedClient := gclient.New(gconn.New("tcp", nestedGardenAddress))
		nestedContainer, err := nestedClient.Create(garden.ContainerSpec{})
		Expect(err).ToNot(HaveOccurred())

		nestedOutput := gbytes.NewBuffer()
		_, err = nestedContainer.Run(garden.ProcessSpec{
			User: "root",
			Path: "/bin/echo",
			Args: []string{
				"I am nested!",
			},
		}, garden.ProcessIO{Stdout: nestedOutput, Stderr: nestedOutput})
		Expect(err).ToNot(HaveOccurred())

		Eventually(nestedOutput, "60s").Should(gbytes.Say("I am nested!"))
	})

	Context("when cgroup limits are applied to the parent garden process", func() {
		devicesCgroupNode := func() string {
			contents, err := ioutil.ReadFile("/proc/self/cgroup")
			Expect(err).ToNot(HaveOccurred())
			for _, line := range strings.Split(string(contents), "\n") {
				if strings.Contains(line, "devices:") {
					lineParts := strings.Split(line, ":")
					Expect(lineParts).To(HaveLen(3))
					return lineParts[2]
				}
			}
			Fail("could not find devices cgroup node")
			return ""
		}

		It("passes on these limits to the child container", func() {
			// When this test is run in garden (e.g. in Concourse), we cannot create more permissive device cgroups
			// than are allowed in the outermost container. So we apply this rule to the outermost container's cgroup
			cmd := exec.Command(
				"sh",
				"-c",
				fmt.Sprintf("echo 'b 7:200 r' > /tmp/garden-%d/cgroup/devices%s/devices.allow", GinkgoParallelNode(), devicesCgroupNode()),
			)
			cmd.Stdout = GinkgoWriter
			cmd.Stderr = GinkgoWriter
			Expect(cmd.Run()).To(Succeed())

			gardenInContainer, nestedGardenAddress := startNestedGarden()
			defer client.Destroy(gardenInContainer.Handle())

			postProc, err := gardenInContainer.Run(garden.ProcessSpec{
				Path: "bash",
				User: "root",
				Args: []string{"-c",
					`
				cgroup_path_segment=$(cat /proc/self/cgroup | grep devices: | cut -d ':' -f 3)
				echo "b 7:200 r" > /tmp/garden-n/cgroup/devices${cgroup_path_segment}/devices.allow
				`},
			}, garden.ProcessIO{
				Stdout: GinkgoWriter,
				Stderr: GinkgoWriter,
			})
			Expect(err).ToNot(HaveOccurred())
			Expect(postProc.Wait()).To(Equal(0))

			nestedClient := gclient.New(gconn.New("tcp", nestedGardenAddress))
			nestedContainer, err := nestedClient.Create(garden.ContainerSpec{
				Privileged: true,
			})
			Expect(err).ToNot(HaveOccurred())

			nestedProcess, err := nestedContainer.Run(garden.ProcessSpec{
				User: "root",
				Path: "sh",
				Args: []string{"-c", `
				mknod ./foo b 7 200
				cat foo > /dev/null
				`},
			}, garden.ProcessIO{
				Stdout: GinkgoWriter,
				Stderr: GinkgoWriter,
			})
			Expect(err).ToNot(HaveOccurred())

			Expect(nestedProcess.Wait()).To(Equal(0))
		})
	})
})
