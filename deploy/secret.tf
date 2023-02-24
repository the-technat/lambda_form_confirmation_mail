resource "random_string" "random" {
  length  = 8
  special = false
}

resource "aws_secretsmanager_secret" "secret" {
  name                    = "${local.resource_name}-${random_string.random.result}"
  recovery_window_in_days = 0

  tags = local.common_tags
}

resource "aws_secretsmanager_secret_version" "secret_data" {
  secret_id     = aws_secretsmanager_secret.secret.id
  secret_string = jsonencode(local.secret_data)
}
