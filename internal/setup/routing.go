package setup

import "net/http"

const (
	blockChainParserGetBlockNumberPath = "/block/number"
	blockChainParserSubscribePath      = "/address/subscribe"
	blockChainParserGetTransaction     = "/address/transaction"
)

var (
	routes = map[string]map[string]struct{}{
		blockChainParserGetBlockNumberPath: {
			http.MethodGet: struct{}{},
		},
		blockChainParserGetTransaction: {
			http.MethodGet: struct{}{},
		},
		blockChainParserSubscribePath: {
			http.MethodPost: struct{}{},
		},
	}
)
