output "repository_urls" {
  description = "ECR repository URLs (key: repo name)"
  value       = { for k, v in aws_ecr_repository.repos : k => v.repository_url }
}

output "repository_arns" {
  description = "ECR repository ARNs (key: repo name)"
  value       = { for k, v in aws_ecr_repository.repos : k => v.arn }
}
