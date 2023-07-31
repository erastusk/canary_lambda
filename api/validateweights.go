package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/route53"

	"github/erastusk/canary_lambda/env"
	"github/erastusk/canary_lambda/types"
)

// Get formatted record set inputs to be used to query route53. Validates if record exists before making changes
// If not return error and fail.
func ValidateWeights(s *env.EnvVariablesLoad, client *route53.Client) error {
	params := GetListRecordSetInput(s)
	// Default timeout limit is 60s if not set, setting it to 5s
	// ctx, cancel := context.WithTimeout(s.Ctx, time.Duration(time.Millisecond))
	// defer cancel()
	resp, err := client.ListResourceRecordSets(context.Background(), params)
	if err != nil {
		return err
	}
	if len(resp.ResourceRecordSets) == 0 {
		s.Log.Fatal("No RecordSets found")
	}
	out, err := UnmarshalRecordSetRes(s, resp)
	if err != nil {
		return err
	}

	//Set s.balanced to return value of GetWeighState. If true weights are balanced otherwise in a failover state.
	s.Balanced = GetWeightState(s, out)
	return nil
}

// Unmarshall route53 json response into a struct
func UnmarshalRecordSetRes(s *env.EnvVariablesLoad, resp *route53.ListResourceRecordSetsOutput) ([]types.ListRecordsResponse, error) {
	var out []types.ListRecordsResponse
	e, err := json.Marshal(&resp.ResourceRecordSets)
	if err != nil {
		s.Log.Fatal("Couldn't Marshal Resource Record Sets")
		return nil, err
	}
	err = json.Unmarshal(e, &out)
	if err != nil {
		s.Log.Fatal("Couldn't UnMarshal Resource Record Sets")
		return nil, err
	}
	return out, nil
}

// Check if response contains records with provided zoneid and set id
// Check weight state, if in failover state (100/0) return false otherwise return true
func GetWeightState(s *env.EnvVariablesLoad, out []types.ListRecordsResponse) bool {
	var defreg, altreg int
	for _, record := range out {
		if record.SetIdentifier == s.SetID+s.DefaultRegion {
			s.Log.Printf("Current %s weight = %d", s.DefaultRegion, record.Weight)
			defreg = record.Weight
		}
		if record.SetIdentifier == s.SetID+s.AltRegion {
			s.Log.Printf("Current %s weight = %d", s.AltRegion, record.Weight)
			altreg = record.Weight
		}
	}
	// If lambda in east has encountered failures and failed over to west, and the next time the lambda is invoked
	// and latency no longer exists, rebalance weights to 50/50, return false
	// if defreg == 0 && altreg == 100 {
	// 	return false
	// }
	// If weights are balanced and there are no failures, Do Nothing
	if defreg == 50 && altreg == 50 && !s.Failure {
		s.Log.Println("Weight validation completed, weights in normal state with no failures, nothing to do....")
		return true
	}
	alt_reg_status := GetAltRegionResp(s)
	// If traffic has been failed over to a region, example west due to failures in east. And lambda in west
	// also encounters latency, meaning both regions would be experiencing latency issues. Do nothing
	// Prevent ping-pong scenario where both lambda's trigger failovers in a loop
	if !alt_reg_status {
		s.Log.Printf("%s is in a failure state, can't failover regardless of default region state, Do Nothing", s.AltRegion)
		return true
	}
	// If Both regions are accessible, rebalance weights.
	if defreg == 100 || altreg == 100 && alt_reg_status && !s.Failure {
		s.Log.Printf("%s is in a non-failure state, Rebalance", s.AltRegion)
		return false
	}
	if defreg == 100 || altreg == 100 && alt_reg_status && s.Failure {
		return false
	}
	return true
}

func GetAltRegionResp(s *env.EnvVariablesLoad) bool {
	s.Log.Printf("Checking if alternate region %s is healthy", s.AltRegion)
	// Determine endpoint type to monitor, ALB or Dynamodb and pass to TryFunc
	switch s.EndpointType {
	case "alb":
		resp, err := net.LookupHost(s.AltUrl)
		if err != nil {
			fmt.Printf("+++++++++++++++++++++++++++++++++++++++++++++++++++++++\n")
			s.Log.Printf("\nLookup failed : %v\n", err)
			fmt.Printf("+++++++++++++++++++++++++++++++++++++++++++++++++++++++\n")
			s.AltRegState = false
			return false
		}
		if len(resp) > 0 {
			return true
		}
	case "dynamodb":
		if !strings.Contains(s.AltRegion, "http") {
			return false
		}
		_, err := http.Get(s.AltUrl)
		if err != nil {
			fmt.Printf("+++++++++++++++++++++++++++++++++++++++++++++++++++++++\n")
			s.Log.Printf("\nGet Url failed : %v\n", err)
			fmt.Printf("+++++++++++++++++++++++++++++++++++++++++++++++++++++++\n")
			s.AltRegState = false
			return false
		}
		return true
	}
	return true
}
