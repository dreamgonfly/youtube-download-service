package youtubefile

import (
	"io"
	"os/exec"
)

func Command(name string, arg ...string) Outputer {
	return &Cmd{exec.Command(name, arg...)}
}

type Cmd struct{ Cmd *exec.Cmd }

func (c *Cmd) Output() ([]byte, error) {
	return c.Cmd.Output()
}

func (c *Cmd) CombinedOutput() ([]byte, error) {
	return c.Cmd.CombinedOutput()
}

func (c *Cmd) StdoutPipe() (io.ReadCloser, error) {
	return c.Cmd.StdoutPipe()
}

func (c *Cmd) StderrPipe() (io.ReadCloser, error) {
	return c.Cmd.StderrPipe()
}

func (c *Cmd) Start() error {
	return c.Cmd.Start()
}

func (c *Cmd) Wait() error {
	return c.Cmd.Wait()
}

func (c *Cmd) String() string {
	return c.Cmd.String()
}