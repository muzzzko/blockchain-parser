package httpclient

type ethereumRequestBody struct {
	Version string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
	ID      int32         `json:"id"`
}
