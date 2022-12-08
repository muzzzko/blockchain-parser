package config

import (
	"log"
	"os"
)

type Server struct {
	Host string
	Port string
}

func parseServer() Server {
	var (
		ok bool
	)

	serverCfg := Server{}
	serverCfg.Host, ok = os.LookupEnv("BLOCKCHAIN_PARSER_SERVER_HOST")
	if !ok {
		log.Fatalf("BLOCKCHAIN_PARSER_SERVER_HOST is required")
	}

	serverCfg.Port, ok = os.LookupEnv("BLOCKCHAIN_PARSER_SERVER_PORT")
	if !ok {
		log.Fatalf("BLOCKCHAIN_PARSER_SERVER_PORT is required")
	}

	return serverCfg
}
