package config

type Config struct {
	EthereumHttpClient EthereumHttpClient
	ParserWorker       ParserWorker
	Server             Server
}

func Parse() Config {
	return Config{
		EthereumHttpClient: parseEthereumHttpClient(),
		ParserWorker:       parseParserWorker(),
		Server:             parseServer(),
	}
}
