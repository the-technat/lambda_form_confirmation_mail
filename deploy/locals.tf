locals {
  resource_name = format("%s%s", var.resource_prefix, "form_confirmation_mail")
  common_tags = {
    "name"    = local.resource_name
    "project" = "https://github.com/the-technat/lambda_form_confirmation_mail"
  }
  secret_data = {
    "MAIL_HOST" = var.mail_host
    "MAIL_FROM" = var.mail_from
    "MAIL_USER" = var.mail_user
    "MAIL_PW"   = var.mail_pw
    "MAIL_PORT" = var.mail_port
    "MAIL_MSG"  = var.mail_msg
  }
}
