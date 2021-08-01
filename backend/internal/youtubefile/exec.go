package youtubefile

import "os/exec"

func Command(name string, arg ...string) Outputer {
	return &Cmd{exec.Command(name, arg...)}
}

type Cmd struct{ Cmd *exec.Cmd }

func (c *Cmd) Output() ([]byte, error) {
	return c.Cmd.Output()
}
