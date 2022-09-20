package main

import (
	"context"
	"godemos/pkg/log"
	"godemos/pkg/otel"
	"os"
	"os/signal"
)

func main() {
	log.Init("[otel]")

	ctx := context.Background()

	shutdownTraceProvider, err := otel.InstallFilePipeline(ctx, "otel_example")
	if err != nil {
		log.L.Fatalf("failed to initialize stdout export pipeline: %v", err)
	}
	defer shutdownTraceProvider()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	errCh := make(chan error)
	app := NewApp(os.Stdin, log.L)
	go func() {
		errCh <- app.Run(ctx)
	}()

	select {
	case <-sigCh:
		log.L.Println("\ngoodbye")
		return
	case err := <-errCh:
		if err != nil {
			log.L.Fatal(err)
		}
	}
}
