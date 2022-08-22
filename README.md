# ssm-sync

## Env Variables

| Name | Description | 
|-------------|-------------------------------------------------------------------|
| AWS_REGION  | required, the region we are replicating FROM                      |
| SSM_PATH    | required, this specifies the source ssm path, ex: "/my-ssm-path/" |
| LOCAL       | if set we are running on a server this is optional                |
| AWS_PROFILE | required if LOCAL is set                                          |
| HOSTNAME    | required of LOCAL is unset. Used in IRSA session name.            |

## Parameter tags

Parameters must have a `ssm-replicate-regions` tag or they are ignored. The tag value for
`ssm-replicate-regions` is a `:` separated list of region to replicate the parameter _to_.

Example: 

```
ssm-replicate-regions: us-east-1:us-west-1
```

## Required Permissions

Below is a sample AWS policy that shows the required permissions for syncing parameters:

```json
{
   "Statement": [
       {
           "Action": [
               "ssm:GetParametersByPath",
               "ssm:GetParameter",
               "ssm:PutParameter",
               "ssm:GetParameters",
               "ssm:ListTagsForResource",
               "ssm:AddTagsToResource"
           ],
           "Effect": "Allow",
           "Resource": "arn:aws:ssm:*:<account number>:parameter/*",
           "Sid": "ExternalSecrets"
       }
   ],
   "Version": "2012-10-17"
}
```
