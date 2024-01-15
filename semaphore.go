package main

type Semaphore chan struct{}

func NewSemaphore(size int) Semaphore {
	return make(chan struct{}, size)
}

func (s Semaphore) Acquire() {
	s <- struct{}{}
	// log.Println("Acquired")
}

func (s Semaphore) Release() {
	<-s
	// log.Println("Released")
}
