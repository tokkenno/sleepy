package event

type Args interface{}

// Callback function definition
type Callback func(sender interface{}, args Args)

// Callback container for event handler
type Container struct {
	callback Callback
	handler  *Handler
}

// Create a new event container from a callback
func newEvent(handler *Handler, callback Callback) *Container {
	return &Container{
		callback: callback,
		handler:  handler,
	}
}

// Execute the callback
func (cont *Container) shot(sender interface{}, args Args) {
	cont.callback(sender, args)
}

// Unsubscribe a callback through his [container] instance
func (cont *Container) Ignore() {
	cont.handler.Ignore(cont)
}
