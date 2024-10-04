Running tool

Step 1. 
    Build Tool
```    
    go build -o partnerkeygen .
```
Step 2.
    Run tool
```
    ./partnerkeygen -issuer="<value>" -sec_kid="<value>" -aws_key_id="<value>" -aws_secret_key="<value>" -aws_token="<value>" -region="<value>" -tenant_id="<value>" -user_id="<value>" -environment="<value>" -product_id="<value>" -service_id="<value>"
```
  
| Flag           | Meaning                                                    |
|----------------|------------------------------------------------------------|
| issuer         | domain value like "api-dev01-user7.project-amber-smas.com" |
| sec_kid        | secret id      (*A)                                        |
| aws_key_id     | aws key id     (*B)                                        |
| aws_secret_key | aws secret key (*B)                                        |
| aws_token      | aws token      (*B)                                        |                        
| region         | aws region                                                 |                       
| tenant_id      | parter tenant id                                           |                     
| user_id        | partener user id                                           |                     
| environment    | full environment name like "dev01-user7"                   | 
| product_id     | product id                                                 | 
| service_id     | Service id                                                 |

```
*A: sec_id: This is aws secret id where we store certificate to sign token. for test1 you can user "amber-apigw-jwt-signing-us-east-1-test1"
    This value can be found in environment variable (**KEY_ID**) of lambda function of given environment in aws console.
    e.g. https://us-east-1.console.aws.amazon.com/lambda/home?region=us-east-1#/functions/amber_attestation_authorizer-us-east-1-test1-user1?tab=configure
```
```
*B: Following step are the same step which we are using to get access key and secret key and aws token to connect to aks cluster

Step 1 : Fetch IAM user iam_devops access keys
    export AWS_ACCESS_KEY_ID= < Use access key id from devops >
    export AWS_SECRET_ACCESS_KEY= < Use Secret access key  from devops >
    export AWS_DEFAULT_REGION='us-east-1'

Step 2 : Assume IAM role
    aws sts assume-role --role-arn "arn:aws:iam::737606249163:role/iam_devops_role" --role-session-name AWSCLI-Session --duration-seconds 28800

Step 3 : Set the environment variables
    Use the out put value from command in step 2 above

```