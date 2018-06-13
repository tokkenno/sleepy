package event

// Emit events
type Emitter struct {
	handler Handler
}

// Create a new event emitter
func NewEvent() *Emitter {
	return &Emitter{handler: Handler{}}
}

// Emit an event asynchronous
func (ee *Emitter) Emit(sender interface{}, args Args) {
	go ee.EmitSync(sender, args)
}

// Emit an event synchronous
func (ee *Emitter) EmitSync(sender interface{}, args Args) {
	for _, eventPtr := range ee.handler.events {
		eventPtr.shot(sender, args)
	}
}

// Get the associated handler instance
func (ee *Emitter) GetHandler() *Handler {
	return &ee.handler
}