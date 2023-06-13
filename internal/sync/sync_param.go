package sync

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	ssmType "github.com/aws/aws-sdk-go-v2/service/ssm/types"
)

func syncParam(ssmClient *ssm.Client, name string, arn string, destRegion string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*5))
	defer cancel()

	sourceRegion := GetRegion(arn)

	//ssmClient, cfg := auth.AwsAuth()
	getParamInput := ssm.GetParameterInput{
		Name:           &name,
		WithDecryption: aws.Bool(true),
	}

	decryptedSourceParam, err := ssmClient.GetParameter(ctx, &getParamInput)
	if err != nil {
		log.Printf("ERR: unable to retrieve paramenter: [%s] from [%s]", name, GetRegion(arn))
		log.Printf("ERR: %s", err)
		return
	}

	destParamName := strings.Replace(name, sourceRegion, destRegion, -1)
	destParameterInput := ssm.PutParameterInput{
		Name:     &destParamName,
		Value:    decryptedSourceParam.Parameter.Value,
		DataType: decryptedSourceParam.Parameter.DataType,
		Type:     decryptedSourceParam.Parameter.Type,
	}

	exists, sync := lookupDestinationParam(ctx, ssmClient, destParamName, destRegion, *decryptedSourceParam.Parameter.Value, sourceRegion)
	if exists && sync {
		// overwrite, no tagging
		destParameterInput.Overwrite = aws.Bool(true)
	} else if !exists && sync {
		// new parameter with tag
		tagKey := "ssm-replicated-from"
		destParameterInput.Tags = []ssmType.Tag{
			{
				Key:   &tagKey,
				Value: &sourceRegion,
			},
		}
	}

	if sync {
		_, err = ssmClient.PutParameter(ctx, &destParameterInput, func(options *ssm.Options) { options.Region = destRegion })
		if err != nil {
			log.Printf("ERR: unable to put paramenter: [%s] in [%s]", name, destRegion)
			log.Printf("ERR: %s", err)
			return
		}

		log.Printf("INFO: successful sync: [%s] region [%s] -> region [%s]", name, sourceRegion, destRegion)
	}
}
