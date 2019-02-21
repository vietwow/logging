package consumer

import (
    "fmt"
    "os"
    "os/signal"
    "syscall"
    "strings"
    "github.com/confluentinc/confluent-kafka-go/kafka"
)

var c *kafka.Consumer

func InitKafka(broker string) error {
    group := os.Getenv("GROUP") // myGroup

    fmt.Printf("Creating consumer to broker %v with group %v\n", broker, group)

    var err error
    c, err = kafka.NewConsumer(&kafka.ConfigMap{
        "bootstrap.servers":               broker,
        "group.id":                        group,
        "auto.offset.reset":    "earliest"})

    return err
}

func Consume(topic string) {
    fmt.Printf("=> Created Consumer %v\n", c)
    sigchan := make(chan os.Signal, 1)
    signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

    err := c.SubscribeTopics(strings.Fields(topic), nil)
    if err != nil {
        fmt.Println("Unable to subscribe to topic " + topic + " due to error - " + err.Error())
        os.Exit(1)
    } else {
        fmt.Println("subscribed to topic ", topic)
    }

    for {
        msg, err := c.ReadMessage(-1)
        if err == nil {
            fmt.Printf("Message on %s: %s\n", msg.TopicPartition, string(msg.Value))
        } else {
            // The client will automatically try to recover from all errors.
            fmt.Printf("Consumer error: %v (%v)\n", err, msg)
        }
    }

    fmt.Printf("Closing consumer\n")
    c.Close()
}