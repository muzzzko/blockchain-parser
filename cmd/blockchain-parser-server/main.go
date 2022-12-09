package main

import (
	"blockchain-parser/internal/setup"
	"context"
	"os"
	"os/signal"
)

func main() {
	srv := setup.Server{}
	srv.Configure()

	srv.Start(context.Background())

	sigch := make(chan os.Signal, 1)
	signal.Notify(sigch, os.Interrupt, os.Kill)

	<-sigch

	srv.Stop()
}
