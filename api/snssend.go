package api

import (
	"context"
	"github/erastusk/canary_lambda/env"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

// PublishMessage publishes a message to an Amazon Simple Notification Service (Amazon SNS) topic
// Inputs:
//
//	c is the context of the method call, which includes the Region
//	api is the interface that defines the method call
//	input defines the input arguments to the service call.
//
// Output:
//
//	If success, a PublishOutput object containing the result of the service call and nil
//	Otherwise, nil and an error from the call to Publish

func SNSSendEmail(ch chan string, s *env.EnvVariablesLoad, subj string, msg string) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(s.DefaultRegion))
	if err != nil {
		s.Log.Println("SNS cfg configuration failed, " + err.Error())
	}
	client := sns.NewFromConfig(cfg)

	// Get SNS topic ARN
	topic, err := SNSGetTopic(s, client)
	if err != nil {
		s.Log.Printf("%v", err)
	}
	if topic == "" {
		s.Log.Println("No SNS topic found")
	}
	input := &sns.PublishInput{
		Message:  &msg,
		Subject:  &subj,
		TopicArn: &topic,
	}
	// Publish to SNS Topic
	_, err = client.Publish(context.TODO(), input)
	if err != nil {
		s.Log.Println("Could not send/publish email..")
		s.Log.Println(err)
	}
	ch <- "Successfully published SNS message"
}

// Get SNS Topics and filter out alerts-notification arn. Default topic created with all AWS accounts
func SNSGetTopic(s *env.EnvVariablesLoad, c *sns.Client) (string, error) {
	input := &sns.ListTopicsInput{}

	result, err := c.ListTopics(context.TODO(), input)
	if err != nil {
		s.Log.Println("Error while listing Topics")
		return "", err
	}
	if len(result.Topics) == 0 {
		return "", nil
	}
	for _, t := range result.Topics {
		if strings.Contains(*t.TopicArn, "alerts-notification") {
			return *t.TopicArn, nil
		}
	}
	return "", nil
}
