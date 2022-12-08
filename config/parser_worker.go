package config

import (
	"log"
	"os"
	"strconv"
	"time"
)

type ParserWorker struct {
	CountWorkers     int
	Interval         time.Duration
	StartBlockNumber int64
}

func parseParserWorker() ParserWorker {
	var (
		ok  bool
		err error
	)

	parserWorkerCfg := ParserWorker{}
	parserWorkerCfgCountWorkers, ok := os.LookupEnv("BLOCKCHAIN_PARSER_PARSER_WORKER_COUNT_WORKERS")
	if !ok {
		log.Fatalf("BLOCKCHAIN_PARSER_PARSER_WORKER_COUNT_WORKERS is required")
	}
	parserWorkerCfg.CountWorkers, err = strconv.Atoi(parserWorkerCfgCountWorkers)
	if err != nil {
		log.Fatalf("BLOCKCHAIN_PARSER_PARSER_WORKER_COUNT_WORKERS is not integer: %s", err)
	}

	parserWorkerCfgInterval, ok := os.LookupEnv("BLOCKCHAIN_PARSER_PARSER_WORKER_INTERVAL")
	if !ok {
		log.Fatalf("BLOCKCHAIN_PARSER_PARSER_WORKER_INTERVAL is required")
	}
	parserWorkerCfg.Interval, err = time.ParseDuration(parserWorkerCfgInterval)
	if err != nil {
		log.Fatalf("BLOCKCHAIN_PARSER_PARSER_WORKER_INTERVAL is not duration: %s", err)
	}

	parserWorkerCfgStartBlockNumber, ok := os.LookupEnv("BLOCKCHAIN_PARSER_PARSER_WORKER_START_BLOCK_NUMBER")
	if ok {
		parserWorkerCfg.StartBlockNumber, err = strconv.ParseInt(parserWorkerCfgStartBlockNumber, 0, 64)
		if err != nil {
			log.Fatalf("BLOCKCHAIN_PARSER_PARSER_WORKER_START_BLOCK_NUMBER is not hex: %s", err)
		}
	} else {
		parserWorkerCfg.StartBlockNumber = -1
	}

	return parserWorkerCfg
}
