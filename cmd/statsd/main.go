package main

import (
	"fmt"
	"godemos/pkg/log"
	"godemos/pkg/statsd"
	"os"
	"os/signal"
)

func main() {
	log.Init("[statsd]")

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	errCh := make(chan error)
	go func() {
		errCh <- run()
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

func run() (err error) {
	cli, err := statsd.NewClient("127.0.0.1:8125", "godemos")
	if err != nil {
		return err
	}
	defer cli.Close()

	var str string
	for {
		log.L.Println("input metric")
		_, err = fmt.Fscanf(os.Stdin, "%s\n", &str)
		if err != nil {
			return err
		}
		cli.Increment("thisisatest_" + str)
		log.L.Printf("'%s' send\n", str)
	}
}
