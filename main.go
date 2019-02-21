package main

import (
    "fmt"
    "github.com/confluentinc/confluent-kafka-go/kafka"
    "github.com/vietwow/logging/consumer"
    "github.com/vietwow/logging/producer"
)

func main() {
    brokers := os.Getenv("BROKERS") // localhost:29092
    topic := os.Getenv("TOPIC") // heroku_logs
    group := os.Getenv("GROUP") // myGroup

    producerErr := producer.Produce(strings.Fields(topic), string(messageJson))
    if producerErr != nil {
        log.Print(err)
    } else {
        messageResponse := fmt.Sprintf("Produced [%s] successfully", string(messageJson))
        fmt.Println(messageResponse)
    }

}
