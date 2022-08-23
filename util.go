package main

import (
	"context"
	"errors"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/ssm"
	ssmType "github.com/aws/aws-sdk-go-v2/service/ssm/types"
)

// GetRegion returns the region found in the parameter's ARN
func GetRegion(arn string) string {
	return strings.Split(arn, ":")[3]
}

// Check for existing destination parameter and `ssm-replicated-from` tag. Return true if found and the value of
// `ssm-replicated-from` tag.
func lookupDestinationParam(ctx context.Context, client *ssm.Client, name string, region string) (exists bool, sourceRegion string) {
	getParamInput := ssm.GetParameterInput{
		Name:           &name,
		WithDecryption: false,
	}

	param, err := client.GetParameter(ctx, &getParamInput)
	if err != nil {
		var pnf *ssmType.ParameterNotFound
		if errors.As(err, &pnf) {
			log.Printf("INFO: paramenter does not exist: [%s] in [%s]", name, region)
		} else {
			log.Printf("ERR: unable to retrieve paramenter [%s] from [%s]", name, region)
			log.Printf("ERR: %s", err)
		}
		return
	}

	exists = true

	output, err := client.ListTagsForResource(ctx, &ssm.ListTagsForResourceInput{ResourceId: param.Parameter.Name, ResourceType: "Parameter"})
	if err != nil {
		log.Printf("failed to get tags for %s\nerr: %s", name, err)
	}

	for _, tag := range output.TagList {
		if *tag.Key == "ssm-replicated-from" {
			sourceRegion = *tag.Value
		}
	}
	return
}
