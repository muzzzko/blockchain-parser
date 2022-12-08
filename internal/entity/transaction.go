package entity

type Transaction struct {
	From             string
	To               string
	Value            string
	BlockNumber      int
	TransactionIndex int
}
