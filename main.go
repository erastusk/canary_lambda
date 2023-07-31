package main

import (
	"fmt"
	"github/erastusk/canary_lambda/api"
	"github/erastusk/canary_lambda/env"
	"log"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/retry"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/route53"
)

func main() {
	lambda.Start(Handler)
}
func Handler(request events.CloudWatchEvent) {
	//Create Logger
	logger := log.New(os.Stdout, "CANARY LAMBDA ---- ", log.Ldate|log.Ltime|log.Lshortfile)
	logger.Printf("Loading Environment...\n")

	// Populate Environment struct
	s, err := env.NewEnvLoad(logger)
	if err != nil {
		logger.Fatalf("%+v", err)
	}

	// On/Off Option
	if s.IsEnabled == "false" {
		logger.Println("Lambda function is disabled: s.IsEnabled Env viariable is false:", s.IsEnabled)
		return
	}

	//Load Default sdk behaviour configurations
	// Includes retry with backoff delay of 5s
	cfg, err := config.LoadDefaultConfig(s.Ctx, config.WithRegion("us-east-1"), config.WithRetryer(func() aws.Retryer {
		return retry.AddWithMaxBackoffDelay(retry.NewStandard(), time.Second*5)
	}))
	if err != nil {
		logger.Fatalf("%v", err)
	}

	// New route53 session from above config
	svc := route53.NewFromConfig(cfg)

	// Determine endpoint type to monitor, ALB or Dynamodb and pass to TryFunc
	switch s.EndpointType {
	case "alb":
		api.TryFunc(api.GetAlbStatus, s, svc)
	case "dynamodb":
		api.TryFunc(api.GetDynamoDbStatus, s, svc)
	default:
		logger.Fatalf("Endpoint missing or undefined: %v", s.EndpointType)
	}
	logger.Printf("Endpoint status check end...\n")
	logger.Println("Beginning weight validation")
	fmt.Println("-----------------------------------------------------------")

	// Validate current reights, normal 50/50 or failed over 100/0
	err = api.ValidateWeights(s, svc)
	if err != nil {
		logger.Fatalf("\n\t\t**********Could not validate Weights**********\n%v\n\n", err.Error())
	}
	// Take action, Failover or rebalance based on Failure and Balanced vars
	err = api.HandleWeightState(s, svc)
	if err != nil {
		logger.Fatalf("\n\t\t**********Weight Validation failed with the following error**********\n%v\n\n", err.Error())
	}
}
