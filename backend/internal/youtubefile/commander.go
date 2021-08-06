package youtubefile

import "io"

type Commander func(name string, arg ...string) Outputer

type Outputer interface {
	Output() ([]byte, error)
	CombinedOutput() ([]byte, error)
	StdoutPipe() (io.ReadCloser, error)
	StderrPipe() (io.ReadCloser, error)
	Start() error
	Wait() error
	String() string
}
