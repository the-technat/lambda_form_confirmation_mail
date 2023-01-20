# lambda_form_confirmation_mail

Simple lambda written in Go that sends a confirmation mail using AWS SES to a given mail with the content of the form that was submitted as json in an HTTP POST request.

## Docs

- <https://docs.aws.amazon.com/lambda/latest/dg/urls-invocation.html>

## Test

```bash
curl -X POST -H "content-type: application/json"  -d '{"Submission Date":"02.06.2016 10:23:54","Form Title":"Contact","Name":"Tim Schmitt","E-Mail":"test@beispiel.de","Phone":"0123/456789","Message":"Webhook-Formular-Submission!"}' https://f4sqdd35mf57m4msx3z3nr4c36priot.lambda-url.sa-east-1.on.aws
```

## Input Model

The following JSON is at least required to send a correct mail, all other fields are parsed as HTML table into the mail.

```json
{
  "Form Title": "Registration for Event XYZ", -> Mail subject line
  "E-Mail": "technat@technat.ch", -> Receiver of the mail
}
```

## Configuration

The following environment variables are needed for the function to know which SES to use:

```bash

```

## Deploy

Built as container image and deployed using [Terraform](https://github.com/alleaffengaffen/aws_baseline/blob/main/lambda.tf)
