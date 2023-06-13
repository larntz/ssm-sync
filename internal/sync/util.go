package sync

import (
	"context"
	"errors"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	ssmType "github.com/aws/aws-sdk-go-v2/service/ssm/types"
)

// GetRegion returns the region found in the parameter's ARN
func GetRegion(arn string) string {
	return strings.Split(arn, ":")[3]
}

// Check for existing destination parameter and `ssm-replicated-from` tag.
// exists = true if the parameter exists in the destination region
// sync = true if source != destination value.
func lookupDestinationParam(ctx context.Context, client *ssm.Client, name string, destRegion string, sourceValue string, sourceRegion string) (exists bool, sync bool) {
	getParamInput := ssm.GetParameterInput{
		Name:           &name,
		WithDecryption: aws.Bool(true),
	}

	param, err := client.GetParameter(ctx, &getParamInput, func(options *ssm.Options) { options.Region = destRegion })
	if err != nil {
		var pnf *ssmType.ParameterNotFound
		if errors.As(err, &pnf) {
			log.Printf("INFO: paramenter does not exist: [%s] in [%s]", name, destRegion)
			sync = true
		} else {
			log.Printf("ERR: unable to retrieve paramenter [%s] from [%s]", name, destRegion)
			log.Printf("ERR: %s", err)
		}
		return
	}

	exists = true

	output, err := client.ListTagsForResource(ctx, &ssm.ListTagsForResourceInput{ResourceId: param.Parameter.Name, ResourceType: "Parameter"}, func(options *ssm.Options) { options.Region = destRegion })
	if err != nil {
		log.Printf("failed to get tags for %s\nerr: %s", name, err)
	}

	tagRegion := ""
	for _, tag := range output.TagList {
		if *tag.Key == "ssm-replicated-from" {
			tagRegion = *tag.Value
		}
	}

	// this paramter was replicated from a different region or is not tagged `ssm-replicated-from`
	if tagRegion != sourceRegion {
		if tagRegion == "" {
			log.Printf("WARN: parameter exists, but is not tagged: [%s] in [%s]", name, destRegion)
		} else {
			log.Printf("WARN: parameter exists, but was replicated from a different region: [%s] in [%s], ['%s' != '%s']", name, destRegion, tagRegion, sourceRegion)
		}
		return
	}

	// source and destination values match, do not sync
	if *param.Parameter.Value == sourceValue {
		log.Printf("INFO: parameter exists and values match: [%s] in [%s]", name, destRegion)
		return
	}

	log.Printf("INFO: parameter exists, but needs updated: [%s] in [%s]", name, destRegion)
	sync = true
	return
}
