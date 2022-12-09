package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"blockchain-parser/internal/constant"
	"blockchain-parser/internal/entity"
	errorpkg "blockchain-parser/internal/error"
)

type ParserWorker struct {
	txnRepo          TransactionRepository
	subscriberRepo   SubscriberRepository
	blockRepo        BlockRepository
	blockChainClient BlockChainClient
	locker           Locker
}

func NewParserWorker(
	txnRepo TransactionRepository,
	subscriberRepo SubscriberRepository,
	blockRepo BlockRepository,
	blockChainClient BlockChainClient,
	locker Locker,
) *ParserWorker {
	return &ParserWorker{
		txnRepo:          txnRepo,
		subscriberRepo:   subscriberRepo,
		blockRepo:        blockRepo,
		blockChainClient: blockChainClient,
		locker:           locker,
	}
}

func (w *ParserWorker) Run(ctx context.Context) error {
	blockNumber, err := w.blockChainClient.GetBlockNumber(ctx)
	if err != nil {
		return fmt.Errorf("fail get block number in ParserWorker: %w", err)
	}

	block, err := w.getProcessingBlock(ctx, blockNumber)
	if err != nil {
		return err
	}

	if err := w.processBlock(ctx, block); err != nil {
		w.failBlockProcessing(ctx, block)

		return err
	}

	w.markBlockAsParsed(ctx, block)

	return nil
}

func (w *ParserWorker) processBlock(ctx context.Context, block entity.Block) error {
	txns, err := w.blockChainClient.GetTxnsByBlockByNumber(ctx, block.Number)
	if err != nil {
		return fmt.Errorf("fail get transactions in ParserWorker: %w", err)
	}

	for _, txn := range txns {
		toOk, err := w.checkSubscription(ctx, txn.To)
		if err != nil {
			return fmt.Errorf("fail get check subscription in ParserWorker: %w", err)
		}

		fromOk, err := w.checkSubscription(ctx, txn.From)
		if err != nil {
			return fmt.Errorf("fail get check subscription in ParserWorker: %w", err)
		}

		if toOk || fromOk {
			if err := w.txnRepo.Save(ctx, txn); err != nil {
				return fmt.Errorf("fail save trasaction in ParserWorker: %w", err)
			}
		}
	}

	return nil
}

func (w *ParserWorker) getProcessingBlock(ctx context.Context, blockNumber int) (entity.Block, error) {
	w.locker.Lock()
	defer w.locker.Unlock()

	block, err := w.blockRepo.GetFailedBlock(ctx)
	if errors.Is(err, errorpkg.BlockNotFound) {
		lastBlock, err := w.blockRepo.GetLastBlock(ctx)
		if err != nil {
			return entity.Block{}, fmt.Errorf("fail get last block in getProcessingBlock: %w", err)
		}

		if lastBlock.Number == blockNumber {
			return entity.Block{}, fmt.Errorf("no block for parsing: %w", errorpkg.NoBlockForParsing)
		}

		block = entity.Block{
			Number: lastBlock.Number + 1,
		}
	}

	block.Status = constant.BlockStatusProcessing
	block.UpdatedAt = time.Now()

	if err = w.blockRepo.Upsert(ctx, block); err != nil {
		return entity.Block{}, fmt.Errorf("fail save block in getProcessingBlock: %w", err)
	}

	return block, nil
}

func (w *ParserWorker) failBlockProcessing(ctx context.Context, block entity.Block) {
	block.Status = constant.BlockStatusFailed
	block.UpdatedAt = time.Now()

	if err := w.blockRepo.Upsert(ctx, block); err != nil {
		log.Printf("fail save block (%d) in failBlockProcessing: %s\n", block.Number, err)
	}
}

func (w *ParserWorker) markBlockAsParsed(ctx context.Context, block entity.Block) {
	block.Status = constant.BlockStatusParsed
	block.UpdatedAt = time.Now()

	if err := w.blockRepo.Upsert(ctx, block); err != nil {
		log.Printf("fail save block (%d) in markBlockAsParsed: %s\n", block.Number, err)
	}
}

func (w *ParserWorker) checkSubscription(ctx context.Context, address string) (bool, error) {
	_, err := w.subscriberRepo.Get(ctx, address)
	if err != nil {
		if errors.Is(err, errorpkg.SubscriberNotFound) {
			return false, nil
		}

		return false, err
	}

	return true, err
}
