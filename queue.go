package dynamicqueue

import (
	"sync"
)

type DynamicQueue interface {
	Add(item interface{})
	Stop()
}

type QueueHandler interface {
	OnItem(item interface{}) // will be run in parallel
}

func NewDynamicQueue(workersCount int, handler QueueHandler) DynamicQueue {
	q := &dynamicQueue{
		handler:        handler,
		newItems:       make(chan interface{}),
	}
	q.itemsToProcess = NewDynamicallyBufferedChannel(q.newItems)
	q.WaitGroup.Add(workersCount)
	for i := 0; i < workersCount; i++ {
		go q.worker()
	}

	return q
}

type dynamicQueue struct {
	newItems       chan interface{}
	itemsToProcess <- chan interface{}
	handler        QueueHandler
	sync.WaitGroup
}

func (dq *dynamicQueue) worker() {
	defer dq.Done()
	for item := range dq.itemsToProcess {
		dq.handler.OnItem(item)
	}
}

func (dq *dynamicQueue) Add(item interface{}) {
	dq.newItems <- item
}

func (dq *dynamicQueue) Stop() {
	close(dq.newItems)
	dq.Wait()
}
