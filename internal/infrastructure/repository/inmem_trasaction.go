package repository

import (
	"context"
	"fmt"
	"sync"

	"blockchain-parser/internal/entity"
)

type InMemTransaction struct {
	data map[string]map[string]*entity.Transaction
	mu   sync.RWMutex
}

func NewInMemTransaction() *InMemTransaction {
	return &InMemTransaction{
		data: map[string]map[string]*entity.Transaction{},
	}
}

func (r *InMemTransaction) Save(_ context.Context, transaction entity.Transaction) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.data[transaction.To]; !ok {
		r.data[transaction.To] = make(map[string]*entity.Transaction)
	}
	if _, ok := r.data[transaction.From]; !ok {
		r.data[transaction.From] = make(map[string]*entity.Transaction)
	}

	txnID := fmt.Sprintf("%d_%d", transaction.BlockNumber, transaction.TransactionIndex)
	r.data[transaction.To][txnID] = &transaction
	r.data[transaction.From][txnID] = &transaction

	return nil
}

func (r *InMemTransaction) GetTxnsByAddress(_ context.Context, address string) ([]entity.Transaction, error) {
	r.mu.RLock()
	txns := r.data[address]
	r.mu.RUnlock()

	txnscopy := make([]entity.Transaction, 0, len(txns))
	for _, txn := range txns {
		txnscopy = append(txnscopy, *txn)
	}

	return txnscopy, nil
}
