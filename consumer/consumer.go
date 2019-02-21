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

    err := c.SubscribeTopics([]string{topic, "^aRegex.*[Tt]opic"}, nil)
    if err != nil {
        fmt.Println("Unable to subscribe to topic " + topic + " due to error - " + err.Error())
        os.Exit(1)
    } else {
        fmt.Println("subscribed to topic :", topic)
    }

    go func() {
        defer c.Close()

    CONSUMER_FOR:
        for {
            select {
            case <-ctx.Done():
                break CONSUMER_FOR
            default:
                msg, err := c.ReadMessage(-1)
                if err == nil {
                    var consumed ConsumedMessage
                    if err := json.Unmarshal(msg.Value, &consumed); err != nil {
                        fmt.Println(err)
                    }
                    fmt.Printf("success consume. message: %s, timestamp: %d\n", consumed.Message, consumed.Timestamp)
                } else {
                    fmt.Printf("fail consume. reason: %s\n", err.Error())
                }
            }
        }
    }()

    fmt.Println("confluent-kafka-go-example start.")

    <-signals


    // run := true

    // for run == true {
    //     select {
    //     case sig := <-sigchan:
    //         fmt.Printf("Caught signal %v: terminating\n", sig)
    //         run = false

    //     case ev := <-c.Events():
    //         switch e := ev.(type) {
    //         case kafka.AssignedPartitions:
    //             fmt.Fprintf(os.Stderr, "%% %v\n", e)
    //             c.Assign(e.Partitions)
    //         case kafka.RevokedPartitions:
    //             fmt.Fprintf(os.Stderr, "%% %v\n", e)
    //             c.Unassign()
    //         case *kafka.Message:
    //             fmt.Printf("%% Message on %s:\n%s\n",
    //                 e.TopicPartition, string(e.Value))
    //         case kafka.PartitionEOF:
    //             fmt.Printf("%% Reached %v\n", e)
    //         case kafka.Error:
    //             // Errors should generally be considered as informational, the client will try to automatically recover
    //             fmt.Fprintf(os.Stderr, "%% Error: %v\n", e)
    //         }
    //     }
    // }

    fmt.Printf("Closing consumer\n")
    c.Close()
}