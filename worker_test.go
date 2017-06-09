package worker

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	MESSAGE_COUNT = 5
	WORKERS_COUNT = 2
)

type (
	FakeReceiver struct{}
	FakeHandler  struct {
		MessagesHandled int
		TimesCalled     int
	}
)

func (r *FakeReceiver) Receive() []Message {
	var m []Message
	for i := 0; i < MESSAGE_COUNT; i++ {
		str := "message: " + strconv.Itoa(i)
		m = append(m, Message{Body: &str})
	}

	return m
}

func (h *FakeHandler) Handle(m *[]Message) (bool, error) {
	h.TimesCalled = h.TimesCalled + 1
	h.MessagesHandled = h.MessagesHandled + len(*m)
	return true, nil
}

func TestShouldProcessMessages(t *testing.T) {
	w := &Worker{MaxWorkers: WORKERS_COUNT, GracefulBreak: true}
	h := &FakeHandler{}
	r := &FakeReceiver{}

	w.Run(r, h)

	assert.Equal(t, MESSAGE_COUNT*WORKERS_COUNT, h.MessagesHandled, "Number of messages handled")
	assert.Equal(t, WORKERS_COUNT, h.TimesCalled, "Number of workers called")
}
