package config

import (
	"log"
	"os"
	"time"
)

type EthereumHttpClient struct {
	Host    string
	Timeout time.Duration
}

func parseEthereumHttpClient() EthereumHttpClient {
	var (
		ok  bool
		err error
	)

	ethereumHttpClientCfg := EthereumHttpClient{}
	ethereumHttpClientCfg.Host, ok = os.LookupEnv("BLOCKCHAIN_PARSER_ETH_HTTP_CLIENT_HOST")
	if !ok {
		log.Fatalf("BLOCKCHAIN_PARSER_ETH_HTTP_CLIENT_HOST is required")
	}
	ethereumHttpClientCfgTimeout, ok := os.LookupEnv("BLOCKCHAIN_PARSER_ETH_HTTP_CLIENT_TIMEOUT")
	if !ok {
		log.Fatalf("BLOCKCHAIN_PARSER_ETH_HTTP_CLIENT_TIMEOUT is required")
	}
	ethereumHttpClientCfg.Timeout, err = time.ParseDuration(ethereumHttpClientCfgTimeout)
	if err != nil {
		log.Fatalf("BLOCKCHAIN_PARSER_ETH_HTTP_CLIENT_TIMEOUT is not duration: %s", err)
	}

	return ethereumHttpClientCfg
}
