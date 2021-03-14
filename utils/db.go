package main

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	url := "mongodb+srv://turingcoffee:test12345@cluster0.m5idb.gcp.mongodb.net/testturingcoffee?retryWrites=true&w=majority"
	ctx := context.TODO()
	options := options.Client()
	options.SetMaxPoolSize(2)
	options.ApplyURI(url)

	fmt.Printf("Connecting to databa-se %s\n", url)
	client, err := mongo.Connect(ctx, options)
	if err != nil {
		fmt.Errorf("mongo.Connect() failed: %v", err.Error())
	}

	fmt.Printf("Ping...\n")
	err = client.Ping(ctx, nil)
	if err != nil {
		fmt.Errorf("Ping mongodb failed: %v", err.Error())
	}

	fmt.Printf("Disconnect...\n")
	client.Disconnect(ctx)
}
