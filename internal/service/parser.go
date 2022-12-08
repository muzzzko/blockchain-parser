package service

import (
	"context"
	"log"

	"blockchain-parser/internal/entity"
)

// in my view methods shoud return error
type Parser struct {
	txnRepo        TransactionRepository
	subscriberRepo SubscriberRepository
	blockRepo      BlockRepository
}

func NewParser(
	txnRepo TransactionRepository,
	subscriberRepo SubscriberRepository,
	blockRepo BlockRepository,
) *Parser {

	return &Parser{
		txnRepo:        txnRepo,
		subscriberRepo: subscriberRepo,
		blockRepo:      blockRepo,
	}
}

func (p *Parser) GetCurrentBlock() int {
	block, err := p.blockRepo.GetLastParsedBlock(context.Background())
	if err != nil {
		log.Printf("fail get current block: %s\n", err)

		return 0
	}

	return block.Number
}

func (p *Parser) Subscribe(address string) bool {
	subscriber := entity.Subscriber{
		Address: address,
	}
	err := p.subscriberRepo.Save(context.Background(), subscriber)
	if err != nil {
		log.Printf("fail subscribe address (%s): %s\n", address, err)

		return false
	}

	return true
}

func (p *Parser) GetTransactions(address string) []entity.Transaction {
	txns, err := p.txnRepo.GetTxnsByAddress(context.Background(), address)
	if err != nil {
		log.Printf("fail get trasactions for address (%s): %s\n", address, err)
	}

	return txns
}
