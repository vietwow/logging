package producer

import (
    "fmt"
    "os"

    "github.com/confluentinc/confluent-kafka-go/kafka"
)

var p *kafka.Producer

func InitKafka() error {
    brokers := os.Getenv("BROKERS") // localhost:29092

    var err error
    p, err = kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": brokers})
    return err
}

func Produce(topics string, message string) error {
	doneChan := make(chan bool)

	go func() {
		defer close(doneChan)
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				m := ev
				if m.TopicPartition.Error != nil {
					fmt.Printf("Delivery failed: %v\n", m.TopicPartition.Error)
				} else {
					fmt.Printf("Delivered message to topic %s [%d] at offset %v\n",
						*m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
				}
				return

			default:
				fmt.Printf("Ignored event: %s\n", ev)
			}
		}
	}()

	value := "Hello Go!"
	p.ProduceChannel() <- &kafka.Message{TopicPartition: kafka.TopicPartition{Topic: &topics, Partition: kafka.PartitionAny}, Value: []byte(value)}

	// wait for delivery report goroutine to finish
	_ = <-doneChan

	p.Close()
}