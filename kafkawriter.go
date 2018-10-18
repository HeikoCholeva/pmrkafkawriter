package main

import (
	"log"
	"sync"
	"strings"
	"github.com/Shopify/sarama"
)

var (
	wg sync.WaitGroup
	successes, errors int
)

func newProducer() sarama.AsyncProducer {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true

        if saslEnabled {
                config.Net.SASL.Enable = saslEnabled
                config.Net.SASL.User = cfg.SASL.Username
                config.Net.SASL.Password = cfg.SASL.Password
        }

	broker := strings.Split(strings.Replace(cfg.Broker, " ", "", -1), ",")
	log.Printf("Broker: %v\n", broker)

	producer, err := sarama.NewAsyncProducer(broker, config)
	if err != nil {
		log.Fatalln("Failed to start Sarama producer:", err)
	}

	wg.Add(1)
        go func() {
                defer wg.Done()
                for range producer.Successes() { successes++ }
        }()

        wg.Add(1)
        go func() {
                defer wg.Done()
                for err := range producer.Errors() {
                        log.Println(err)
                        errors++
                }
        }()

	return producer
}

func writeToKafka(report string) {
	producer.Input() <- &sarama.ProducerMessage{Topic: cfg.Topic, Value: sarama.StringEncoder(report)}
}
