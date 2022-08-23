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
		log.Printf("ERR: unable to retrieve paramenter [%s] from [%s]", name, GetRegion(arn))
		log.Printf("ERR: %s", err)
		return
	}

	destParamName := strings.Replace(name, sourceRegion, region, -1)
	destCfg := cfg
	destCfg.Region = region
	destSsmClient := ssm.NewFromConfig(destCfg)
	destParameterInput := ssm.PutParameterInput{
		Name:     &destParamName,
		Value:    decryptedSourceParam.Parameter.Value,
		DataType: decryptedSourceParam.Parameter.DataType,
		Type:     decryptedSourceParam.Parameter.Type,
	}

	exists, tagRegion := lookupDestinationParam(ctx, destSsmClient, destParamName, region)
	if exists {
		if tagRegion != sourceRegion {
			// this paramter was replicated from a different region. Skip it.
			log.Printf("WARN: parameter [%s] in [%s] exists, but was replicated from a different region or is not tagged. ['%s' != '%s']", name, region, tagRegion, sourceRegion)
			return
		}
		// put with overwrite no tagging
		destParameterInput.Overwrite = true
	} else {
		// put new parameter with tag
		tagKey := "ssm-replicated-from"
		destParameterInput.Tags = []ssmType.Tag{
			{
				Key:   &tagKey,
				Value: &sourceRegion,
			},
		}
	}

	_, err = destSsmClient.PutParameter(ctx, &destParameterInput)
	if err != nil {
		log.Printf("ERR: unable to put paramenter [%s] in [%s]", name, region)
		log.Printf("ERR: %s", err)
		return
	}

	log.Printf("INFO: successfully syncd [%s] from region [%s] to region [%s]", name, sourceRegion, region)
}
