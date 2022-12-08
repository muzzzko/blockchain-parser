package main

import (
	"blockchain-parser/config"
	"blockchain-parser/internal/constant"
	"blockchain-parser/internal/entity"
	"blockchain-parser/internal/infrastructure/httpclient"
	lockerpkg "blockchain-parser/internal/infrastructure/locker"
	"blockchain-parser/internal/infrastructure/repository"
	"blockchain-parser/internal/service"
	"blockchain-parser/tools/job"
	"context"
	"log"
	"os"
	"os/signal"
	"time"
)

func main() {
	cfg := config.Parse()

	txnRepo := repository.NewInMemTransaction()
	subscriberRepo := repository.NewInMemSubscriber()
	blockRepo := repository.NewInMemBlock()

	ethereumClient := httpclient.NewEthereum(cfg.EthereumHttpClient)

	locker := lockerpkg.NewInMem()

	setupStartBlockNumber(ethereumClient, blockRepo, cfg.ParserWorker)

	service.NewParser(txnRepo, subscriberRepo, blockRepo)
	parserWorker := service.NewParserWorker(txnRepo, subscriberRepo, blockRepo, ethereumClient, locker)

	jobs := job.Jobs{}
	for i := 0; i < cfg.ParserWorker.CountWorkers; i++ {
		jobs.Add(job.NewJob(
			parserWorker.Run,
			constant.ParserWorkerJobName,
			cfg.ParserWorker.Interval,
		))
	}

	jobs.Start(context.Background())

	sigch := make(chan os.Signal, 1)
	signal.Notify(sigch, os.Interrupt, os.Kill)

	<-sigch

	jobs.Stop()
}

func setupStartBlockNumber(
	ethereumClient *httpclient.Ethereum,
	blockRepo *repository.InMemBlock,
	cfg config.ParserWorker,
) {
	block := entity.Block{
		Number:    int(cfg.StartBlockNumber),
		Status:    constant.BlockStatusParsed,
		UpdatedAt: time.Now(),
	}

	if block.Number == -1 {
		var err error
		block.Number, err = ethereumClient.GetBlockNumber(context.Background())
		if err != nil {
			log.Fatalf("fail get block number: %s", err)
		}
	}

	if err := blockRepo.Upsert(context.Background(), block); err != nil {
		log.Fatalf("fail save block: %s", err)
	}
}
