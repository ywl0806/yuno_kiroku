output "cluster_arn" {
  value = aws_ecs_cluster.main.arn
}

output "cluster_name" {
  value = aws_ecs_cluster.main.name
}

output "task_definition_arn" {
  value = aws_ecs_task_definition.ai.arn
}

output "security_group_id" {
  value = aws_security_group.ecs_ai.id
}
