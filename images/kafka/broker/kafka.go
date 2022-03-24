package broker

import (
	"context"
	"os"

	singe "github.com/Adversary-Informed-Defense/singe/pkg/singe"
	kafka "github.com/confluentinc/confluent-kafka-go/kafka"
	sigma "github.com/markuskont/go-sigma-rule-engine/pkg/sigma/v2"
	logrus "github.com/sirupsen/logrus"
)

type OutputMessage struct {
	Event  sigma.Event
	Result sigma.Results
}

var brokerAddresses []string = []string{string(os.Getenv("KAFKA_SERVICE_HOST") + ":" + os.Getenv("KAFKA_SERVICE_PORT"))}

// CreateKafkaTopics creates Kafka topics from a list of topic names and returns the results
func CreateKafkaTopics(topics []string) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	admin, err := kafka.NewAdminClient(&kafka.ConfigMap{
		"bootstrap.servers": brokerAddresses[0],
	})
	if err != nil {
		logrus.Errorf("Error Creating Topics: %s", err)
	}
	topicSpecs := []kafka.TopicSpecification{}
	for _, topic := range topics {
		topicSpecs = append(
			topicSpecs,
			kafka.TopicSpecification{
				Topic:             topic,
				NumPartitions:     10, // TODO: Will probably be determined by config
				ReplicationFactor: 1,
			})
	}
	results, err := admin.CreateTopics(ctx, topicSpecs, kafka.SetAdminOperationTimeout(0))
	if err != nil {
		logrus.Errorf("Error Creating Topics: %s", err)
	}
	logrus.Infof("Topics Created: %s", results)
}

// ProduceKafkaMsg produces every byte array that is sent to the msgs channel as a Kafka message
func ProduceKafkaMsg(msgs chan []byte, producer *kafka.Producer, topic *string, errs chan error) {
	for msg := range msgs {
		// Produce messages to Kafka topic
		if err := producer.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{
				Topic:     topic,
				Partition: kafka.PartitionAny,
			},
			Value: msg,
		}, nil); err != nil {
			// Send errors to error channel
			errs <- err
		}
	}
}

// CreateConsumer loops the creation of a Kafka consumer until success
func CreateConsumer(topics []string, groupId string) *kafka.Consumer {
	var err error
	var consumer *kafka.Consumer
	for consumer == nil {
		consumer, err = kafka.NewConsumer(&kafka.ConfigMap{
			"bootstrap.servers": brokerAddresses[0],
			"group.id":          groupId,
			"auto.offset.reset": "latest",
		})
		if err != nil {
			logrus.Debugf("Error creating consumer: %s", err)
			continue
		}
	}
	logrus.Infof("Created Kafka consumer")
	Subscribe(consumer, topics)
	return consumer
}

// Subscribe subscribes consumer to Kafka topics
func Subscribe(consumer *kafka.Consumer, topics []string) {
	var err error
	if err = consumer.SubscribeTopics(topics, nil); err != nil {
		logrus.Errorf("Error subscribing to topics: %s", err)
	}
	logrus.Infof("Consumer subscribed to topics")
}

// CreateProducer loops the creation of a Kafka producer until success
func CreateProducer() *kafka.Producer {
	var err error
	var producer *kafka.Producer
	for producer == nil {
		producer, err = kafka.NewProducer(&kafka.ConfigMap{
			"bootstrap.servers": brokerAddresses[0],
		})
		if err != nil {
			logrus.Debugf("Error creating producer: %s", err)
			continue
		}
	}
	logrus.Infof("Created producer")
	return producer
}

var engine singe.SigmaEngine

func init() {
	engine = singe.CreateEngine(string(os.Getenv("RULESET_PATH")))
}

// HandleEvent processes the log event contained in the Kafka message
// The event runs through the Sigma engine and if at least one match occurs, an ouput struct containing the event data and the Sigma engine results is returned
func HandleEvent(msg *kafka.Message) ([]byte, bool, error) {
	return engine.Match(string(msg.Value), *msg.TopicPartition.Topic)
}
