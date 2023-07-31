package test

import (
	"github/erastusk/canary_lambda/api"
	"github/erastusk/canary_lambda/env"
	"log"
	"os"
	"testing"
)

var logger = log.New(os.Stdout, "Canary Lambda Unit Testing---- ", log.Ldate|log.Ltime|log.Lshortfile)
var s = &env.EnvVariablesLoad{
	Url:            "web.us-west-2.sre.wsf.dev.npd.bfsaws.net",
	DefaultRegion:  "us-east-1",
	AltRegion:      "us-west-2",
	EnvAlias:       "sre",
	Env:            "dev",
	ZoneId:         "Z07944991LMIQ6M5XRXDI",
	Log:            logger,
	SetID:          "web.",
	HostedZoneName: "wsf.dev.npd.bfsaws.net",
}

func TestAlb(t *testing.T) {
	tables := []struct {
		url string
		res bool
	}{

		{"git.devops.broadridge.net", true},
		{"nexus.devops.broadridge.net", true},
		{"", false},
	}

	for _, table := range tables {
		s.Url = table.url
		ret := api.GetAlbStatus(s)
		if ret != table.res {
			t.Errorf("Expected: %v, got: %v", table.res, ret)

		}

	}

}
func TestDynamodb(t *testing.T) {
	tables := []struct {
		url string
		res bool
	}{
		{"https://git.devops.broadridge.net", true},
		{"https://nexus.devops.broadridge.net", true},
		{"", false},
	}

	for _, table := range tables {
		s.Url = table.url
		ret := api.GetDynamoDbStatus(s)
		if ret != table.res {
			t.Errorf("Expected: %v, got: %v", table.res, ret)

		}

	}

}
