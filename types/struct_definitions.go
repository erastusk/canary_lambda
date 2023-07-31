package types

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/sns"
)

type ListRecordsResponse struct {
	AliasTarget struct {
		DNSName              string `json:"DNSName"`
		EvaluateTargetHealth bool   `json:"EvaluateTargetHealth"`
		HostedZoneID         string `json:"HostedZoneId"`
	} `json:"AliasTarget"`
	Name          string `json:"Name"`
	SetIdentifier string `json:"SetIdentifier"`
	Type          string `json:"Type"`
	Weight        int    `json:"Weight"`
}

type RecordUpdateVars struct {
	DnsName       string
	RecName       string
	SetIdentifier string
	TargetRegion  string
}

// SNSListTopicsAPI defines the interface for the ListTopics function.
// We use this interface to test the function using a mocked service.
type SNSListTopicsAPI interface {
	ListTopics(ctx context.Context,
		params *sns.ListTopicsInput,
		optFns ...func(*sns.Options)) (*sns.ListTopicsOutput, error)
}

// SNSPublishAPI defines the interface for the Publish function.
// We use this interface to test the function using a mocked service.
type SNSPublishAPI interface {
	Publish(ctx context.Context,
		params *sns.PublishInput,
		optFns ...func(*sns.Options)) (*sns.PublishOutput, error)
}
