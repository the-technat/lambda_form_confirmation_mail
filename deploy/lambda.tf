resource "aws_iam_policy" "execution_policy" {
  name        = local.resource_name
  path        = "/"
  description = "Policy for aws lamda ${local.resource_name}"

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
          aws_secretsmanager_secret.secret.arn
        ]
      }
    ]
  })

  tags = local.common_tags
}

resource "aws_iam_role" "execution_role" {
  name = local.resource_name
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

resource "aws_iam_role_policy_attachment" "execution_binding" {
  role       = aws_iam_role.execution_role.name
  policy_arn = aws_iam_policy.execution_policy.arn
}

resource "null_resource" "lambda_code" {
  triggers = {
    always_run = timestamp() # this will always run
  }

  provisioner "local-exec" {
    command = "wget https://github.com/the-technat/lambda_form_confirmation_mail/releases/download/${var.lambda_archive_version}/main.zip"
  }
}

resource "aws_lambda_function" "fcm" {
  function_name = local.resource_name
  role          = aws_iam_role.execution_role.arn

  filename = "${path.module}/main.zip"

  runtime = "go1.x"
  handler = "main"
  timeout = "60"

  environment {
    variables = {
      SECRET = "${local.resource_name}-${random_string.random.result}"
    }
  }

  depends_on = [
    null_resource.lambda_code
  ]
}

resource "aws_lambda_function_url" "webhook" {
  function_name      = aws_lambda_function.fcm.function_name
  authorization_type = "NONE"
}
