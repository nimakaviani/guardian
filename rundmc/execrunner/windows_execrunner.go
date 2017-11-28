package execrunner

import (
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"syscall"

	"github.com/pkg/errors"

	"code.cloudfoundry.org/commandrunner"
	"code.cloudfoundry.org/garden"
	"code.cloudfoundry.org/guardian/logging"
	"code.cloudfoundry.org/lager"
)

type DirectExecRunner struct {
	RuntimePath   string
	CommandRunner commandrunner.CommandRunner
}

func (e *DirectExecRunner) Run(
	log lager.Logger, processID, processPath, sandboxHandle, _ string,
	_, _ uint32, pio garden.ProcessIO, _ bool, procJSON io.Reader,
	_ func() error,
) (garden.Process, error) {
	log = log.Session("execrunner")

	log.Info("start")
	defer log.Info("done")

	specPath := filepath.Join(processPath, "spec.json")
	specFile, err := os.OpenFile(specPath, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return nil, errors.Wrap(err, "opening process spec file for writing")
	}
	defer specFile.Close()
	if _, err := io.Copy(specFile, procJSON); err != nil {
		return nil, errors.Wrap(err, "writing process spec")
	}

	logPath := filepath.Join(processPath, "exec.log")

	proc := &process{id: processID}

	cmd := exec.Command(e.RuntimePath, "--debug", "--log", logPath, "--log-format", "json", "exec", "-p", specPath, "--pid-file", filepath.Join(processPath, "pidfile"), sandboxHandle)
	cmd.Stdout = pio.Stdout
	cmd.Stderr = pio.Stderr
	if err := e.CommandRunner.Start(cmd); err != nil {
		return nil, errors.Wrap(err, "execing runtime plugin")
	}

	proc.mux.Lock()

	go func() {
		if err := e.CommandRunner.Wait(cmd); err != nil {
			if exitErr, ok := err.(*exec.ExitError); ok {
				if status, ok := exitErr.Sys().(syscall.WaitStatus); ok {
					proc.exitCode = status.ExitStatus()
				} else {
					proc.exitCode = 1
					proc.exitErr = errors.New("couldn't get WaitStatus")
				}
			} else {
				proc.exitCode = 1
				proc.exitErr = err
			}
		}
		forwardLogs(log, logPath)
		proc.mux.Unlock()
	}()

	return proc, nil
}

func (e *DirectExecRunner) Attach(log lager.Logger, processID string, io garden.ProcessIO, processesPath string) (garden.Process, error) {
	panic("not supported on this platform")
}

type process struct {
	id       string
	exitCode int
	exitErr  error
	mux      sync.RWMutex
}

func (p *process) ID() string {
	return p.id
}

func (p *process) Wait() (int, error) {
	p.mux.RLock()
	defer p.mux.RUnlock()

	return p.exitCode, p.exitErr
}

func (p *process) SetTTY(ttySpec garden.TTYSpec) error {
	return nil
}

func (p *process) Signal(signal garden.Signal) error {
	return nil
}

func forwardLogs(log lager.Logger, logPath string) {
	defer os.Remove(logPath)

	buff, readErr := ioutil.ReadFile(logPath)
	if readErr != nil {
		log.Error("reading log file", readErr)
	}

	logging.ForwardRuncLogsToLager(log, "exec", buff)
}
