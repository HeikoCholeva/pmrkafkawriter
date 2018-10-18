package main

import (
	"io"
	"os"
	"log"
	"strconv"
	"github.com/Shopify/sarama"
)

var (
	cfg Config
	saslEnabled bool
	tlsEnabled bool
	producer sarama.AsyncProducer
)

func main() {
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
		StartServer(port, cfg.WebServer.Path, true, cfg.WebServer.TLS.Cert, cfg.WebServer.TLS.Key)
	} else {
		StartServer(port, cfg.WebServer.Path, false, "", "")
	}
}
