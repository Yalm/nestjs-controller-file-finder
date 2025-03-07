# Bitbucket Pipelines Pipe: AWS API GATEWAY V1

This project provides a binary that can be used in a Bitbucket pipeline to interact with AWS API Gateway.


## YAML Definition

Add the following snippet to the `script` section of your `bitbucket-pipelines.yml` file:

```yaml
  - curl -L -o pipe https://github.com/Yalm/nestjs-controller-file-finder/releases/download/v0.0.1/pipe
  - chmod +x pipe
  - ./pipe
```

## Variables

| Variable                         | Usage                                                                                                                                                                                                                                                                        |
| -------------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| AWS_ACCESS_KEY_ID (\*\*)         | AWS access key id                                                                                                                                                                                                                                                            |
| AWS_SECRET_ACCESS_KEY (\*\*)     | AWS secret key                                                                                                                                                                                                                                                               |
| AWS_DEFAULT_REGION (\*\*)        | AWS region                                                                                                                                                                                                                                                                   |
| AWS_OIDC_ROLE_ARN                | The ARN of the role used for web identity federation or OIDC. See **Authentication**.                                                                                                                                                                                        |
| AWS_ROLE_ARN                     | Specifies the Amazon Resource Name (ARN) of an IAM role with a web identity provider that you want to use to run the AWS commands                                                                                                                                            |
| AWS_ROLE_SESSION_NAME            | Specifies the name to attach to the role session. This value is provided to the RoleSessionName parameter when the AWS CLI calls the AssumeRole operation, and becomes part of the assumed role user ARN: arn:aws:sts::123456789012:assumed-role/role_name/role_session_name |
| REST_API_ID\*| The id of the API 
| BACKEND_URL\* | The URL of the backend

_(\*) = required variable. This variable needs to be specified always when using the pipe._
_(\*\*) = required variable. If this variable is configured as a repository, account or environment variable, it doesnâ€™t need to be declared in the pipe as it will be taken from the context. It can still be overridden when using the pipe._

## Prerequisites

To use this pipe you should have a IAM user configured with programmatic access or Web Identity Provider (OIDC) role, with the necessary permissions to Authenticate on codeArtifact repository.

## Authentication

Supported options:

1. Environment variables: AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY. Default option.

2. Assume role provider with OpenID Connect (OIDC). More details in the Bitbucket Pipelines Using OpenID Connect guide [Integrating aws bitbucket pipeline with oidc][aws-oidc]. Make sure that you set up OIDC before:
   - configure Bitbucket Pipelines as a Web Identity Provider in AWS
   - attach to provider your AWS role with required policies in AWS
   - set up a build step with `oidc: true` in your Bitbucket Pipelines
   - pass AWS_OIDC_ROLE_ARN (\*) variable that represents role having appropriate permissions to execute actions on AWS CodeArtifact resources

### Basic example:

Authenticate with default options:

```yaml
- step:
    script:
      - curl -L -o pipe https://github.com/Yalm/nestjs-controller-file-finder/releases/download/0.0.1/pipe
      - chmod +x pipe
      - ./pipe
```