package utils

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/apigateway"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

func defaultAuth() *apigateway.Options {
	return &apigateway.Options{
		Credentials: credentials.NewStaticCredentialsProvider(
			os.Getenv("AWS_ACCESS_KEY_ID"),
			os.Getenv("AWS_SECRET_ACCESS_KEY"),
			"",
		),
		Region: os.Getenv("AWS_DEFAULT_REGION"),
	}
}

func NewSession(ctx context.Context) *apigateway.Options {
	roleArn := os.Getenv("AWS_OIDC_ROLE_ARN")
	webIdentityTokenFile := os.Getenv("BITBUCKET_STEP_OIDC_TOKEN")
	if webIdentityTokenFile == "" && roleArn == "" {
		return defaultAuth()
	}

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatalf("Error loading base configuration: %v", err)
	}

	stsClient := sts.NewFromConfig(cfg)

	output, err := stsClient.AssumeRoleWithWebIdentity(ctx, &sts.AssumeRoleWithWebIdentityInput{
		RoleArn:          aws.String(roleArn),
		RoleSessionName:  aws.String("build-session"),
		WebIdentityToken: aws.String(webIdentityTokenFile),
	})

	if err != nil {
		log.Fatalf("Error assuming role: %v", err)
	}

	return &apigateway.Options{
		Credentials: credentials.NewStaticCredentialsProvider(
			*output.Credentials.AccessKeyId,
			*output.Credentials.SecretAccessKey,
			*output.Credentials.SessionToken,
		),
		Region: os.Getenv("AWS_DEFAULT_REGION"),
	}
}
