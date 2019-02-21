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

    // Optional delivery channel, if not specified the Producer object's
    // .Events channel is used.
    deliveryChan := make(chan kafka.Event)

    value := "Hello Go!"
    err := p.Produce(&kafka.Message{
        TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
        Value:          []byte(value),
        Headers:        []kafka.Header{{Key: "myTestHeader", Value: []byte("header values are binary")}},
    }, deliveryChan)

    e := <-deliveryChan
    m := e.(*kafka.Message)

    if m.TopicPartition.Error != nil {
        fmt.Printf("Delivery failed: %v\n", m.TopicPartition.Error)
    } else {
        fmt.Printf("Delivered message to topic %s [%d] at offset %v\n",
            *m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
    }

    close(deliveryChan)
}