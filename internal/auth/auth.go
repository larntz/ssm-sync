package auth

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
)

func AwsAuth() (*ssm.Client, aws.Config) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*2))
	defer cancel()

	// running locally or in cluster?
	_, local := os.LookupEnv("LOCAL")
	if local {
		cfg, err := config.LoadDefaultConfig(ctx, config.WithSharedConfigProfile(os.Getenv("AWS_PROFILE")))
		if err != nil {
			log.Fatal("config error:", err)
		}
		return ssm.NewFromConfig(cfg), cfg
	}

	os.Getenv("AWS_REGION")
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(os.Getenv("AWS_REGION")),
		config.WithWebIdentityRoleCredentialOptions(func(options *stscreds.WebIdentityRoleOptions) {
			options.RoleSessionName = "IRSA_SSM_SYNC@" + os.Getenv("HOSTNAME")
		}))
	if err != nil {
		log.Fatal("config error:", err)
	}
	return ssm.NewFromConfig(cfg), cfg
}
