package test

import (
	"encoding/json"
	"github/erastusk/canary_lambda/api"
	"github/erastusk/canary_lambda/types"
	"io/ioutil"
	"testing"
)

var out []types.ListRecordsResponse

func TestGetWeightState(t *testing.T) {
	r, err := ioutil.ReadFile("../data/route53aliasstruct.json")
	if err != nil {
		t.Logf("Unable to Read File: %+v", err)
	}
	err = json.Unmarshal(r, &out)
	if err != nil {
		t.Fatal("Couldn't UnMarshal Resource Record Sets")
	}
	z := api.GetWeightState(s, out)
	if z != true {
		t.Errorf("Got: %v, expected %v", true, z)
	}
}
