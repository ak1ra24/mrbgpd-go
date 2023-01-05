package event

type EventQueue []Event

func NewEventQueue() *EventQueue {
	eventQueue := EventQueue{}
	return &eventQueue
}

func (eq *EventQueue) Enqueue(event Event) {
	*eq = append(*eq, event)
}

func (eq *EventQueue) Dequeue() Event {
	result := (*eq)[0]
	*eq = (*eq)[1:]
	return result
}
