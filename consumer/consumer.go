package consumer

import (
    "fmt"
    "os"
    "os/signal"
    // "strings"
    // "context"
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
    signal.Notify(sigchan, os.Interrupt)

    err := c.SubscribeTopics([]string{topic, nil)
    if err != nil {
        fmt.Println("Unable to subscribe to topic " + topic + " due to error - " + err.Error())
        os.Exit(1)
    } else {
        fmt.Println("subscribed to topic :", topic)
    }

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