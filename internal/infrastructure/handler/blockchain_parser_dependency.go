package handler

import "blockchain-parser/internal/entity"

type Parser interface {
	GetCurrentBlock() int
	Subscribe(address string) bool
	GetTransactions(address string) []entity.Transaction
}
