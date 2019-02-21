package main

import (
    "fmt"
    "os"
    // "github.com/vietwow/logging/consumer"
    "github.com/vietwow/logging/producer"
)

func main() {
    topic := os.Getenv("TOPIC") // heroku_logs
    broker := os.Getenv("BROKERS") // localhost:29092
    // Initialize kafka producer
    err := producer.InitKafka(broker)
    if err != nil {
        fmt.Printf("Failed to create producer: %s\n", err)
        os.Exit(1)
    }

    message := "Hello Go!"

    producer.Produce(topic, message)
}
