package jwt

import (
	"context"
	"crypto/rsa"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/partnerkey/models"
	"strings"
	"time"
)

var pk *rsa.PrivateKey

type SecretsClient interface {
	GetSecretValue(ctx context.Context, params *secretsmanager.GetSecretValueInput, optFns ...func(options *secretsmanager.Options)) (*secretsmanager.GetSecretValueOutput, error)
}

func GenerateJwtKey(ctx context.Context, claims models.ApiKeyTags, domain, kid string, permissions models.Permission,
	scopes []string, secretsClient SecretsClient) (string, error) {
	token := jwt.New(jwt.SigningMethodPS384)
	issuedAt := time.Now()
	duration, _ := time.ParseDuration("3h")
	expiresAt := issuedAt.Add(duration)

	token.Claims = models.CustomClaims{
		RegisteredClaims: &jwt.RegisteredClaims{
			Issuer:    domain,
			Subject:   domain,
			NotBefore: &jwt.NumericDate{Time: time.Now()},
			Audience:  jwt.ClaimStrings{"amber.servicemesh"},
			ExpiresAt: &jwt.NumericDate{Time: expiresAt},
			IssuedAt:  &jwt.NumericDate{Time: issuedAt},
		},
		ApiKeyTags:  claims,
		Scope:       strings.Join(scopes, " "),
		Permissions: permissions,
	}

	if len(token.Header) == 0 {
		token.Header = map[string]interface{}{}
	}
	token.Header["kid"] = kid
	var err error
	if pk == nil {
		pk, err = getPk(ctx, secretsClient, kid)
		if err != nil {
			return "", err
		}
	}
	return token.SignedString(pk)
}

func getPk(ctx context.Context, secretsClient SecretsClient, keyId string) (*rsa.PrivateKey, error) {
	// AWS client for Secrets Manager
	//secretsClient := secretsmanager.NewFromConfig(cfg)

	secOp, err := secretsClient.GetSecretValue(ctx, &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(keyId),
		VersionStage: aws.String("AWSCURRENT"),
	})

	if err != nil || secOp == nil {
		//log.WithContext(ctx).WithError(err).Errorf("Not able to get private key kid: %s", keyId)
		return nil, err
	}
	signBytes := []byte(*secOp.SecretString)

	return jwt.ParseRSAPrivateKeyFromPEM(signBytes)
}
