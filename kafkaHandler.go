package main

import (
	"encoding/json"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"io/ioutil"
	"log"
	"os"
	"sync"
)

// Struct for kafkaConfig.json
type Config struct {
	Broker         string   `json:"broker"`
	ProducerTopic  string   `json:"producer_topic"`
	ConsumerTopics []string `json:"consumer_topics"`
	Group          string   `json:"group"`
}

// Return config for setting up Kafka Producer and Consumer
func configKafka() (broker string, producerTopic string, consumerTopics []string, group string) {
	log.Printf("Get config for current request")

	file, _ := os.Open("./kafkaConfig.json")
	defer file.Close()

	b, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	var config Config
	json.Unmarshal(b, &config)

	log.Printf("Kafka Config -> Broker: `%v`, Producer Topic: `%v`, Consumer Topics: `%v`, Group: `%v`",
		config.Broker, config.ProducerTopic, config.ConsumerTopics, config.Group)
	return config.Broker, config.ProducerTopic, config.ConsumerTopics, config.Group
}

func producer(wg *sync.WaitGroup, broker string, topic string, message <-chan string) {
	log.Println("Producer started!")

	// Setting up Consumer (Kafka) config
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": broker,
	})
	if err != nil {
		panic(err)
	}
	defer p.Close()

	// Run go routine for produce available event to Kafka
	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					log.Printf("Produce failed: %v\n", ev.TopicPartition)
				} else {
					log.Printf("Produced message to %v. Message: %s (Header: %s)\n", ev.TopicPartition, ev.Value, ev.Headers)
				}
			}
		}
	}()

	// Setting up kafka message to get ready to be produced
	// message to be produced
	msg := <-message

	// header for the message
	header := map[string]string{
		"key":   "testHeader",
		"value": "headers value are binary",
	}
	p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          []byte(msg),
		Headers:        []kafka.Header{{Key: header["key"], Value: []byte(header["value"])}},
	}, nil)

	// Wait for message deliveries before shutting down
	p.Flush(3 * 1000)
	log.Println("Producer closing!")

	// Done with worker
	wg.Done()
}

func consumer(broker string, topics []string, group string) {
	log.Println("Consumer (Kafka) started!")

	// Setting up Consumer (Kafka) config
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": broker,
		"group.id":          group,
		"auto.offset.reset": "earliest",
	})

	if err != nil {
		panic(err)
	}

	// Subscribe to topics
	c.SubscribeTopics(topics, nil)

	for {
		msg, err := c.ReadMessage(-1)
		if err == nil {
			log.Println("New Request from Kafka")
			log.Printf("Message consumed on %s: %s\n", msg.TopicPartition, string(msg.Value))

			// Send any consumed event to consumerChan
			consumerChan <- string(msg.Value)
		} else {
			log.Printf("Consumer error: %v (%v)\n", err, msg)
		}
	}

	c.Close()
}
