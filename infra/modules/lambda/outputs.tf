output "api_lambda_function_arn" {
  value = aws_lambda_function.api.arn
}

output "api_lambda_function_name" {
  value = aws_lambda_function.api.function_name
}

output "api_lambda_function_invoke_arn" {
  value = aws_lambda_function.api.invoke_arn
}

output "resize_lambda_function_url" {
  value = aws_lambda_function_url.resize.function_url
}