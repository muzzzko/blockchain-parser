package config

type Config struct {
	EthereumHttpClient EthereumHttpClient
	ParserWorker       ParserWorker
}

func Parse() Config {
	return Config{
		EthereumHttpClient: parseEthereumHttpClient(),
		ParserWorker:       parseParserWorker(),
	}
}
