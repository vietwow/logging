package producer

import (
    "fmt"

    "github.com/confluentinc/confluent-kafka-go/kafka"
)

var p *kafka.Producer

func InitKafka(broker string) error {
    var err error
    fmt.Printf("Creating consumer to broker %v\n", broker)
    p, err = kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": broker})
    return err
}

func Produce(topic string, message string) {
    fmt.Printf("=> Created Producer %v\n", p)

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
    p.ProduceChannel() <- &kafka.Message{TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny}, Value: []byte(value)}

    // wait for delivery report goroutine to finish
    _ = <-doneChan

    p.Close()
}