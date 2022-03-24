package main

import (
	"os"

	broker "github.com/Adversary-Informed-Defense/k8s-go-sigma-streamer/images/kafka/broker"
	tools "github.com/Adversary-Informed-Defense/singe/pkg/singe/tools"

	logrus "github.com/sirupsen/logrus"
)

func main() {

	// Define topics
	outputTopics := []string{os.Getenv("OUTPUT_TOPIC")}
	consumeTopics := tools.GetListEnvVar("INGEST_TOPICS")
	allTopics := append(outputTopics, consumeTopics...)
	broker.CreateKafkaTopics(allTopics)

	consumer := broker.CreateConsumer(consumeTopics, "sigma-engine")
	defer consumer.Close()

	producer := broker.CreateProducer()
	defer producer.Close()

	// Make output alerts channel and error channel
	outputMsgs := make(chan []byte)
	defer close(outputMsgs)
	outputErrs := make(chan error)
	defer close(outputErrs)

	// Produce output messages to Kafka
	go broker.ProduceKafkaMsg(outputMsgs, producer, &outputTopics[0], outputErrs)

	// Handle errors handling events
	// TODO: Create error struct containing error context
	// go ErrorHandler(outputErrs)

	logrus.Infof("Reading Messages")
	for {
		// TODO: Add condition to stop waiting on logs and exit gracefully
		// Consume logs via kafka
		msg, err := consumer.ReadMessage(-1) // No timeout waiting for kafka message
		if err != nil {
			logrus.Infof("Error Reading Message: %s", err)
			continue
		}

		// Handle log event
		// Goroutines are ~3x faster on test data
		go func(out chan []byte, errs chan error) {
			output, ok, err := broker.HandleEvent(msg)
			switch {
			case err != nil:
				// Send errors to error channel for handling
				errs <- err
			case ok:
				// If match occurred, send event and match results data alerts output channel
				out <- output
			}
		}(outputMsgs, outputErrs)
	}
}
