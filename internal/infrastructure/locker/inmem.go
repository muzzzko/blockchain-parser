package locker

import "sync"

type InMem struct {
	mu sync.Mutex
}

func NewInMem() *InMem {
	return &InMem{}
}

func (l *InMem) Lock() {
	l.mu.Lock()
}

func (l *InMem) Unlock() {
	l.mu.Unlock()
}
