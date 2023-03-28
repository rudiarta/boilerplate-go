package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	handler := func(ctx context.Context, event events.APIGatewayProxyRequest) {

	}

	lambda.Start(handler)
}
