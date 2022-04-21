package event_publisher

type EventPublisher struct {
	subscribers map[string][]int
	callbacks   []func([]byte)
}

func New() *EventPublisher {
	return &EventPublisher{
		subscribers: make(map[string][]int),
		callbacks:   make([]func([]byte), 0),
	}
}

func (ep *EventPublisher) Subscribe(event string, cb func([]byte)) func() {
	if _, ok := ep.subscribers[event]; !ok {
		ep.subscribers[event] = make([]int, 0, 1)
	}

	id := len(ep.callbacks)
	ep.callbacks = append(ep.callbacks, cb)
	ep.subscribers[event] = append(ep.subscribers[event], id)

	return func() {
		// unsubscribe
		ep.callbacks[id] = nil
		for i := range ep.subscribers[event] {
			if ep.subscribers[event][i] == id {
				ep.subscribers[event][i] = -1
			}
		}
	}
}

func (ep *EventPublisher) Publish(event string, payload []byte) {
	if _, ok := ep.subscribers[event]; !ok {
		return
	}

	for i := range ep.subscribers[event] {
		id := ep.subscribers[event][i]
		if id != -1 {
			ep.callbacks[id](payload)
		}
	}
}
