package test

import (
	"github/erastusk/canary_lambda/api"
	"testing"
)

func TestSetUpVars(t *testing.T) {
	tables := []struct {
		region  string
		dnsname string
		recname string
		setid   string
	}{
		{"us-east-1", "web.us-east-1.sre.wsf.dev.npd.bfsaws.net", "web.sre.wsf.dev.npd.bfsaws.net", "web.us-east-1"},
		{"us-west-2", "web.us-west-2.sre.wsf.dev.npd.bfsaws.net", "web.sre.wsf.dev.npd.bfsaws.net", "web.us-west-2"},
	}
	for _, table := range tables {
		res := api.SetUpVars(s, table.region, "web.")
		if res.TargetRegion != table.region {
			t.Errorf("Region, Expected: %v, Got: %v", table.region, res.TargetRegion)
		}
		if res.DnsName != table.dnsname {
			t.Errorf("DNSName, Expected: %v, Got: %v", table.dnsname, res.DnsName)
		}
		if res.RecName != table.recname {
			t.Errorf("RECName, Expected: %v, Got: %v", table.recname, res.RecName)
		}
		if res.SetIdentifier != table.setid {
			t.Errorf("SetID, Expected: %v, Got: %v", table.setid, res.SetIdentifier)
		}
	}
}
