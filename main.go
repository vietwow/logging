package main

import (
    "fmt"
    "github.com/confluentinc/confluent-kafka-go/kafka"
    "os"
    "os/signal"
    "syscall"
    "strings"
    "github.com/vietwow/logging/consumer"
)

func main() {
    // brokers := os.Getenv("BROKERS") // localhost:29092
    // topic := os.Getenv("TOPIC") // heroku_logs
    // group := os.Getenv("GROUP") // myGroup

    brokers := os.Getenv("BROKERS") // localhost:29092
    topic := os.Getenv("TOPIC") // heroku_logs
    group := os.Getenv("GROUP") // myGroup

    sigchan := make(chan os.Signal, 1)
    signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

    c, err := kafka.NewConsumer(&kafka.ConfigMap{
        "bootstrap.servers":               brokers,
        "group.id":                        group,
        "session.timeout.ms":              6000,
        "go.events.channel.enable":        true,
        "go.application.rebalance.enable": true,
        // Enable generation of PartitionEOF when the
        // end of a partition is reached.
        "enable.partition.eof": true,
        "auto.offset.reset":    "earliest"})

    if err != nil {
        fmt.Fprintf(os.Stderr, "Failed to create consumer: %s\n", err)
        os.Exit(1)
    }

    fmt.Printf("Created Consumer %v\n", c)

    err = c.SubscribeTopics(strings.Fields(topic), nil)

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
