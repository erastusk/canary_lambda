package api

import (
	"github.com/aws/aws-sdk-go-v2/service/route53"

	"github/erastusk/canary_lambda/env"
	"github/erastusk/canary_lambda/types")

// Function manages weights.
// Incase of an endpoint failure, failover. Move traffic to 0 in failing region.
// Incase of unbalanced weights and no latency, rebalance weights. Restore traffic weights to 50/50 in both regions.
// Otherwise do nothing.
func HandleWeightState(s *env.EnvVariablesLoad, client *route53.Client) error {
	ch := make(chan string, 1)
	params_def_region := SetUpVars(s, s.DefaultRegion, s.SetID)
	params_alt_region := SetUpVars(s, s.AltRegion, s.SetID)
	wf_params_def_region := SetUpVars(s, s.DefaultRegion, s.WFSetID)
	wf_params_alt_region := SetUpVars(s, s.AltRegion, s.WFSetID)
	subj := "Workstation Framework Active Active Canary Lambda Notification: " + s.DefaultRegion
	// If endpoint failure is true AND balanced (50/50) AND Alternate failover region is live -> Failover
	if s.Failure && s.Balanced && s.AltRegState {
		msg := s.DefaultRegion + " ALB: " + s.Url + " Breached Failure Thresholds\n\nFailover Triggered: FROM " + s.DefaultRegion + " TO " + s.AltRegion
		s.Log.Println("Failover State detected, modifying weights")
		s.Log.Printf("Setting default region %s to %d", s.DefaultRegion, 0)
		s.Log.Printf("Setting Alternate region %s to %d", s.AltRegion, 100)
		r_east := GetChangeResourceSetsInputsCname(s, params_def_region, 0)
		r_west := GetChangeResourceSetsInputsCname(s, params_alt_region, 100)
		r_east_wf := GetChangeResourceSetsInputsCname(s, wf_params_def_region, 0)
		r_west_wf := GetChangeResourceSetsInputsCname(s, wf_params_alt_region, 100)
		_, err := client.ChangeResourceRecordSets(s.Ctx, r_east)
		if err != nil {
			return err
		}
		_, err = client.ChangeResourceRecordSets(s.Ctx, r_west)
		if err != nil {
			return err
		}
		_, err = client.ChangeResourceRecordSets(s.Ctx, r_east_wf)
		if err != nil {
			return err
		}
		_, err = client.ChangeResourceRecordSets(s.Ctx, r_west_wf)
		if err != nil {
			return err
		}
		go SNSSendEmail(ch, s, subj, msg)
		s.Log.Println(<-ch)
		close(ch)

		s.Log.Println("Weights modification completed")

		s.Balanced = false
		return nil
	}
	// If endpoint failure is true AND NOT balanced (0/100) -> Do nothing
	if s.Failure && !s.Balanced {
		s.Log.Println("App already in a failed over state with failures still occuring. Nothing to do...")
		return nil
	}
	// if weights not balanced -> rebalance
	if !s.Balanced {
		msg := s.DefaultRegion + " ALB: " + s.Url + " is no longer experiencing failures\n\nRolling Back Failover state: FROM " + s.AltRegion + " TO " + s.DefaultRegion + "\n\nReturning Weights back to 50/50"
		s.Log.Println("Unbalanced State detected, modifying weights")
		s.Log.Printf("Setting default region %s to %d", s.DefaultRegion, 50)
		s.Log.Printf("Setting Alternate region %s to %d", s.AltRegion, 50)
		r_east := GetChangeResourceSetsInputsCname(s, params_def_region, 50)
		r_west := GetChangeResourceSetsInputsCname(s, params_alt_region, 50)
		r_east_wf := GetChangeResourceSetsInputsCname(s, wf_params_def_region, 50)
		r_west_wf := GetChangeResourceSetsInputsCname(s, wf_params_alt_region, 50)
		_, err := client.ChangeResourceRecordSets(s.Ctx, r_east)
		if err != nil {
			return err
		}
		_, err = client.ChangeResourceRecordSets(s.Ctx, r_west)
		if err != nil {
			return err
		}
		_, err = client.ChangeResourceRecordSets(s.Ctx, r_east_wf)
		if err != nil {
			return err
		}
		_, err = client.ChangeResourceRecordSets(s.Ctx, r_west_wf)
		if err != nil {
			return err
		}
		go SNSSendEmail(ch, s, subj, msg)
		s.Log.Println(<-ch)
		close(ch)
		s.Log.Println("Weights modification completed")
		s.Balanced = true
		s.Failure = false
		return nil
	}
	return nil
}

// Instantiate Struct to be used for Route53 changes.

func SetUpVars(s *env.EnvVariablesLoad, region string, set_id string) *types.RecordUpdateVars {
	return &types.RecordUpdateVars{
		TargetRegion:  region,
		DnsName:       set_id + region + "." + s.EnvAlias + "." + s.HostedZoneName,
		RecName:       set_id + s.EnvAlias + "." + s.HostedZoneName,
		SetIdentifier: set_id + region,
	}
}
