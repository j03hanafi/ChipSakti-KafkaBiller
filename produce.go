package main

import (
	"errors"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"log"
	"sync"
)

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

func doProducer(broker string, topic string, iso string) (err error) {

	log.Printf("Starting producer\n")

	// Setup producer config
	cm := kafka.ConfigMap{
		"bootstrap.servers": broker,
	}

	// Create new producer
	if p, err := kafka.NewProducer(&cm); err != nil {
		if ke, ok := err.(kafka.Error); ok == true {
			switch ec := ke.Code(); ec {
			case kafka.ErrInvalidArg:
				return fmt.Errorf("Can't create producer because wrong configuration (code: %d)!\n\t%v\n\nTo see the configuration options, refer to https://github.com/edenhill/librdkafka/blob/master/CONFIGURATION.md\n", ec, err)
			default:
				return fmt.Errorf("Can't create producer (code: %d)!\n\t%v\n", ec, err)
			}
		} else {
			// Not Kafka Error occurs
			return fmt.Errorf("Can't create producer because generic error! \n\t%v\n", err)
		}
	} else {

		// For signalling termination from main to go-routine
		termChan := make(chan bool, 1)
		// For signalling that termination is done from go-routine to main
		doneChan := make(chan bool)
		// For capturing errors from the go-routine
		errorChan := make(chan string, 8)

		// build the messages
		m := kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic},
			Value:          []byte(iso),
			Headers:        []kafka.Header{{Key: "testHeader", Value: []byte("headers value are binary")}},
		}

		// handle any event that we get
		go func() {
			doTerm := false
			for !doTerm {
				select {
				case ev := <-p.Events():
					switch ev.(type) {

					// Look at the type of Event we've received
					case *kafka.Message:
						// delivery report
						km := ev.(*kafka.Message)
						if km.TopicPartition.Error != nil {
							errorChan <- fmt.Sprintf("Failed to send message '%v' \n\tto topic '%v'\n\tError: %v\n",
								string(km.Value),
								*km.TopicPartition.Topic,
								km.TopicPartition.Error)
						} else {
							log.Printf("Message '%v' \n\tdelivered to topic '%v' (partition %d at offset %d)\n",
								string(km.Value),
								*km.TopicPartition.Topic,
								km.TopicPartition.Partition,
								km.TopicPartition.Offset)
						}

					case kafka.Error:
						// Kafka error
						em := ev.(kafka.Error)
						errorChan <- fmt.Sprintf("Caught an Error:\n\t%v\n", em)

					default:
						// Not anything we expected occur
						errorChan <- fmt.Sprintf("Got an event but not an Error \n\t%v\n", ev)
					}

				case <-termChan:
					doTerm = true
				}
			}

			close(errorChan)
			close(doneChan)
		}()

		// produce the message
		if e := p.Produce(&m, nil); e != nil {
			errorChan <- fmt.Sprintf("Error producing the message! %v\n", e)
		}

		// flush the Producer's queue
		t := 5000
		if r := p.Flush(t); r > 0 {
			errorChan <- fmt.Sprintf("Failed to flush all messages after %d seconds. %d message(s) remain\n", t, r)
		} else {
			log.Printf("All messages flused from queue\n")
		}

		// signal termination to go-routine
		termChan <- true
		// wait for go-routine to terminate
		<-doneChan

		// check error when done
		done := false
		var err string
		for !done {
			if t, o := <-errorChan; o == false {
				done = true
			} else {
				err += t
			}
		}

		if len(err) > 0 {
			// if error not nil
			log.Printf("returning an error\n")
			return errors.New(err)
		}
		log.Printf("Closing producer...\n")
		p.Close()

	}

	return nil

}
