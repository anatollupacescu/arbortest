package runner

// Testable exported.
type Testable interface {
	Failed() bool
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Log(args ...interface{})
}

// T exported.
type T struct {
	proxy Testable
}

// NewT exported.
func NewT(t Testable) *T {
	return &T{proxy: t}
}

// Failed exported.
func (a *T) Failed() bool {
	return a.proxy.Failed()
}

// Error exported.
func (a *T) Error(args ...interface{}) {
	a.proxy.Error(args...)
}

// Errorf exported.
func (a *T) Errorf(format string, args ...interface{}) {
	a.proxy.Errorf(format, args...)
}

// Log exported.
func (a *T) Log(args ...interface{}) {
	a.proxy.Log(args...)
}
