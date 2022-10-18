package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

func usage() {
	fmt.Println("USAGE: sqs-getmessages -q <queue-name>")
}

func main() {
	queue := flag.String("q", "", "SQS Queue Name")
	flag.Parse()

	if *queue == "" {
		fmt.Println("SQS Queue name not specified!")
		usage()
		os.Exit(1)
	}

	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Fatalf("Error: %s", err)
		os.Exit(1)
	}

	client := sqs.NewFromConfig(cfg)

	urlResults, err := client.GetQueueUrl(context.Background(), &sqs.GetQueueUrlInput{
		QueueName: queue,
	})
	if err != nil {
		log.Fatalf("Error: %s", err)
		os.Exit(1)
	}
	queueUrl := urlResults.QueueUrl

	rmInput := &sqs.ReceiveMessageInput{
		MessageAttributeNames: []string{
			"Body",
		},
		QueueUrl:            queueUrl,
		MaxNumberOfMessages: 10,
		VisibilityTimeout:   30,
	}

	for {

		msgResults, err := client.ReceiveMessage(context.Background(), rmInput)
		if err != nil {
			log.Fatalf("Error: %s", err)
			os.Exit(1)
		}

		if msgResults.Messages == nil {
			log.Printf("Queue %s has no more visible messages.", *queue)
			break
		}
		for _, v := range msgResults.Messages {
			fmt.Println(*v.Body)
		}
	}
}
