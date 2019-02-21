package main

import (
    "context"
    "encoding/json"
    "flag"
    "fmt"
    "os"
    "os/signal"
    "time"

    "github.com/confluentinc/confluent-kafka-go/kafka"
)

var (
    bootstrapServers = flag.String("bootstrapServers", "kafka:29092", "kafka address")
)

// SendMessage 送信メッセージ
type SendMessage struct {
    Message   string `json:"message"`
    Timestamp int64  `json:"timestamp"`
}

// ConsumedMessage 受信メッセージ
type ConsumedMessage struct {
    Message   string `json:"message"`
    Timestamp int64  `json:"timestamp"`
}

func main() {
    flag.Parse()

    if *bootstrapServers == "" {
        flag.PrintDefaults()
        os.Exit(1)
    }

    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    signals := make(chan os.Signal, 1)
    signal.Notify(signals, os.Interrupt)

    topic := "test.D"

    // Consumer
    c, err := kafka.NewConsumer(&kafka.ConfigMap{
        "bootstrap.servers": *bootstrapServers,
        "group.id":          "test",
        "auto.offset.reset": "earliest",
    })
    if err != nil {
        panic(err)
    }

    // 購読開始
    c.SubscribeTopics([]string{topic}, nil)

    // 受信
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

    fmt.Println("confluent-kafka-go-example stop.")
}