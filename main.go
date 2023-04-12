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
	deadline, _ := ctx.Deadline()
	deadline = deadline.Add(-100 * time.Millisecond)
	timeoutChannel := time.After(time.Until(deadline))

	for {
		select {
		case <-timeoutChannel:
			logger.Printf("Deadline exceeded: %v", deadline)
			return events.APIGatewayProxyResponse{Body: "Internal Server Error", StatusCode: 500}, nil
		default:
			report, err := lido.FetchRewardsReport("0x03d04a5F3cc050aB69A84eB0Da3242cd84DBf724")
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
