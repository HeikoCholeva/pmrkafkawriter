package main

import (
	"io"
	"os"
	"log"
	"sync"
	"strconv"
	"os/signal"
	"github.com/Shopify/sarama"
)

var (
	cfg Config
	saslEnabled bool
	tlsEnabled bool
	producer sarama.AsyncProducer
	swg sync.WaitGroup
)

func main() {
	go signalHandler()

	f, err := os.OpenFile("server.log", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
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

	saslEnabled, _ = strconv.ParseBool(cfg.SASL.Enabled)
	tlsEnabled, _ = strconv.ParseBool(cfg.WebServer.TLS.Enabled)
	port := ":" + string(cfg.WebServer.Port)

	producer = newProducer()

	if tlsEnabled {
		startServer(port, cfg.WebServer.Path, true, cfg.WebServer.TLS.Cert, cfg.WebServer.TLS.Key)
	} else {
		startServer(port, cfg.WebServer.Path, false, "", "")
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
