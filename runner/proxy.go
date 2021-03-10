package runner

import "testing"

// Testable exported.
type Testable interface {
	Failed() bool
	Error(...interface{})
	Errorf(string, ...interface{})
	Log(...interface{})
	Run(string, func(t *testing.T)) bool
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

// Run exported.
func (a *T) Run(name string, f func(t *testing.T)) bool {
	f(nil)
	return true
}
