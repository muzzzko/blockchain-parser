package httpclient

type EthereumError struct {
	Code    int64
	Message string
}

type EthereumTxn struct {
	From  string
	To    string
	Value string
}

type EthereumGetBlockByNumberResult struct {
	Transactions []EthereumTxn
}

type ethereumGetBlockNumberResponse struct {
	ID     int32
	Result *string        `json:",omitempty"`
	Error  *EthereumError `json:",omitempty"`
}

type ethereumGetBlockByNumberResponse struct {
	ID     int32
	Result EthereumGetBlockByNumberResult `json:",omitempty"`
	Error  *EthereumError                 `json:",omitempty"`
}
