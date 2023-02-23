##############
# Variables
# Passed using TFC, CI/CD or CLI
##############
variable "fcm_lambda_archive_version" {
  type    = string
  default = "v0.0.2"
}
variable "region" {
  type = string
}
variable "resource_prefix" {
  type = string
}
variable "mail_host" {
  type = string
}
variable "mail_from" {
  type = string
}
variable "mail_user" {
  type = string
}
variable "mail_pw" {
  type      = string
  sensitive = true
}
variable "mail_port" {
  type = number
}
variable "mail_msg" {
  type = string
}

##############
# Terraform Config
##############
terraform {
  required_providers {
    random = {
      source = "hashicorp/random"
    }
    null = {
      source = "hashicorp/null"
    }
    aws = {
      source = "hashicorp/aws"
    }
  }
}

##############
# Locals
# Internally used
##############
locals {
  fcm_resource_name = format("%s%s", var.resource_prefix, "form_confirmation_mail")
  common_tags = {
    "name"    = local.fcm_resource_name
    "project" = "https://github.com/the-technat/lambda_form_confirmation_mail"
  }
  fcm_secret_data = {
    "MAIL_HOST" = var.mail_host
    "MAIL_FROM" = var.mail_from
    "MAIL_USER" = var.mail_user
    "MAIL_PW"   = var.mail_pw
    "MAIL_PORT" = var.mail_port
    "MAIL_MSG"  = var.mail_msg
  }
}

##############
# Data
##############
data "aws_caller_identity" "current" {}

##############
# Outputs
# can be used
##############
output "fcm_webhook" {
  value = aws_lambda_function_url.fcm_webhook.function_url
}

##############
# Resources
##############
resource "random_string" "random" {
  length  = 8
  special = false
}

resource "aws_secretsmanager_secret" "fcm_secret" {
  name                    = "${local.fcm_resource_name}-${random_string.random.result}"
  recovery_window_in_days = 0

  tags = local.common_tags
}

resource "aws_secretsmanager_secret_version" "fcm_secret_data" {
  secret_id     = aws_secretsmanager_secret.fcm_secret.id
  secret_string = jsonencode(local.fcm_secret_data)
}

resource "aws_iam_policy" "fcm_execution_policy" {
  name        = local.fcm_resource_name
  path        = "/"
  description = "Policy for aws lamda ${local.fcm_resource_name}"

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = [
          "logs:CreateLogGroup",
          "logs:CreateLogStream",
          "logs:PutLogEvents",
        ],
        Effect   = "Allow"
        Resource = "arn:aws:logs:${var.region}:${data.aws_caller_identity.current.account_id}:*"
      },
      {
        Action = [
          "secretsmanager:GetSecretValue"
        ]
        Effect = "Allow"
        Resource = [
          aws_secretsmanager_secret.fcm_secret.arn
        ]
      }
    ]
  })

  tags = local.common_tags
}

resource "aws_iam_role" "fcm_execution_role" {
  name = local.fcm_resource_name
  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Sid    = ""
        Principal = {
          Service = "lambda.amazonaws.com"
        }
      },
    ]
  })
  tags = local.common_tags
}

resource "aws_iam_role_policy_attachment" "fcm_execution_binding" {
  role       = aws_iam_role.fcm_execution_role.name
  policy_arn = aws_iam_policy.fcm_execution_policy.arn
}

resource "null_resource" "fcm_lambda_code" {
  triggers = {
    always_run = timestamp() # this will always run
  }

  provisioner "local-exec" {
    command = "wget https://github.com/the-technat/lambda_form_confirmation_mail/releases/download/${var.fcm_lambda_archive_version}/main.zip"
  }
}

resource "aws_lambda_function" "fcm" {
  function_name = local.fcm_resource_name
  role          = aws_iam_role.fcm_execution_role.arn

  filename = "${path.module}/main.zip"

  runtime = "go1.x"
  handler = "main"
  timeout = "60"

  environment {
    variables = {
      SECRET = "${local.fcm_resource_name}-${random_string.random.result}"
    }
  }

  depends_on = [
    null_resource.fcm_lambda_code
  ]
}

resource "aws_lambda_function_url" "fcm_webhook" {
  function_name      = aws_lambda_function.fcm.function_name
  authorization_type = "NONE"
}
