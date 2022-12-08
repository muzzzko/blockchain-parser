package service

import (
	"context"

	"blockchain-parser/internal/entity"
)

type TransactionRepository interface {
	GetTxnsByAddress(ctx context.Context, address string) ([]entity.Transaction, error)
	Save(_ context.Context, transaction entity.Transaction) error
}

type SubscriberRepository interface {
	Save(_ context.Context, subscriber entity.Subscriber) error
	Get(_ context.Context, address string) (entity.Subscriber, error)
}

type BlockRepository interface {
	GetLastParsedBlock(ctx context.Context) (entity.Block, error)
	GetLastBlock(ctx context.Context) (entity.Block, error)
	GetFailedBlock(ctx context.Context) (entity.Block, error)

	Upsert(ctx context.Context, block entity.Block) error
}

type BlockChainClient interface {
	GetBlockNumber(ctx context.Context) (int, error)
	GetTxnsByBlockByNumber(ctx context.Context, blockNumber int) ([]entity.Transaction, error)
}

type Locker interface {
	Lock()
	Unlock()
}
