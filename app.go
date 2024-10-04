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
	if err != nil || len(keys) == 0 {
		logrus.Errorf("No keys found, Crating")
		keys, err := api_key.CreateApiKey(ctx, jwt, *environment, *serviceId, *productId)
		if err != nil {
			fmt.Printf("Error in create api key %v", err)
			return
		}
		fmt.Printf("Key %v", keys)
		return
	}
	fmt.Printf("Key %v", keys)
}
