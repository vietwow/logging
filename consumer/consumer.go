package consumer

import (
    "fmt"
    "os"
    "os/signal"
    "syscall"
    // "strings"
    "github.com/confluentinc/confluent-kafka-go/kafka"
)

var c *kafka.Consumer

func InitKafka(broker string) error {
    group := os.Getenv("GROUP") // myGroup

    fmt.Printf("Creating consumer to broker %v with group %v\n", broker, group)

    var err error
    c, err = kafka.NewConsumer(&kafka.ConfigMap{
        "bootstrap.servers":  broker,
        "group.id":           group,
        "auto.offset.reset":  "earliest"})

    return err
}

func Consume(topic string) {
    fmt.Printf("=> Created Consumer %v\n", c)
    sigchan := make(chan os.Signal, 1)
    signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

    err := c.SubscribeTopics([]string{"heroku_logs", "^aRegex.*[Tt]opic"}, nil)
    if err != nil {
        fmt.Println("Unable to subscribe to topic " + topic + " due to error - " + err.Error())
        os.Exit(1)
    } else {
        fmt.Println("subscribed to topic :", topic)
    }

    fmt.Printf("Closing consumer\n")
    c.Close()
}