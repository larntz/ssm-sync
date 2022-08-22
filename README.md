# ssm-sync

## Env Variables

### Running in cluster

| Name | Description | 
|-------------|--------------------------------------------------------------------|
| AWS_REGION  | required, the region we are replicating from.                      |
| SSM_PATH    | required, this specifies the source ssm path, ex: "/my-ssm-path/". |
| HOSTNAME    | required, used in IRSA session name.                               |

### Local testing
| Name | Description | 
|-------------|-------------------------------------------------------------------|
| LOCAL       | if set run using a local aws profile for testing.                 |
| AWS_PROFILE | name of the aws profile to use when running locally.              |

## Parameter tags

### Source parameter

Parameters must have a `ssm-replicate-regions` tag or they are ignored. The tag value for
`ssm-replicate-regions` is a `:` separated list of region to replicate the parameter _to_.

Example: 

```
ssm-replicate-regions: us-east-1:us-west-1
```

### Destination parameter

Replicated parameters get tagged with `ssm-replicated-from`. The value of these tags is the source region the parameter was replicated from.

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
           "Sid": "SSMSync"
       }
   ],
   "Version": "2012-10-17"
}
```
