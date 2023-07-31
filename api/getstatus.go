package api

import (
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github/erastusk/canary_lambda/env"

	"github.com/aws/aws-sdk-go-v2/service/route53"
)

// TryFunc manages endpoint retries 5 times within a minute with a 10s sleep time.
func TryFunc(m func(s *env.EnvVariablesLoad) bool, s *env.EnvVariablesLoad, svc *route53.Client) {
	failure, tries := 0, 0
	//Convert Env strings to INT, otherwise set defaults
	sleep_t, err := strconv.Atoi(s.Wait)
	if err != nil {
		sleep_t = 10
	}
	fail_count, err := strconv.Atoi(s.FailureCount)
	if err != nil {
		fail_count = 3
	}
	duration, err := strconv.Atoi(s.Duration)
	if err != nil {
		duration = 50
	}
	now := time.Now()
	then := now.Add(time.Second * time.Duration(duration))

	s.Log.Printf("\nBegninning Endpoint status check...\n")
	for then.After(time.Now()) {
		tries++
		s.Log.Println("Try..", tries)
		time.Sleep(time.Second * time.Duration(sleep_t))

		// m function is either GetAlbStatus or GetDynamoDbStatus function passed from main.
		result := m(s)
		if !result {
			failure++
		}
		// If failure exceeds or equals a value, set s.Failure to true to failover.
		if failure >= fail_count {
			s.Failure = true
			fmt.Println("----------------------------------------------------------")
			s.Log.Printf("\n Endpoint status exceeded failure Tolerance of %v tries ", s.FailureCount)
			fmt.Println("----------------------------------------------------------")
			break
		}
	}
}

// First check if url containts http/https, if not fail immmediately.
// If no err or 200 response codes return a true  otherwise return false to fail over.
func GetDynamoDbStatus(s *env.EnvVariablesLoad) bool {
	if !strings.Contains(s.Url, "http") {
		return false
	}
	resp, err := http.Get(s.Url)
	if err != nil {
		fmt.Printf("+++++++++++++++++++++++++++++++++++++++++++++++++++++++\n")
		s.Log.Printf("\nGet Url failed : %v\n", err)
		fmt.Printf("+++++++++++++++++++++++++++++++++++++++++++++++++++++++\n")
		return false
	}
	s.Log.Printf("Get URL successfull for %v, response code : %v", s.Url, resp.StatusCode)
	return true
}

// If no err or a valid Host IP return a true  otherwise return false to fail over.
func GetAlbStatus(s *env.EnvVariablesLoad) bool {
	resp, err := net.Dial("tcp", s.Url+":https")
	if err != nil || resp == nil {
		fmt.Printf("+++++++++++++++++++++++++++++++++++++++++++++++++++++++\n")
		s.Log.Printf("\nError: %v\n", err)
		fmt.Printf("+++++++++++++++++++++++++++++++++++++++++++++++++++++++\n")
		return false
	}
	s.Log.Printf("TCP connection succeeded for : %v\n", s.Url)
	return true
}
