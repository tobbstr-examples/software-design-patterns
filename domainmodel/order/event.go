package order

type Message string

func NewMessage(msg string) Message {
	return Message(msg)
}
