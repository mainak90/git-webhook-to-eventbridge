# Lambda function
data "aws_caller_identity" "current" {
}

resource "aws_iam_role" "golang_lambda_execution_role" {
  description        = "golang_lambda_execution_role"
  tags               = var.tags
  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "lambda.amazonaws.com"
      },
      "Effect": "Allow"
    }
  ]
}
EOF

}

resource "aws_iam_role_policy" "golang_lambda_execution_policy" {
  role   = aws_iam_role.golang_lambda_execution_role.id
  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
	    "events:PutEvents"
      ],
      "Resource": [
        "*"
      ]
    },
    {
        "Effect": "Allow",
        "Action": "logs:CreateLogGroup",
        "Resource": "arn:aws:logs:${var.region}:*:*"
    },
    {
        "Effect": "Allow",
        "Action": [
            "logs:CreateLogStream",
            "logs:PutLogEvents"
        ],
        "Resource": [
            "arn:aws:logs:${var.region}:${var.account_id}:*"
        ]
    }
  ]
}
EOF

}

resource "aws_lambda_function" "gitwebhook-putevents" {
  filename         = var.lambda_zip_path
  function_name    = var.lambda_function_name
  role             = aws_iam_role.golang_lambda_execution_role.arn
  handler          = var.lambda_handler
  source_code_hash = base64sha256(filebase64(var.lambda_zip_path))
  runtime          = var.lambda_runtime
  timeout          = var.lambda_timeout
  memory_size      = var.lambda_memory_size


  #vpc_config {
  #  security_group_ids = [
  #    aws_security_group.sg_lambda.id
  #  ]
  #  subnet_ids = var.lambda_subnet_ids
  #}

  tags = var.tags
}