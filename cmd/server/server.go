// Command server launch the gRPC AML approval server.
//
// It will start the grpc server on a specific port,
// be aware that the transaction state is not persistent.
//
// Usage:
//
//	server -p 1000
//
// Flags:
//
//	-p int
//	   gRPC server port to listen to
package main

import (
	"context"
	"flag"
	"fmt"
	"jonatak/aml/internal/aml"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	port := flag.Int("p", 8080, "GRpc server port to listen to.")
	flag.Parse()
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	ctxCancel, cancel := context.WithCancel(ctx)

	txChan := make(chan aml.TransactionQuery, aml.MaxBurstGrpcQuery)

	var wg sync.WaitGroup

	errChan := make(chan error, 2)

	go func() {
		if err := <-errChan; err != nil {
			log.Printf("server exited with error: %v", err)
			cancel()
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := aml.StartServer(ctxCancel, *port, txChan); err != nil {
			errChan <- fmt.Errorf("server exited with error: %w", err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		checker := &aml.AMLChecker{}
		if err := checker.StartLoop(ctxCancel, txChan); err != nil {
			errChan <- fmt.Errorf("checker exited with error: %w", err)
		}
	}()

	wg.Wait()
}
