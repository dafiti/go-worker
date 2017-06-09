package worker

import (
	"sync"
)

type (
	Worker struct {
		MaxWorkers    int
		GracefulBreak bool
	}

	Message struct {
		Body *string
	}

	MessagesReceiver interface {
		Receive() []Message
	}

	MessagesHandler interface {
		Handle(messages *[]Message) (bool, error)
	}
)

func (w *Worker) process(wg *sync.WaitGroup, r MessagesReceiver, h MessagesHandler) {
	defer wg.Done()

	messages := r.Receive()
	h.Handle(&messages)
}

func (w *Worker) Run(receiver MessagesReceiver, handler MessagesHandler) {
	var wg sync.WaitGroup

	for {
		wg.Add(w.MaxWorkers)

		for i := 0; i < w.MaxWorkers; i++ {
			go w.process(&wg, receiver, handler)
		}

		wg.Wait()

		if w.GracefulBreak == true {
			break
		}
	}
}
