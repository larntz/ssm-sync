package sync

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
)

// Parameters grabs parameters from ssm
func Parameters(ssmClient *ssm.Client, path string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*45))
	defer cancel()

	ssmPathInput := ssm.GetParametersByPathInput{
		Path:           &path,
		Recursive:      aws.Bool(true),
		WithDecryption: aws.Bool(false),
	}

	for {
		ssmParameters, err := ssmClient.GetParametersByPath(ctx, &ssmPathInput)
		if err != nil {
			fmt.Println("Error getting ssmParameters")
			panic(err)
		}

		for _, param := range ssmParameters.Parameters {
			output, err := ssmClient.ListTagsForResource(ctx, &ssm.ListTagsForResourceInput{ResourceId: param.Name, ResourceType: "Parameter"})
			if err != nil {
				log.Printf("failed to get tags for %s\nerr: %s", *param.ARN, err)
			}

			for _, tag := range output.TagList {
				// sync to destination regions
				if *tag.Key == "ssm-replicate-regions" {
					for _, region := range strings.Split(*tag.Value, ":") {
						go syncParam(ssmClient, *param.Name, *param.ARN, region)
					}
				}
			}
		}

		if ssmParameters.NextToken == nil {
			break
		} else {
			ssmPathInput.NextToken = ssmParameters.NextToken
		}
	}
}
