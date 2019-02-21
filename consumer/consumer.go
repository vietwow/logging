package consumer

import (
    "fmt"
    "os"
    "os/signal"
    "syscall"
    "strings"
    "github.com/confluentinc/confluent-kafka-go/kafka"
)

var consumer *kafka.Consumer

func InitKafka() error {
    broker := os.Getenv("BROKERS") // localhost:29092
    topic := os.Getenv("TOPIC") // heroku_logs
    group := os.Getenv("GROUP") // myGroup

    var err error
    c, err := kafka.NewConsumer(&kafka.ConfigMap{
        "bootstrap.servers":  broker,
        "group.id":           group,
        "session.timeout.ms": 6000,
        "auto.offset.reset":  "earliest"})

    return err
}

func Consume(topics string, message string) error {
    sigchan := make(chan os.Signal, 1)
    signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

    err := c.SubscribeTopics(strings.Fields(topic), nil)

    run := true

    for run == true {
        select {
        case sig := <-sigchan:
            fmt.Printf("Caught signal %v: terminating\n", sig)
            run = false
        default:
            ev := c.Poll(100)
            if ev == nil {
                continue
            }

            switch e := ev.(type) {
            case *kafka.Message:
                fmt.Printf("%% Message on %s:\n%s\n",
                    e.TopicPartition, string(e.Value))
                if e.Headers != nil {
                    fmt.Printf("%% Headers: %v\n", e.Headers)
                }
            case kafka.Error:
                // Errors should generally be considered as informational, the client will try to automatically recover
                fmt.Fprintf(os.Stderr, "%% Error: %v\n", e)
            default:
                fmt.Printf("Ignored %v\n", e)
            }
        }
    }

    fmt.Printf("Closing consumer\n")
    c.Close()
}