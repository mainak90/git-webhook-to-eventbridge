resource "aws_apigatewayv2_api" "api" {
  name          = "gitwebhook-put-event"
  protocol_type = "HTTP"
  cors_configuration {
    allow_methods = ["*"]
    allow_origins = ["*"]
  }
  target = aws_lambda_function.findall.arn
}

resource "aws_lambda_permission" "apigw" {
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.findall.arn
  principal     = "apigateway.amazonaws.com"

  source_arn = "${aws_apigatewayv2_api.api.execution_arn}/*/*"
}

resource "aws_apigatewayv2_api" "api2" {
  name          = "api-insert"
  protocol_type = "HTTP"
  cors_configuration {
    allow_methods = ["POST", "PUT"]
    allow_origins = ["*"]
  }
  target = aws_lambda_function.insertItem.arn
}

resource "aws_lambda_permission" "apigw2" {
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.insertItem.arn
  principal     = "apigateway.amazonaws.com"

  source_arn = "${aws_apigatewayv2_api.api2.execution_arn}/*/*"
}
