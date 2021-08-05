package youtubefile

type Commander func(name string, arg ...string) Outputer

type Outputer interface {
	Output() ([]byte, error)
	CombinedOutput() ([]byte, error)
}
