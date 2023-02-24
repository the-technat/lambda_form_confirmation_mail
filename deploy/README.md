# terraform-aws-lambda-confirmation-mail

TF Modul to deploy the lambda

## Usage

```hcl
module "fcm_function" {
  source = "git::https://github.com/the-technat/lambda_form_confirmation_mail?ref=main"

  region = "sa-east-1"
  resource_prefix = "event_xyz"

  mail_user = "banane@alleaffengaffen.ch"
  mail_pw = "SuperSecure123"
  mail_host = "mail.gmail.com"
  mail_port = 587
  mail_from = "technat@technat.ch"
}

```

<!-- BEGIN_TF_DOCS -->
## Requirements

No requirements.

## Providers

| Name | Version |
|------|---------|
| <a name="provider_aws"></a> [aws](#provider\_aws) | 4.56.0 |
| <a name="provider_null"></a> [null](#provider\_null) | 3.2.1 |
| <a name="provider_random"></a> [random](#provider\_random) | 3.4.3 |

## Modules

No modules.

## Resources

| Name | Type |
|------|------|
| [aws_iam_policy.execution_policy](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/iam_policy) | resource |
| [aws_iam_role.execution_role](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/iam_role) | resource |
| [aws_iam_role_policy_attachment.execution_binding](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/iam_role_policy_attachment) | resource |
| [aws_lambda_function.fcm](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/lambda_function) | resource |
| [aws_lambda_function_url.webhook](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/lambda_function_url) | resource |
| [aws_secretsmanager_secret.secret](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/secretsmanager_secret) | resource |
| [aws_secretsmanager_secret_version.secret_data](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/secretsmanager_secret_version) | resource |
| [null_resource.lambda_code](https://registry.terraform.io/providers/hashicorp/null/latest/docs/resources/resource) | resource |
| [random_string.random](https://registry.terraform.io/providers/hashicorp/random/latest/docs/resources/string) | resource |
| [aws_caller_identity.current](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/data-sources/caller_identity) | data source |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_mail_from"></a> [mail\_from](#input\_mail\_from) | FROM mail address | `string` | n/a | yes |
| <a name="input_mail_host"></a> [mail\_host](#input\_mail\_host) | SMTP mail host | `string` | n/a | yes |
| <a name="input_mail_msg"></a> [mail\_msg](#input\_mail\_msg) | You're Mail message, use {{ .user }} and {{ .content }} to render in the form content and username | `string` | `"You're sign-up was successful"` | no |
| <a name="input_mail_port"></a> [mail\_port](#input\_mail\_port) | SMTP por to use | `number` | `465` | no |
| <a name="input_mail_pw"></a> [mail\_pw](#input\_mail\_pw) | Password for the mail account | `string` | n/a | yes |
| <a name="input_mail_user"></a> [mail\_user](#input\_mail\_user) | User for the mail account | `string` | n/a | yes |
| <a name="input_region"></a> [region](#input\_region) | AWS region you want to deploy the function in | `string` | n/a | yes |
| <a name="input_resource_prefix"></a> [resource\_prefix](#input\_resource\_prefix) | n/a | `string` | `"Name of the form you are using this function"` | no |

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_function_url"></a> [function\_url](#output\_function\_url) | n/a |
<!-- END_TF_DOCS -->