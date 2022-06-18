package event_publisher

type EventPublisher struct {
	subscriberId eventSubscriberId
	subscribers  map[string][]eventSubscriberId
	callbacks    map[eventSubscriberId]func([]byte)
}

type eventSubscriberId int

func New() *EventPublisher {
	return &EventPublisher{
		subscriberId: 0,
		subscribers:  make(map[string][]eventSubscriberId),
		callbacks:    make(map[eventSubscriberId]func([]byte)),
	}
}

func (ep *EventPublisher) Subscribe(event string, cb func([]byte)) func() {
	id := ep.subscriberId
	ep.subscriberId++

	if _, ok := ep.subscribers[event]; !ok {
		ep.subscribers[event] = make([]eventSubscriberId, 0, 1)
	}
	ep.callbacks[id] = cb
	ep.subscribers[event] = append(ep.subscribers[event], id)

	unsubscribe := func() { delete(ep.callbacks, id) }
	return unsubscribe
}

func (ep *EventPublisher) Publish(event string, payload []byte) {
	if _, ok := ep.subscribers[event]; !ok {
		return
	}

	for i := range ep.subscribers[event] {
		id := ep.subscribers[event][i]
		if ep.callbacks[id] != nil {
			go ep.callbacks[id](payload)
		}
	}
}
