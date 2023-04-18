package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/mramsden/lido-rewards-exporter/lido"
)

var logger = log.New(os.Stdout, "[Handler] ", log.LstdFlags)

func handleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	address := request.PathParameters["address"]
	if address == "" {
		logger.Printf("Bad request missing address path parameter in path %s", request.Path)
		return events.APIGatewayProxyResponse{Body: "Bad Request", StatusCode: 400}, nil
	}

	deadline, _ := ctx.Deadline()
	deadline = deadline.Add(-100 * time.Millisecond)
	timeoutChannel := time.After(time.Until(deadline))

	for {
		select {
		case <-timeoutChannel:
			logger.Printf("Deadline exceeded: %v", deadline)
			return events.APIGatewayProxyResponse{Body: "Internal Server Error", StatusCode: 500}, nil
		default:
			report, err := lido.FetchRewardsReport(address)
			if err != nil {
				return events.APIGatewayProxyResponse{Body: "Internal Server Error", StatusCode: 500}, err
			}

			b, err := json.Marshal(report)
			if err != nil {
				return events.APIGatewayProxyResponse{Body: "Internal Server Error", StatusCode: 500}, err
			}
			return events.APIGatewayProxyResponse{Body: string(b), StatusCode: 200}, nil
		}
	}
}

func main() {
	lambda.Start(handleRequest)
}
