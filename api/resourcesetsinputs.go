package api

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	r53_types "github.com/aws/aws-sdk-go-v2/service/route53/types"

	"github/erastusk/canary_lambda/env"
	"github/erastusk/canary_lambda/types"
)

//Leaving Function uncommented, to be used if/when CNAME changes are required.

func GetChangeResourceSetsInputsCname(s *env.EnvVariablesLoad, p *types.RecordUpdateVars, weight int64) *route53.ChangeResourceRecordSetsInput {
	return &route53.ChangeResourceRecordSetsInput{
		ChangeBatch: &r53_types.ChangeBatch{ // Required
			Changes: []r53_types.Change{ // Required
				{ // Required
					Action: r53_types.ChangeActionUpsert, // Required
					ResourceRecordSet: &r53_types.ResourceRecordSet{ // Required
						ResourceRecords: []r53_types.ResourceRecord{
							{Value: &p.DnsName},
						},
						Name:          aws.String(p.RecName), // Required
						Type:          r53_types.RRTypeCname, // Required
						TTL:           aws.Int64(1),
						Weight:        aws.Int64(weight),
						SetIdentifier: aws.String(p.SetIdentifier),
					},
				},
			},
			Comment: aws.String("Sample update."),
		},
		HostedZoneId: aws.String(s.ZoneId), // Required
	}
}

// Required Route53 Inputs configs set up to make weight changes.
func GetChangeResourceSetsInputsAlias(s *env.EnvVariablesLoad, p *types.RecordUpdateVars, weight int64) *route53.ChangeResourceRecordSetsInput {
	return &route53.ChangeResourceRecordSetsInput{
		ChangeBatch: &r53_types.ChangeBatch{ // Required
			Changes: []r53_types.Change{ // Required
				{ // Required
					Action: r53_types.ChangeActionUpsert, // Required
					ResourceRecordSet: &r53_types.ResourceRecordSet{ // Required
						AliasTarget: &r53_types.AliasTarget{
							DNSName:              aws.String(p.DnsName),
							EvaluateTargetHealth: true,
							HostedZoneId:         aws.String(s.ZoneId),
						},
						Name:          aws.String(p.RecName), // Required
						Type:          r53_types.RRTypeA,     // Required
						Weight:        aws.Int64(weight),
						SetIdentifier: aws.String(p.SetIdentifier),
					},
				},
			},
			Comment: aws.String("Sample update."),
		},
		HostedZoneId: aws.String(s.ZoneId), // Required
	}
}

// Required Route53 Inputs configs required for all GET record set listings from a given Zone ID.
func GetListRecordSetInput(s *env.EnvVariablesLoad) *route53.ListResourceRecordSetsInput {
	params_def_region := SetUpVars(s, s.DefaultRegion, s.SetID)
	return &route53.ListResourceRecordSetsInput{
		HostedZoneId:    aws.String(s.ZoneId),
		StartRecordName: aws.String(params_def_region.RecName),
		MaxItems:        aws.Int32(s.MaxItems),
	}
}
