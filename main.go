package main

import (
    "fmt"
    "github.com/confluentinc/confluent-kafka-go/kafka"
    "github.com/vietwow/logging/producer"
    "github.com/vietwow/logging/consumer"
)

func main() {
    // brokers := os.Getenv("BROKERS") // localhost:29092
    // topic := os.Getenv("TOPIC") // heroku_logs
    // group := os.Getenv("GROUP") // myGroup

    // Initialize kafka producer
    err = producer.InitKafka()
    if err != nil {
        log.Fatal("Kafka producer ERROR: ", err)
    }

    producerErr := producer.Produce(topics, string(messageJson))
    if producerErr != nil {
        log.Print(err)
    } else {
        messageResponse := fmt.Sprintf("Produced [%s] successfully", string(messageJson))
        fmt.Println(messageResponse)
    }

}
