package repository

import (
	errorpkg "blockchain-parser/internal/error"
	"context"
	"sync"

	"blockchain-parser/internal/entity"
)

type InMemSubscriber struct {
	data map[string]entity.Subscriber
	mu   sync.RWMutex
}

func NewInMemSubscriber() *InMemSubscriber {
	return &InMemSubscriber{
		data: map[string]entity.Subscriber{},
	}
}

func (r *InMemSubscriber) Save(_ context.Context, subscriber entity.Subscriber) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.data[subscriber.Address] = subscriber

	return nil
}

func (r *InMemSubscriber) Get(_ context.Context, address string) (entity.Subscriber, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	subscriber, ok := r.data[address]
	if !ok {
		return subscriber, errorpkg.SubscriberNotFound
	}

	return subscriber, nil
}
