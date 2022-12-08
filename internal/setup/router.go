package setup

import "net/http"

const (
	blockChainParserGetBlockNumberPath = "/block/number"
	blockChainParserSubscribePath      = "/address/subscribe"
	blockChainParserGetTransaction     = "/address/transaction"
)

var (
	routes = map[string]map[string]struct{}{
		http.MethodGet: {
			blockChainParserGetBlockNumberPath: struct{}{},
			blockChainParserGetTransaction:     struct{}{},
		},
		http.MethodPost: {
			blockChainParserSubscribePath: struct{}{},
		},
	}
)
