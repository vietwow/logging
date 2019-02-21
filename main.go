package main

import (
    "fmt"
    "os"
    // "github.com/vietwow/logging/consumer"
    "github.com/vietwow/logging/producer"
    // "strings"
)

func main() {
    // brokers := os.Getenv("BROKERS") // localhost:29092
    topic := os.Getenv("TOPIC") // heroku_logs
    group := os.Getenv("GROUP") // myGroup

    // Initialize kafka producer
    err := producer.InitKafka()
    if err != nil {
        fmt.Printf("Failed to create producer: %s\n", err)
        os.Exit(1)
    }

    message := "Hello Go!"

    producer.Produce(topic, message)
}
