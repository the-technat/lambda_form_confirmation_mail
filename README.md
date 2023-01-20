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

## Deploy

1. Create lambda named `blabla` using the following settings:

- Runtime: `go1.x`
- Handler: `main`
- Function URL: yes, using Auth Type `NONE`
- Env:
  - SECRET=mySecret

2. Create a secrets in aws secretsmanager with the name of the `SECRET` env of the function. Add the following self-explaining keys:

- `MAIL_FROM`
- `MAIL_USER`
- `MAIL_PASSWORD`
- `MAIL_HOST`

3. Attach the `SecretsManagerReadWrite` to the execution role of your function
Note: this grants your function access to all secrets which is not what you want in production!

Finally upload the code using:

```bash
FUNCTION=blabla make deploy
```

Assuming that:

- you got a local go env
- you have the aws-cli installed and configured
