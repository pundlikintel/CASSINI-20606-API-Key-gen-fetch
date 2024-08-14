package aws

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
)

type CredentialProvider struct {
	accessKeyId, secretAccessKey, sessionToken string
}

func (cp *CredentialProvider) SetCredentials(accessKeyId, secretAccessKey, sessionToken string) {
	cp.secretAccessKey = secretAccessKey
	cp.accessKeyId = accessKeyId
	cp.sessionToken = sessionToken
}

func (cp *CredentialProvider) Retrieve(ctx context.Context) (aws.Credentials, error) {
	return aws.Credentials{
		AccessKeyID:     cp.accessKeyId,
		SecretAccessKey: cp.secretAccessKey,
		SessionToken:    cp.sessionToken,
	}, nil
}
