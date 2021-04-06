package main

import "testing"

func TestConfigKafka(t *testing.T) {
	expectedBroker := "localhost:9092"
	expectedProducerTopic := "goroutine-biller"
	expectedConsumerTopics := []string{"goroutine-channel"}
	expectedGroup := "test-go"

	broker, producerTopic, consumerTopics, groups := configKafka()

	if broker != expectedBroker {
		t.Errorf("broker value at configKafka() failed. Expected: %v. Got: %v", expectedBroker, broker)
	} else {
		t.Log("broker at configKafka() success")
	}

	if producerTopic != expectedProducerTopic {
		t.Errorf("producerTopic value at configKafka() failed. Expected: %v. Got: %v", expectedProducerTopic, producerTopic)
	} else {
		t.Log("producerTopic at configKafka() success")
	}

	if groups != expectedGroup {
		t.Errorf("groups value at configKafka() failed. Expected: %v. Got: %v", expectedGroup, groups)
	} else {
		t.Log("groups at configKafka() success")
	}

	if len(consumerTopics) != len(expectedConsumerTopics) {
		t.Errorf("consumerTopics at configKafka() failed. Expected length: %v. Got length: %v", len(expectedConsumerTopics), len(consumerTopics))
	} else {
		isPassed := true
	test:
		for index, _ := range consumerTopics {
			if consumerTopics[index] != expectedConsumerTopics[index] {
				t.Errorf("consumerTopics value at configKafka() failed. Expected: %v. Got: %v", expectedConsumerTopics[index], consumerTopics[index])
				isPassed = false
				break test
			}
		}
		if isPassed {
			t.Log("consumerTopics at configKafka() success")
		}
	}
}
