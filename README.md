# lambda_form_confirmation_mail

![artifacts workflow](https://github.com/the-technat/lambda_form_confirmation_mail/actions/workflows/artifacts.yml/badge.svg)
![go version](https://img.shields.io/github/go-mod/go-version/the-technat/lambda_form_confirmation_mail)

Simple lambda written in Go that sends a mail using SMTP to a given mail with the content of the JSON that was submitted. Used to send confirmation mails for web forms that can do webhooks but don't implement confirmations mails on their own.

## Usage

Function expects a simple REST API call using POST and a JSON that contains at least the following fields:

```json
{
  "Form Title": "Registration for Event XYZ", // Mail subject line
  "E-Mail": "technat@technat.ch", // Receiver of the mail
  "Name": "Nathanael Liechti", // Greeted in mail
}
```

For example you could use the following curl command to trigger a new mail:

```bash
curl -X POST -H "content-type: application/json"  -d '{"Submission Date":"02.06.2016 10:23:54","Form Title":"Contact","Name":"Tim Schmitt","E-Mail":"technat@technat.ch","Phone":"0123/456789","Message":"Webhook-Formular-Submission!"}' https://f4sqdd35mf57m4msx3z3nr4c36priot.lambda-url.sa-east-1.on.aws
```

## Configuration

The lambda reads all configuration from a secret in AWS SecretsManager. The name of the secret is looked up from env var `SECRET`. The secret itself should have the following keys:

- `MAIL_FROM` -> From mail address
- `MAIL_USER` -> User for the mail account
- `MAIL_PASSWORD` -> Password for the mail account
- `MAIL_HOST` -> SMTP host
- `MAIL_PORT` -> Port of your SMTP host
- `MAIL_COPY` -> Whether you shall receive a copy of the mail
- `MAIL_MSG` -> Go template how your mail should look like. Must be formatted as HTML and can contain any number of keys from the JSON above
  - Example:
    ```html
    Hello {{ .Name }}

    You signed up successfully for the Event XYZ

    You're submited data:

    {{ .items }}

    Regards
    Event Team
    ```

### Function settings

Use the following settings when configuring a function:

- Runtime: `go1.x`
- Handler: `main`
- Architecture: `x86_64`
- Environment: `SECRET=nameOfYourSecret`
- Execution Policy: create new one using:
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
              "Action": [
                  "secretsmanager:GetSecretValue"
              ],
              "Resource": [
                  "arn:aws:secretsmanager:sa-east-1:298410952490:secret:id_of_secret"
              ]
          }
      ]
  }
  ```

## Deploy

You can take the [main.zip](./main.zip) and upload it to your function.

There is also a [container image](https://github.com/the-technat/lambda_form_confirmation_mail/pkgs/container/lambda_form_confirmation_mail) available if you want to go that way, you just have to push it to an ECR somewhere...

Furthermore theres a [Terraform module](./deploy) to deploy the function.
