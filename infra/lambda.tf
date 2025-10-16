resource "aws_lambda_function" "api_lambda" {
  function_name = "api-lambda"
  role          = aws_iam_role.api_lambda_role.arn
  handler       = "cmd/server/main.main"
  filename      = "cmd/server/main"
  source_code_hash = filebase64sha256("cmd/server/main")
}

