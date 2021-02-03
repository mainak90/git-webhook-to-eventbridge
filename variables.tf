variable "region" {
  description = "region"
  default     = "us-east-1"
}

variable "component" {
  default = "golang"
}
variable "deployment_identifier" {
  default = "gitwebhook-putevents"
}

variable "account_id" {
  description = "AWS account id where the lambda execution"
  default     = ""
}

variable "vpc_id" {
  description = "VPC to deploy the lambda to"
  default     = ""
}

variable "lambda_subnet_ids" {
  description = "Subnet ids to deploy the lambda to"
  type        = list(string)
  default     = ["subnet-0160bb6305b46eaee"]
}


variable "api_gateway_stage_name" {
  default = "staging"
}

variable "resource_path_part" {
  default = "gitwebhook-putevents"
}

variable "lambda_zip_path" {
  default = "build/deployment-git-webhook.zip"
}

variable "lambda_ingress_cidr_blocks" {
  type    = list(string)
  default = ["0.0.0.0/0"]
}

variable "lambda_egress_cidr_blocks" {
  type    = list(string)
  default = ["0.0.0.0/0"]
}


variable "lambda_function_name" {
  default = "gitwebhook-putevents"
}

variable "lambda_handler" {
  default = "main"
}

variable "lambda_runtime" {
  default = "go1.x"
}

variable "lambda_timeout" {
  default = 30
}

variable "lambda_memory_size" {
  default = 128
}

variable "tags" {
  default = {}
}
