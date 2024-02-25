package main

import (
	"context"
	"net/url"
)

func main() {

	// Create a cancellable context.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	url, err := url.Parse("https://example.com")
	if err != nil {
		panic(err)
	}

	// Create a new SQS client.
	client := aws

	// Create a new producer.
	producer := email.NewProducer(client)

}
