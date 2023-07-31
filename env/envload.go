package env

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type EnvVariablesLoad struct {
	Log            *log.Logger
	Url            string
	AltUrl         string
	WFUrl          string
	WFAltUrl       string
	ZoneId         string
	DefaultRegion  string
	AltRegion      string
	EnvAlias       string
	Env            string
	MaxItems       int32
	Balanced       bool
	Failure        bool
	EndpointType   string
	Ctx            context.Context
	IsEnabled      string
	Wait           string
	FailureCount   string
	Duration       string
	SetID          string
	WFSetID        string
	AltRegState    bool
	HostedZoneName string
}

// Environment variables from Env.
func NewEnvLoad(logger *log.Logger) (*EnvVariablesLoad, error) {
	suc := ValidateEnvVars()
	if os.Getenv("enabled") == "" {
		os.Setenv("enabled", "true")
	}
	if os.Getenv("wait") == "" {
		os.Setenv("wait", "10")
	}
	if os.Getenv("failurecount") == "" {
		os.Setenv("failurecount", "3")
	}
	if os.Getenv("duration") == "" {
		os.Setenv("duration", "50")
	}
	if !suc || os.Getenv("zonename") == "" || os.Getenv("setid") == "" || os.Getenv("default_region") == "" || os.Getenv("alt_region") == "" || os.Getenv("envalias") == "" || os.Getenv("env") == "" || os.Getenv("zoneid") == "" {
		return nil, errors.New("unable to load environment, missing environment variables")
	}
	ctx := context.Background()
	return &EnvVariablesLoad{
		Log:      logger,
		Url:      os.Getenv("setid") + os.Getenv("default_region") + "." + os.Getenv("envalias") + "." + os.Getenv("zonename"),
		AltUrl:   os.Getenv("setid") + os.Getenv("alt_region") + "." + os.Getenv("envalias") + "." + os.Getenv("zonename"),
		WFUrl:    os.Getenv("setid") + os.Getenv("default_region") + "." + os.Getenv("envalias") + "." + os.Getenv("zonename"),
		WFAltUrl: os.Getenv("setid") + os.Getenv("alt_region") + "." + os.Getenv("envalias") + "." + os.Getenv("zonename"), DefaultRegion: os.Getenv("default_region"),
		AltRegion:      os.Getenv("alt_region"),
		EnvAlias:       os.Getenv("envalias"),
		Env:            os.Getenv("env"),
		ZoneId:         os.Getenv("zoneid"),
		HostedZoneName: os.Getenv("zonename"),
		MaxItems:       2,
		Balanced:       true,
		Failure:        false,
		EndpointType:   "alb",
		Ctx:            ctx,
		IsEnabled:      os.Getenv("enabled"),
		Wait:           os.Getenv("sleep"),
		FailureCount:   os.Getenv("failurecount"),
		Duration:       os.Getenv("duration"),
		SetID:          os.Getenv("setid"),
		WFSetID:        os.Getenv("wfsetid"),
		AltRegState:    true,
	}, nil
}

func LoadVars() error {
	err := godotenv.Load()
	if err != nil {
		return err
	}
	return nil

}

func ValidateEnvVars() bool {

	_, ok := os.LookupEnv("zonename")
	if !ok {
		fmt.Println("Missing zonename environment variable")
		return false
	}
	_, ok = os.LookupEnv("setid")
	if !ok {
		fmt.Println("Missing setid environment variable")
		return false
	}
	_, ok = os.LookupEnv("default_region")
	if !ok {
		fmt.Println("Missing default_region environment variable")
		return false
	}

	_, ok = os.LookupEnv("alt_region")
	if !ok {
		fmt.Println("Missing alt_region environment variable")
		return false
	}

	_, ok = os.LookupEnv("envalias")
	if !ok {
		fmt.Println("Missing envalias environment variable")
		return false
	}

	_, ok = os.LookupEnv("env")
	if !ok {
		fmt.Println("Missing env environment variable")
		return false
	}

	_, ok = os.LookupEnv("zoneid")
	if !ok {
		fmt.Println("Missing zoneid environment variable")
		return false
	}
	return true
}
