package worker

import (
	"fmt"
	"sync"
)

type Pool struct {
	capacity         int
	messageQueueChan chan interface{}
	wg               sync.WaitGroup
	processFunc      func(message interface{}) error
}

func NewWorkerPool(capacity int, messageQueueSize int, processFunc func(message interface{}) error) *Pool {
	return &Pool{
		capacity:         capacity,
		messageQueueChan: make(chan interface{}, messageQueueSize),
		processFunc:      processFunc,
	}
}

func (p *Pool) PushMessage(message interface{}) {
	p.messageQueueChan <- message
}

func (p *Pool) Start() {
	for i := 0; i < p.capacity; i++ {
		p.wg.Add(1)
		go p.handleMessage(i, p.messageQueueChan)
	}
}

func (p *Pool) handleMessage(id int, channel <-chan interface{}) {
	defer p.wg.Done()
	for message := range channel {
		fmt.Printf("worker %v is processing message %v\n", id, message)
		err := p.processFunc(message)

		if err != nil {
			fmt.Printf("worker %v failed to process message %v with error: %v\n", id, message, err)
		}
	}
}

// Shutdown implement graceful shutdown
func (p *Pool) GracefulStop() {
	fmt.Println("Closing working pool...")
	close(p.messageQueueChan)
	p.wg.Wait()
	fmt.Println("Worker pool gracefully stopped")
}
