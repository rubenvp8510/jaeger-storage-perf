package queue

import (
	"github.com/jaegertracing/jaeger/model"
	"go.uber.org/atomic"
	"sync"
)

type Queue struct {
	Dropped atomic.Uint32
	Added   atomic.Uint32
	workers int
	queue   *chan model.Span
	stop    chan bool
	waiter  sync.WaitGroup
}

func NewQueue(numWorkers int, capacity int) *Queue{
	queue := make(chan model.Span, capacity)
	return &Queue{
		workers:numWorkers,
		queue:&queue,
		stop:make(chan bool, 1),
	}
}

func (q *Queue) Start(writeFunction func(span model.Span) error) {
	var startWg sync.WaitGroup
	for i := 0; i < q.workers; i++ {
		q.waiter.Add(1)
		startWg.Add(1)
		go func() {
			startWg.Done()
			defer q.waiter.Done()
			for span := range *q.queue {
				err := writeFunction(span)
				if err != nil {
					q.Dropped.Inc()
				} else {
					q.Added.Inc()
				}
			}
		}()
	}
	startWg.Wait()
}

func (q *Queue) Stop() {
	close(*q.queue)
	q.waiter.Wait()
}


func (q *Queue) Enqueue(span model.Span) {
	*q.queue <- span
}
