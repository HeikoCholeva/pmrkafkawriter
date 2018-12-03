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
	configFlag  = flag.String("config", "", "path to the config file")
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
	if *configFlag == "" {
		log.Fatal("missing mandatory config path parameter")
		os.Exit(1)
	}
	err := cfg.FromFile(*configFlag)
	if err != nil {
		panic(err)
	}
	go signalHandler()
	logfile := fmt.Sprintf("%s/%s", cfg.Log.Path, cfg.Log.Filename)
	f, err := os.OpenFile(logfile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	defer f.Close()
	var mw io.Writer
	if cfg.Log.STDOUT {
		mw = io.MultiWriter(os.Stdout, f)
	} else {
		mw = io.MultiWriter(f)
	}
	log.SetOutput(mw)

	listen := fmt.Sprintf("%s:%d", cfg.WebServer.Listen, cfg.WebServer.Port)

	producer = newProducer()

	ba := false
	if cfg.WebServer.BasicAuthFile != "" {
		ba = true
	}

	if cfg.WebServer.TLS.Cert != "" && cfg.WebServer.TLS.Key != "" {
		startServer(listen, cfg.WebServer.Path, true, cfg.WebServer.TLS.Cert, cfg.WebServer.TLS.Key, ba)
	} else {
		startServer(listen, cfg.WebServer.Path, false, "", "", ba)
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
