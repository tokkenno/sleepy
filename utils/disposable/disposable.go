package disposable

// Disposable object interface
type Disposable interface {
	Dispose()
}

// Use a disposable variable inside function and then release it
func Using(d Disposable, f func(Disposable)) {
	defer d.Dispose()
	f(d)
}