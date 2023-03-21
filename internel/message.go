package internel

import "fmt"

type Message struct {
	uid  string
	data any
}

func WrapMsg(uid string, data any) Message {
	return Message{
		uid:  uid,
		data: data,
	}
}

func (m Message) String() string {
	return fmt.Sprintf("%v", m.data)
}

// StartMessage Define a start message struct
type StartMessage struct{}

// StopMessage Define a stop message struct
type StopMessage struct{}
