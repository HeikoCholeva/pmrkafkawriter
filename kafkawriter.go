package main

import (
	"log"
	"strings"
	"sync"
	"time"

	"github.com/Shopify/sarama"
)

var (
	wg                         sync.WaitGroup
	successes, errors, resendc int
)

func stopProducer(p sarama.AsyncProducer) {
	p.AsyncClose()
	swg.Done()
}

func newProducer() sarama.AsyncProducer {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.Flush.Messages = 1
	config.Producer.Flush.MaxMessages = 5

	if cfg.SASL.Username != "" && cfg.SASL.Password != "" {
		config.Net.SASL.Enable = true
		config.Net.SASL.User = cfg.SASL.Username
		config.Net.SASL.Password = cfg.SASL.Password
	}

	broker := strings.Split(strings.Replace(cfg.Broker, " ", "", -1), ",")
	log.Printf("Broker: %v\n", broker)
	config.ClientID = "pmrkafkawriter"
	producer, err := sarama.NewAsyncProducer(broker, config)
	if err != nil {
		log.Fatalln("Failed to start Sarama producer:", err)
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		for range producer.Successes() {
			successes++
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for perr := range producer.Errors() {
			errors++
			log.Printf("Error @producer.Errors: %v", perr)
			resend(perr.Msg)
		}
	}()

	return producer
}

func writeToKafka(report string) {
	producer.Input() <- &sarama.ProducerMessage{Topic: cfg.Topic, Value: sarama.StringEncoder(report)}
}

func resend(msg *sarama.ProducerMessage) {
	time.Sleep(3 * time.Second)
	producer.Input() <- msg
	resendc++
}
