package httpclient

import "blockchain-parser/internal/entity"

func mapResponseToTxns(resp ethereumGetBlockByNumberResponse) []entity.Transaction {
	txns := make([]entity.Transaction, 0, len(resp.Result.Transactions))

	for _, resptxn := range resp.Result.Transactions {
		txn := entity.Transaction{
			From:  resptxn.From,
			To:    resptxn.To,
			Value: resptxn.Value,
		}

		txns = append(txns, txn)
	}

	return txns
}
