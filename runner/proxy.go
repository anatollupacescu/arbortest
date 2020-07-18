package runner

type Testable interface {
	Failed() bool
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Log(args ...interface{})
}

type T struct {
	proxy Testable
}

func NewT(t Testable) *T {
	return &T{proxy: t}
}

func (a *T) Failed() bool {
	return a.proxy.Failed()
}

func (a *T) Error(args ...interface{}) {
	a.proxy.Error(args...)
}

func (a *T) Errorf(format string, args ...interface{}) {
	a.proxy.Errorf(format, args...)
}

func (a *T) Log(args ...interface{}) {
	a.proxy.Log(args...)
}
