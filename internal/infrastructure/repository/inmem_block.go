package repository

import (
	"context"
	"sync"
	"time"

	"blockchain-parser/internal/constant"
	"blockchain-parser/internal/entity"
	errorpkg "blockchain-parser/internal/error"
)

const (
	processingTTL = time.Minute * 5
)

type InMemBlock struct {
	failedBlocks     map[int]entity.Block
	processingBlocks map[int]entity.Block
	parsedBlocks     map[int]entity.Block

	parsedBlockNumber     int
	processingBlockNumber int

	mu sync.Mutex
}

func NewInMemBlock() *InMemBlock {
	return &InMemBlock{
		failedBlocks:     map[int]entity.Block{},
		processingBlocks: map[int]entity.Block{},
		parsedBlocks:     map[int]entity.Block{},
	}
}

func (r *InMemBlock) GetLastParsedBlock(_ context.Context) (entity.Block, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	block, ok := r.parsedBlocks[r.parsedBlockNumber]
	if !ok {
		return entity.Block{}, errorpkg.BlockNotFound
	}

	return block, nil
}

func (r *InMemBlock) GetLastBlock(_ context.Context) (entity.Block, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.processingBlockNumber > r.parsedBlockNumber {
		if processingBlock, ok := r.processingBlocks[r.processingBlockNumber]; ok {
			return processingBlock, nil
		}
	}

	if parsedBlock, ok := r.parsedBlocks[r.parsedBlockNumber]; ok {
		return parsedBlock, nil
	}

	return entity.Block{}, errorpkg.BlockNotFound
}

func (r *InMemBlock) GetFailedBlock(_ context.Context) (entity.Block, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, failedBlock := range r.failedBlocks {
		return failedBlock, nil
	}

	for _, processingBlock := range r.processingBlocks {
		if time.Since(processingBlock.UpdatedAt) > processingTTL {
			return processingBlock, nil
		}
	}

	return entity.Block{}, errorpkg.BlockNotFound
}

func (r *InMemBlock) Upsert(_ context.Context, block entity.Block) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	switch block.Status {
	case constant.BlockStatusProcessing:
		delete(r.failedBlocks, block.Number)

		r.processingBlocks[block.Number] = block

		if block.Number > r.processingBlockNumber {
			r.processingBlockNumber = block.Number
		}
	case constant.BlockStatusFailed:
		delete(r.processingBlocks, block.Number)

		r.failedBlocks[block.Number] = block

		if r.processingBlockNumber == block.Number {
			r.processingBlockNumber--
		}
	case constant.BlockStatusParsed:
		delete(r.processingBlocks, block.Number)

		r.parsedBlocks[block.Number] = block

		if block.Number > r.parsedBlockNumber {
			r.parsedBlockNumber = block.Number
		}
	default:
		return errorpkg.UnknownBlockStatus
	}

	return nil
}
