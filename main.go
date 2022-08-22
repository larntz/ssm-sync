package main

import (
	"log"
	"os"
	"time"
)

// ENV Variables
// -----
// AWS_REGION: required, the region we are replicating FROM
// SSM_PATH: required, this specifies the source ssm path, ex: "/my-ssm-path/"
// LOCAL: if set we are running on a server this is optional
// AWS_PROFILE: required if LOCAL is set
// HOSTNAME: required of LOCAL is unset. used in session name
// -----

// Parameter tags
// -----
// Parameters need to have a special tag to be replicated.
// tag key:ssm-replicate-regions
// tag values are 1 or more AWS regions separated by a :
// example: us-west-1:us-west-2
// -----

// Required Permissions
// {
//     "Statement": [
//         {
//             "Action": [
//                 "ssm:GetParametersByPath",
//                 "ssm:GetParameter",
//                 "ssm:PutParameter",
//                 "ssm:GetParameters",
//                 "ssm:ListTagsForResource",
//                 "ssm:AddTagsToResource"
//             ],
//             "Effect": "Allow",
//             "Resource": "arn:aws:ssm:*:<account number>:parameter/*",
//             "Sid": "SSMSync"
//         }
//     ],
//     "Version": "2012-10-17"
// }

func main() {
	log.Printf("ssm-sync starting")

	ssmClient, _ := awsAuth()
	for {
		log.Printf("starting sync... ")
		syncAllParams(ssmClient, os.Getenv("SSM_PATH"))
		time.Sleep(1 * time.Minute)
	}
}
