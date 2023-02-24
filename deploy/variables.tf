variable "region" {
  type        = string
  description = "AWS region you want to deploy the function in"
}
variable "resource_prefix" {
  type    = string
  default = "Name of the form you are using this function"
}
variable "mail_host" {
  type        = string
  description = "SMTP mail host"
}
variable "mail_from" {
  type        = string
  description = "FROM mail address"
}
variable "mail_user" {
  type        = string
  description = "User for the mail account"
}
variable "mail_pw" {
  type        = string
  sensitive   = true
  description = "Password for the mail account"
}
variable "mail_port" {
  type        = number
  default     = 465
  description = "SMTP por to use"
}
variable "mail_msg" {
  type        = string
  default     = "You're sign-up was successful"
  description = "You're Mail message, use {{ .user }} and {{ .content }} to render in the form content and username"
}
