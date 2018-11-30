package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"sync"

	"github.com/Shopify/sarama"
)

var (
	cfg         Config
	producer    sarama.AsyncProducer
	swg         sync.WaitGroup
	versionFlag = flag.Bool("version", false, "print the version of the program")
)
var githash, shorthash, builddate, buildtime string

func main() {
	flag.Parse()
	if *versionFlag {
		fmt.Fprintln(os.Stderr, `pmrkafkawriter`)
		fmt.Fprintf(os.Stderr, "Version  : %s-%s\n", builddate,
			shorthash)
		fmt.Fprintf(os.Stderr, "Git Hash : %s\n", githash)
		fmt.Fprintf(os.Stderr, "Timestamp: %s\n", buildtime)
		os.Exit(0)
	}
	go signalHandler()

	f, err := os.OpenFile("server.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	defer f.Close()
	mw := io.MultiWriter(os.Stdout, f)
	log.SetOutput(mw)

	err = cfg.FromFile("server.properties")
	if err != nil {
		panic(err)
	}

	port := ":" + string(cfg.WebServer.Port)

	producer = newProducer()

	ba := false
	if cfg.WebServer.BasicAuthFile != "" {
		ba = true
	}

	if cfg.WebServer.TLS.Cert != "" && cfg.WebServer.TLS.Key != "" {
		startServer(port, cfg.WebServer.Path, true, cfg.WebServer.TLS.Cert, cfg.WebServer.TLS.Key, ba)
	} else {
		startServer(port, cfg.WebServer.Path, false, "", "", ba)
	}
}

func signalHandler() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

SignalLoop:
	for {
		select {
		case <-signals:
			log.Println("Shutting down server...")
			swg.Add(1)
			stopProducer(producer)
			swg.Wait()
			log.Printf("Errors %v, ResendCount %v", errors, resendc)
			stopServer()
			break SignalLoop
		}
	}
}
