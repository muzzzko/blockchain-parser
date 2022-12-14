package setup

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"blockchain-parser/config"
	"blockchain-parser/internal/constant"
	"blockchain-parser/internal/entity"
	"blockchain-parser/internal/infrastructure/handler"
	"blockchain-parser/internal/infrastructure/httpclient"
	lockerpkg "blockchain-parser/internal/infrastructure/locker"
	"blockchain-parser/internal/infrastructure/repository"
	"blockchain-parser/internal/service"
	"blockchain-parser/tools/job"
)

type Server struct {
	starts []func(ctx context.Context)
	stops  []func()
}

func (s *Server) Configure() {
	cfg := config.Parse()

	//-------------------
	// repositories
	//-------------------

	txnRepo := repository.NewInMemTransaction()
	subscriberRepo := repository.NewInMemSubscriber()
	blockRepo := repository.NewInMemBlock()

	//-------------------
	// http clients
	//-------------------

	ethereumClient := httpclient.NewEthereum(cfg.EthereumHttpClient)

	//-------------------
	// locker
	//-------------------

	locker := lockerpkg.NewInMem()

	//-------------------
	// services
	//-------------------

	parser := service.NewParser(txnRepo, subscriberRepo, blockRepo)
	parserWorker := service.NewParserWorker(txnRepo, subscriberRepo, blockRepo, ethereumClient, locker)

	//-------------------
	// handlers
	//-------------------

	BlockChainParserHandler := handler.NewBlockChainParser(parser)

	mux := http.NewServeMux()
	mux.HandleFunc(blockChainParserGetBlockNumberPath, BlockChainParserHandler.GetCurrentBlock)
	mux.HandleFunc(blockChainParserSubscribePath, BlockChainParserHandler.Subscribe)
	mux.HandleFunc(blockChainParserGetTransaction, BlockChainParserHandler.GetTransactions)

	//-------------------
	// setup server
	//-------------------

	srv := http.Server{
		Addr:    fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port),
		Handler: methodCheckMiddleware(panicRecoveryMiddleware(contentTypeMiddleware(mux))),
	}

	startServer := func(_ context.Context) {
		go func() {
			log.Printf("server started at: %s\n", srv.Addr)

			if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
				log.Printf("fail stop server: %s\n", err)
			}
		}()
	}
	stopServer := func() {
		if err := srv.Shutdown(context.Background()); err != nil {
			log.Printf("fail shutdown: %s", err)
		}
	}

	s.starts = append(s.starts, startServer)
	s.stops = append(s.stops, stopServer)

	//-------------------
	// setup initial state
	//-------------------

	setupStartBlockNumber(ethereumClient, blockRepo, cfg.ParserWorker)
	subscribePredefinedAddress(subscriberRepo, cfg.ParserWorker)

	s.createJobs(cfg, parserWorker)
}

func (s *Server) Start(ctx context.Context) {
	for _, start := range s.starts {
		start(ctx)
	}
}

func (s *Server) Stop() {
	for _, stop := range s.stops {
		stop()
	}
}

func (s *Server) createJobs(cfg config.Config, parserWorker *service.ParserWorker) {
	jobs := job.Jobs{}
	for i := 0; i < cfg.ParserWorker.CountWorkers; i++ {
		jobs.Add(job.NewJob(
			parserWorker.Run,
			constant.ParserWorkerJobName,
			cfg.ParserWorker.Interval,
		))
	}

	s.starts = append(s.starts, jobs.Start)
	s.stops = append(s.stops, jobs.Stop)
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

func subscribePredefinedAddress(subscriberRepo *repository.InMemSubscriber, cfg config.ParserWorker) {
	for _, address := range cfg.PredefinedAddresses {
		subscriber := entity.Subscriber{
			Address: address,
		}
		if err := subscriberRepo.Save(context.Background(), subscriber); err != nil {
			log.Fatalf("fail subscribe predefined address: %s", err)
		}
	}
}
