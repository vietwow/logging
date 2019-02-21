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

    fmt.Println("Creating consumer to broker ", broker)

    var err error
    c, err = kafka.NewConsumer(&kafka.ConfigMap{
        "bootstrap.servers":  broker,
        "group.id":           group,
        "session.timeout.ms": 6000,
        "auto.offset.reset":  "earliest"})

    return err
}

func Consume(topic string) {
    fmt.Printf("Created Consumer %v\n", c)
    sigchan := make(chan os.Signal, 1)
    signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

    err := c.SubscribeTopics(strings.Fields(topic), nil)
    if err != nil {
        fmt.Println("Unable to subscribe to topic " + topic + " due to error - " + err.Error())
        os.Exit(1)
    } else {
        fmt.Println("subscribed to topic ", topic)
    }

    run := true

    for run == true {
        select {
        case sig := <-sigchan:
            fmt.Printf("Caught signal %v: terminating\n", sig)
            run = false

        case ev := <-c.Events():
            switch e := ev.(type) {
            case kafka.AssignedPartitions:
                fmt.Fprintf(os.Stderr, "%% %v\n", e)
                c.Assign(e.Partitions)
            case kafka.RevokedPartitions:
                fmt.Fprintf(os.Stderr, "%% %v\n", e)
                c.Unassign()
            case *kafka.Message:
                fmt.Printf("%% Message on %s:\n%s\n",
                    e.TopicPartition, string(e.Value))
            case kafka.PartitionEOF:
                fmt.Printf("%% Reached %v\n", e)
            case kafka.Error:
                // Errors should generally be considered as informational, the client will try to automatically recover
                fmt.Fprintf(os.Stderr, "%% Error: %v\n", e)
            }
        }
    }

    fmt.Printf("Closing consumer\n")
    c.Close()
}