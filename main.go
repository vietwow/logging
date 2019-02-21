package main

import (
    "fmt"
    "os"
    // "github.com/vietwow/logging/consumer"
    "github.com/vietwow/logging/producer"
    "log"
    // "strings"
)

type producedMessage struct {
    Id string      `json:"id"`
    Message string `json:"message"`
}

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

    messageJson, _ := json.Marshal(&producedMessage{
        Id: message.ID,
        Message: message.Text,
    })

    producerErr := producer.Produce(topic, string(messageJson))
    if producerErr != nil {
        log.Print(err)
    } else {
        messageResponse := fmt.Sprintf("Produced [%s] successfully", string(messageJson))
        fmt.Println(messageResponse)
    }

}
