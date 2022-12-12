package handler

import "blockchain-parser/internal/entity"

type blockChainParserGetCurrentBlockResponse struct {
	Block string `json:"block"`
}

type blockChainParserGetTransactionsTransactions struct {
	From  string `json:"from"`
	To    string `json:"to"`
	Value string `json:"value"`
}

type blockChainParserGetTransactionsResponse struct {
	Transactions []blockChainParserGetTransactionsTransactions `json:"transactions"`
}

func mapTransactionsToGetTransactionsResponse(txns []entity.Transaction) blockChainParserGetTransactionsResponse {
	resp := blockChainParserGetTransactionsResponse{
		Transactions: make([]blockChainParserGetTransactionsTransactions, 0, len(txns)),
	}

	for _, txn := range txns {
		resp.Transactions = append(resp.Transactions, blockChainParserGetTransactionsTransactions{
			From:  txn.From,
			To:    txn.To,
			Value: txn.Value,
		})
	}

	return resp
}
