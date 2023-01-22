# lambda_form_confirmation_mail

![release workflow](https://github.com/the-technat/lambda_form_confirmation_mail/actions/workflows/release.yml/badge.svg)

Simple lambda written in Go that sends a confirmation mail using SMTP to a given mail with the content of the form that was submitted.

## Usage

Function expects a simple REST API call using POST and a JSON that contains at least the following fields:

```json
{
  "Form Title": "Registration for Event XYZ", -> Mail subject line
  "E-Mail": "technat@technat.ch", -> Receiver of the mail
}
```

For example you could use the following curl command:

```bash
curl -X POST -H "content-type: application/json"  -d '{"Submission Date":"02.06.2016 10:23:54","Form Title":"Contact","Name":"Tim Schmitt","E-Mail":"technat@technat.ch","Phone":"0123/456789","Message":"Webhook-Formular-Submission!"}' https://f4sqdd35mf57m4msx3z3nr4c36priot.lambda-url.sa-east-1.on.aws
```

## Configuration

The lambda reads all his configuration from a secret in AWS SecretsManager. The name of the secret is lookup from env var `SECRET`. The secret itself should have the following keys:

- `MAIL_FROM`
- `MAIL_USER`
- `MAIL_PASSWORD`
- `MAIL_HOST`

### Function settings

- Runtime: `go1.x`
- Handler: `main`
- Architecture: `x86_64`
- Environment: `SECRET=nameOfYourSecret`

### Permissions

The lambda needs the following policy in it's execution role:

```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": "logs:CreateLogGroup",
            "Resource": "arn:aws:logs:sa-east-1:298300902191:*"
        },
        {
            "Effect": "Allow",
            "Action": [
                "logs:CreateLogStream",
                "logs:PutLogEvents"
            ],
            "Resource": [
                "arn:aws:logs:sa-east-1:298410952490:log-group:/aws/lambda/form_confirmation_mail:*"
            ]
        },
        {
            "Effect": "Allow",
            "Actions": [
                "secretsmanager:GetSecretValue"
            ],
            "Resource": [
                "arn:aws:secretsmanager:sa-east-1:298410952490:id_of_secret"
            ]
        }
    ]
}
```

## Deploy

For every tag that is manually pushed to the repo, a GH action packages the code as zip file. See the latest release for the zip archive you can upload to your function.

There is also a [container image](https://github.com/the-technat/lambda_form_confirmation_mail/pkgs/container/lambda_form_confirmation_mail) available to download for each release if you want to go that way, you just have to push it an ECR somewhere...

Furthermore checkout this Terraform file to deploy the function and secret: <https://github.com/alleaffengaffen/aws_baseline/blob/main/lambda_form_confirmation_mail.tf>
