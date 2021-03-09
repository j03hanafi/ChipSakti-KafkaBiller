package main

import (
	"log"
	"net/http"
	"os"
	"sync"
)

var (
	billerChan   = make(chan string) // channel for send-receive data from-to `Biller`
	producerChan = make(chan string) // channel for receive data from channelChan and send data to `Producer (Kafka)`
	consumerChan = make(chan string) // channel for receive data from `Consumer (Kafka)` and send data to channelChan
)

func main() {
	// Setting up log file
	// set permission to read/write log file
	// read/write to existing log file, if there is none it will create new log file
	file, err := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal("Found error in log ", err)
	}
	log.SetOutput(file)

	// ChannelKafka started
	log.Println("Service Started!")

	// Setting up HTTP Listener and Handler
	// router will handle any request at any endpoint available in server()
	router := server()
	go func() {
		// listen to specific address and handler
		address := "localhost:6020"
		err := http.ListenAndServe(address, router)
		log.Println("Server started at", address)
		if err != nil {
			log.Fatal(err.Error())
		}
	}()

	// WaitGroup for make sure that an event is produced to Kafka
	var wg sync.WaitGroup

	// Get config for Kafka Producer and Consumer
	broker, producerTopic, consumerTopics, groups := configKafka()

	// Run Consumer (Kafka)
	go consumer(broker, consumerTopics, groups)

	// Run Goroutine for request-response data from-to `Biller`
	go requestHandler()

	// loop for checking if there is any new response from `Biller` that has been sent to channelChan
	for {
		select {
		// execute if there is a new response in channelChan
		case newResponse := <-billerChan:
			log.Println("New response from `Biller` is ready to produce to Kafka")

			// Add WaitGroup counter to wait until Producer (Kafka) finish producing an event
			wg.Add(1)

			// Run Producer (Kafka)
			go producer(&wg, broker, producerTopic, producerChan)

			// Send new response to producerChan, then produce the new response to Kafka
			producerChan <- newResponse

			// Waiting until Producer (Kafka) finish producing event
			wg.Wait()

		// keep looping if there is none new response
		default:
			continue
		}
	}
}
