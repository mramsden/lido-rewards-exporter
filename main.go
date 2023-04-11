package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var logger = log.New(os.Stdout, "[Handler] ", log.LstdFlags)

func handleRequest(ctx context.Context, _ events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	deadline, _ := ctx.Deadline()
	deadline = deadline.Add(-100 * time.Millisecond)
	timeoutChannel := time.After(time.Until(deadline))

	for {
		select {
		case <-timeoutChannel:
			logger.Printf("Deadline exceeded: %v", deadline)
			return events.APIGatewayProxyResponse{Body: "Internal Server Error", StatusCode: 500}, nil
		default:
			return events.APIGatewayProxyResponse{Body: "Hello, world!", StatusCode: 200}, nil
		}
	}
}

func main() {
	lambda.Start(handleRequest)
}
