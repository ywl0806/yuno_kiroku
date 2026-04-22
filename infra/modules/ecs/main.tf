# ── CloudWatch Log Group ─────────────────────────────────────

resource "aws_cloudwatch_log_group" "ai_task" {
  name              = "/ecs/yuno-ai-${var.env}"
  retention_in_days = 30
  tags              = var.common_tags
}

# ── ECS Cluster ──────────────────────────────────────────────

resource "aws_ecs_cluster" "main" {
  name = "yuno-cluster-${var.env}"

  setting {
    name  = "containerInsights"
    value = "disabled" # 비용 절감 (Phase 1)
  }

  tags = var.common_tags
}

resource "aws_ecs_cluster_capacity_providers" "main" {
  cluster_name = aws_ecs_cluster.main.name

  capacity_providers = ["FARGATE_SPOT", "FARGATE"]

  default_capacity_provider_strategy {
    capacity_provider = "FARGATE_SPOT"
    weight            = 1
    base              = 0
  }
}

# ── ECS Task Definition ──────────────────────────────────────

resource "aws_ecs_task_definition" "ai" {
  family                   = "yuno-ai-task-${var.env}"
  requires_compatibilities = ["FARGATE"]
  network_mode             = "awsvpc"
  cpu                      = "2048" # 2 vCPU
  memory                   = "8192" # 8 GB
  task_role_arn            = var.ecs_task_role_arn
  execution_role_arn       = var.ecs_task_execution_role_arn

  runtime_platform {
    operating_system_family = "LINUX"
    cpu_architecture        = "X86_64"
  }

  container_definitions = jsonencode([
    {
      name      = "yuno-ai"
      image     = var.ai_image_uri
      essential = true


      logConfiguration = {
        logDriver = "awslogs"
        options = {
          "awslogs-group"         = aws_cloudwatch_log_group.ai_task.name
          "awslogs-region"        = var.aws_region
          "awslogs-stream-prefix" = "ecs"
        }
      }
    }
  ])

  tags = var.common_tags
}

# ── Security Group (ECS AI Task) ─────────────────────────────
# Phase 1: VPC 없음 구조에서 default VPC 사용, assignPublicIp=ENABLED

resource "aws_security_group" "ecs_ai" {
  name        = "yuno-ecs-ai-${var.env}"
  description = "ECS AI Task security group"
  vpc_id      = var.vpc_id

  # Supabase DB, S3, ECR, Resize Lambda Function URL 접근
  egress {
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
    description = "HTTPS outbound (Supabase, S3, ECR, Lambda Function URL)"
  }

  tags = merge(var.common_tags, { Name = "yuno-ecs-ai-${var.env}" })
}
