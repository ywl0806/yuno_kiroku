output "api_lambda_role_arn" {
    value = aws_iam_role.api_lambda.arn
}

output "resize_lambda_role_arn" {
    value = aws_iam_role.resize_lambda.arn
}

output "ecs_task_execution_role_arn" {
    value = aws_iam_role.ecs_task_execution.arn
}

output "ecs_task_role_arn" {
    value = aws_iam_role.ecs_task.arn
}