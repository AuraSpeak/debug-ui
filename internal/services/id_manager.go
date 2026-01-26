package services

import "sync"

type ids struct {
	nextID int
	mu     sync.Mutex
}

var idCounter ids

func init() {
	idCounter = ids{
		nextID: 0,
		mu:     sync.Mutex{},
	}
}

func (i *ids) getNextID() int {
	i.mu.Lock()
	id := i.nextID
	i.nextID++
	i.mu.Unlock()
	return id
}

// GetNextID returns the next available ID (exported function)
func GetNextID() int {
	return idCounter.getNextID()
}
