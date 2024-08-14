package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	api_key "github.com/partnerkey/api-key"
	aws2 "github.com/partnerkey/aws"
	"github.com/partnerkey/jwt"
	"github.com/partnerkey/models"
	"github.com/sirupsen/logrus"
)

func main() {
	var ctx context.Context = context.Background()
	/*domain := "api-dev01-user7.project-amber-smas.com"
	kid := "amber-apigw-jwt-signing-us-east-1-dev01"
	var AccessKeyId, SecretAccessKey, sessionToken = "ASIA2XPGCFLFWT4MNKMM", "SoF8Nb5uyuSW/G4xE9+7n1CWiUCk1Wb70GikSVUK", "IQoJb3JpZ2luX2VjEMX//////////wEaCXVzLWVhc3QtMSJHMEUCIQD1ozSRb8hrEFyaLJkOfU4kNglyDsKfPkJWLCoMbJ+cTQIgTrgn0Lgn6KjGopMSvBN4E06LsSBnI0oCrXOABxtebJ8qpQIIvf//////////ARABGgw3Mzc2MDYyNDkxNjMiDCpMJKuuWPVVbbC8ZCr5ATMeB5sbNYAF/wkdwfYA+WLAGwDHZ8zvTGnhKRDxP6LkEVtqllkqeyTyh2LhdgFoaeFrxcZlgFJuyw93lg2BbBLVo+U1Od0F/C2zk9isl4KVpzSqv9ECbMG1BaBjdzjYB8G+CV05kMtoxO21afY2lnPw0+xYyYjx8Itz8u98CoyFv/xu7jKOTz1HkiLWsnivnOOBpHkt3Ya+w+ZkBiFoR5KvTBlVuTkdniSTacNKz0dHVllrPmNhOgOvyvYNFu99xv4dp3qgHwb1Hj6mjru4afvrFUzcFEMOm+OzldFJsNVHMQ969cX+2XR2rPykP+2J+7ixl70HU7srETCgwfK1BjqdAcvI8RzzO4gA4nRRJWk5boCqXIos+irviOJJk0J/VulPm/+x2N9xmiLNvBUo+EoV4rQ6UbbzU0P7l26iK22Of8qJlV+i0Ims3J8Q9g/kRzzkxC0SnmK7bwRTov1RHL/2RolJ+aXr+CL+1v6ihTpOSw2IW7HLrSk+QMsGuKMtUbWtjw7nTaVYoQInT4DBdfpSDcuTBlOyj7qD66iU098="
	awsRegion := "us-east-1"
	tenantId := "e6711eab-3416-4f66-b033-6485ae3d40f6"
	UserId := "62ef5530-6c33-4d9a-9816-96c2041647bc"
	environment := "dev01-user7"
	productId := "827be883-7219-4910-aa36-c17e023a5f7d"
	serviceId := "840efc64-3c4b-41f2-8b9e-e1c70781331c"
	*/
	domain := flag.String("issuer", "", "Issuer like api-dev01-user7.project-amber-smas.com")
	kid := flag.String("sec_kid", "", "secret key id")
	AccessKeyId := flag.String("aws_key_id", "", "Access Key Id")
	SecretAccessKey := flag.String("aws_secret_key", "", "Secret Access Key")
	sessionToken := flag.String("aws_token", "", "Session token")
	awsRegion := flag.String("region", "us-east-1", "AWS region")
	tenantId := flag.String("tenant_id", "", "Tenant id")
	UserId := flag.String("user_id", "", "user id")
	environment := flag.String("environment", "", "Full env like dev01-user1")
	productId := flag.String("product_id", "", "Product id")
	serviceId := flag.String("service_id", "", "Service id")

	flag.Parse()

	apikeyTags := models.ApiKeyTags{
		TenantId: tenantId,
		UserRole: nil,
		UserId:   UserId,
	}
	///////////////////////
	permissionJson := []byte(`{
    "partner-customer": {
      "grants": [
        "read",
        "write",
        "modify",
        "remove"
      ]
    },
    "service": {
      "grants": [
        "write",
        "modify",
        "read",
        "remove"
      ],
      "data": {
        "productType": [
          "partner"
        ]
      }
    }
  }`)

	//////////////////////
	permission := models.Permission{}
	err := json.Unmarshal(permissionJson, &permission)

	if err != nil {
		_ = err
	}
	scopes := []string{"partner", "tenant"}

	credProvider := &aws2.CredentialProvider{}
	//TBD: Credential provider need to change based on discussion
	credProvider.SetCredentials(*AccessKeyId, *SecretAccessKey, *sessionToken)

	conf := aws.Config{Credentials: credProvider, Region: *awsRegion, RetryMaxAttempts: 10, RetryMode: aws.RetryModeStandard}
	secretsClient := secretsmanager.NewFromConfig(conf)

	jwt, err := jwt.GenerateJwtKey(ctx, apikeyTags, *domain, *kid, permission, scopes, secretsClient)

	if err != nil {
		logrus.Errorf("Error in generate jwt %v", err)
	}
	logrus.Infof("JWT: %v", jwt)
	keys, err := api_key.GetApiKeys(ctx, jwt, *environment, *serviceId)
	if err != nil {
		return
	}
	if len(keys) == 0 {
		logrus.Errorf("No keys found, Crating")
		keys, err := api_key.CreateApiKey(ctx, jwt, *environment, *serviceId, *productId)
		if err != nil {
			fmt.Printf("Error in create api key %v", err)
			return
		}
		fmt.Printf("Key %v", keys)
	}
	fmt.Printf("Key %v", keys)
}
