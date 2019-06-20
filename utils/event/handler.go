package event

import "sync"

// Handle events
type Handler struct {
	events    []*Container
	eventLock sync.Mutex
}

// Subscribe a callback to be fired when the event halts, and get
// the internal callback container pointer, to unsubscribe
func (eh *Handler) Listen(callback Callback) *Container {
	ev := newEvent(eh, callback)
	eh.eventLock.Lock()
	eh.events = append(eh.events, ev)
	eh.eventLock.Unlock()
	return ev
}

// Unsubscribe a callback through his [container] instance
func (eh *Handler) Ignore(container *Container) {
	eh.eventLock.Lock()
	for eventPos, eventPtr := range eh.events {
		if eventPtr == container {
			eh.events = append(eh.events[:eventPos], eh.events[eventPos+1:]...)
			eh.eventLock.Unlock()
			return
		}
	}
	eh.eventLock.Unlock()
}
