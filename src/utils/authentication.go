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

func defaultAuth(awsProfile, awsRegion string) apigateway.Options {
	if awsProfile != "" && awsProfile != "default" {
		cfg, err := config.LoadDefaultConfig(context.TODO(),
			config.WithSharedConfigProfile(awsProfile),
			config.WithRegion(awsRegion),
		)
		if err != nil {
			log.Printf("Error loading profile %s, falling back to environment variables: %v", awsProfile, err)
			return fallbackToEnvAuth(awsRegion)
		}
		return apigateway.Options{
			Credentials: cfg.Credentials,
			Region:      awsRegion,
		}
	}

	return fallbackToEnvAuth(awsRegion)
}

func fallbackToEnvAuth(awsRegion string) apigateway.Options {
	return apigateway.Options{
		Credentials: credentials.NewStaticCredentialsProvider(
			os.Getenv("AWS_ACCESS_KEY_ID"),
			os.Getenv("AWS_SECRET_ACCESS_KEY"),
			"",
		),
		Region: awsRegion,
	}
}

func NewSession(ctx context.Context, awsProfile, awsRegion string) apigateway.Options {
	roleArn := os.Getenv("AWS_OIDC_ROLE_ARN")
	webIdentityTokenFile := os.Getenv("BITBUCKET_STEP_OIDC_TOKEN")
	if webIdentityTokenFile == "" && roleArn == "" {
		return defaultAuth(awsProfile, awsRegion)
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

	return apigateway.Options{
		Credentials: credentials.NewStaticCredentialsProvider(
			*output.Credentials.AccessKeyId,
			*output.Credentials.SecretAccessKey,
			*output.Credentials.SessionToken,
		),
		Region: awsRegion,
	}
}
