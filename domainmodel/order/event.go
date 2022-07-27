package order

type Event string

func NewEvent(msg string) Event {
	return Event(msg)
}
