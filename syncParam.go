package main

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/ssm"
	ssmType "github.com/aws/aws-sdk-go-v2/service/ssm/types"
)

func syncParam(name string, arn string, region string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*5))
	defer cancel()

	sourceRegion := GetRegion(arn)

	ssmClient, cfg := awsAuth()
	getParamInput := ssm.GetParameterInput{
		Name:           &name,
		WithDecryption: true,
	}

	decryptedSourceParam, err := ssmClient.GetParameter(ctx, &getParamInput)
	if err != nil {
		log.Printf("ERROR: unable to retrieve paramenter [%s] from [%s]", name, GetRegion(arn))
		log.Printf("ERROR: %s", err)
		return
	}

	destinationParamName := strings.Replace(name, GetRegion(arn), region, -1)
	destinationParameterInput := ssm.PutParameterInput{
		Name:      &destinationParamName,
		Value:     decryptedSourceParam.Parameter.Value,
		DataType:  decryptedSourceParam.Parameter.DataType,
		Type:      decryptedSourceParam.Parameter.Type,
		Overwrite: true,
	}

	destCfg := cfg
	destCfg.Region = region
	destSsmClient := ssm.NewFromConfig(destCfg)
	_, err = destSsmClient.PutParameter(ctx, &destinationParameterInput)
	if err != nil {
		log.Printf("ERROR: unable to retrieve paramenter [%s] from [%s]", name, sourceRegion)
		log.Printf("ERROR: %s", err)
		return
	}

	tagKey := "ssm-replicated-from"
	destinationParamTags := []ssmType.Tag{
		{
			Key:   &tagKey,
			Value: &sourceRegion,
		},
	}

	addTagsInput := ssm.AddTagsToResourceInput{
		ResourceId:   &destinationParamName,
		Tags:         destinationParamTags,
		ResourceType: "Parameter",
	}
	_, err = destSsmClient.AddTagsToResource(ctx, &addTagsInput)
	if err != nil {
		log.Printf("ERROR: unable to tag paramenter [%s] from [%s]", name, sourceRegion)
		log.Printf("ERROR: %s", err)
		return
	}

	log.Printf("successfully syncd [%s] from region [%s] to region [%s]", name, sourceRegion, region)
}
