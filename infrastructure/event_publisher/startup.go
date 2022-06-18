package event_publisher

func StartUp() (*EventPublisher, error) {
	ep := New()
	return ep, nil
}
