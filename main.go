package main

import (
    "fmt"
    "os"
    "github.com/vietwow/logging/consumer"
)

func main() {
    topic := os.Getenv("TOPIC") // heroku_logs
    broker := os.Getenv("BROKER") // localhost:29092

    // Initialize kafka consumer
    err = consumer.InitKafka(broker)
    if err != nil {
        fmt.Printf("Failed to create consumer: %s\n", err)
        os.Exit(1)
    }

    consumer.Consume(topic)
}
