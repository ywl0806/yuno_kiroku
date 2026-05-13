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
  memory                   = "4096" # 4 GB
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

      environment = [for k, v in var.app_env_vars : { name = k, value = v }]
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

resource "aws_security_group" "ecs_ai" {
  name        = "yuno-ecs-ai-${var.env}"
  description = "ECS AI Task security group"
  vpc_id      = var.vpc_id

  # ECR 이미지 pull, S3, SQS, CloudWatch Logs 등 AWS 서비스
  egress {
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
    description = "HTTPS"
  }

  # DB접근
  egress {
    from_port   = 5432
    to_port     = 5432
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
    description = "PostgreSQL (Supabase direct)"
  }
  egress {
    from_port   = 6543
    to_port     = 6543
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
    description = "PostgreSQL (Supabase connection pooler)"
  }

  tags = merge(var.common_tags, { Name = "yuno-ecs-ai-${var.env}" })
}

# ── ECS Service ──────────────────────────────────────────────

resource "aws_ecs_service" "ai" {
  name            = "yuno-ai-service-${var.env}"
  cluster         = aws_ecs_cluster.main.id
  task_definition = aws_ecs_task_definition.ai.arn
  desired_count   = 0

  capacity_provider_strategy {
    capacity_provider = "FARGATE_SPOT"
    weight            = 1
    base              = 0
  }

  network_configuration {
    subnets          = var.subnet_ids
    security_groups  = [aws_security_group.ecs_ai.id]
    assign_public_ip = true
  }

  lifecycle {
    ignore_changes = [desired_count]
  }

  tags = var.common_tags
}

# ── App Auto Scaling ───────────────────────────────────────

resource "aws_appautoscaling_target" "ai_task" {
  resource_id        = "service/${aws_ecs_cluster.main.name}/${aws_ecs_service.ai.name}"
  scalable_dimension = "ecs:service:DesiredCount"
  service_namespace  = "ecs"
  min_capacity       = 0
  max_capacity       = 10
}

# ── App Auto Scaling Policies ─────────────────────────────
# CloudWatch 알람은 modules/cloudwatch에서 관리

resource "aws_appautoscaling_policy" "ai_task" {
  name = "yuno-ai-task-scale-out-${var.env}"
  policy_type = "StepScaling"
  resource_id = aws_appautoscaling_target.ai_task.resource_id
  scalable_dimension = aws_appautoscaling_target.ai_task.scalable_dimension
  service_namespace  = "ecs"

  step_scaling_policy_configuration {
    # 정확한 태스크 수를 기준으로 스케일 아웃
    adjustment_type = "ExactCapacity"
    cooldown = 300
    # 최대 동시 실행 태스크 수를 기준으로 스케일 아웃
    metric_aggregation_type = "Maximum"

    step_adjustment {
      scaling_adjustment          = 1
      metric_interval_lower_bound = 0
      metric_interval_upper_bound = 20
    }

    step_adjustment {
      scaling_adjustment          = 2
      metric_interval_lower_bound = 20
      metric_interval_upper_bound = 40
    }

    step_adjustment {
      scaling_adjustment          = 3
      metric_interval_lower_bound = 40
    }
  }
}

resource "aws_appautoscaling_policy" "ai_task_scale_in" {
  name = "yuno-ai-task-scale-in-${var.env}"
  policy_type = "StepScaling"
  resource_id = aws_appautoscaling_target.ai_task.resource_id
  scalable_dimension = aws_appautoscaling_target.ai_task.scalable_dimension
  service_namespace  = "ecs"

  step_scaling_policy_configuration {
    adjustment_type = "ExactCapacity"
    cooldown = 300
    metric_aggregation_type = "Maximum"

    step_adjustment {
      scaling_adjustment = 0
      metric_interval_upper_bound = 0
    }
  }
}

# ── CloudWatch Log Group (Video Worker) ──────────────────────

resource "aws_cloudwatch_log_group" "video_task" {
  name              = "/ecs/yuno-video-${var.env}"
  retention_in_days = 30
  tags              = var.common_tags
}

# ── ECS Task Definition (Video Processing) ───────────────────

resource "aws_ecs_task_definition" "video" {
  family                   = "yuno-video-task-${var.env}"
  requires_compatibilities = ["FARGATE"]
  network_mode             = "awsvpc"
  cpu                      = "2048" # 2 vCPU
  memory                   = "4096" # 4 GB
  task_role_arn            = var.ecs_task_role_arn
  execution_role_arn       = var.ecs_task_execution_role_arn

  runtime_platform {
    operating_system_family = "LINUX"
    cpu_architecture        = "X86_64"
  }

  container_definitions = jsonencode([
    {
      name      = "yuno-video-worker"
      image     = var.video_image_uri
      essential = true

      environment = [for k, v in var.app_env_vars : { name = k, value = v }]
      logConfiguration = {
        logDriver = "awslogs"
        options = {
          "awslogs-group"         = aws_cloudwatch_log_group.video_task.name
          "awslogs-region"        = var.aws_region
          "awslogs-stream-prefix" = "ecs"
        }
      }
    }
  ])

  ephemeral_storage {
    size_in_gib = 50
  }

  tags = var.common_tags
}

# ── ECS Service (Video Processing) ───────────────────────────

resource "aws_ecs_service" "video" {
  name            = "yuno-video-service-${var.env}"
  cluster         = aws_ecs_cluster.main.id
  task_definition = aws_ecs_task_definition.video.arn
  desired_count   = 0

  capacity_provider_strategy {
    capacity_provider = "FARGATE_SPOT"
    weight            = 1
    base              = 0
  }

  network_configuration {
    subnets          = var.subnet_ids
    security_groups  = [aws_security_group.ecs_ai.id]
    assign_public_ip = true
  }

  lifecycle {
    ignore_changes = [desired_count]
  }

  tags = var.common_tags
}

# ── App Auto Scaling (Video) ──────────────────────────────────

resource "aws_appautoscaling_target" "video_task" {
  resource_id        = "service/${aws_ecs_cluster.main.name}/${aws_ecs_service.video.name}"
  scalable_dimension = "ecs:service:DesiredCount"
  service_namespace  = "ecs"
  min_capacity       = 0
  max_capacity       = 5
}

resource "aws_appautoscaling_policy" "video_task_scale_out" {
  name               = "yuno-video-task-scale-out-${var.env}"
  policy_type        = "StepScaling"
  resource_id        = aws_appautoscaling_target.video_task.resource_id
  scalable_dimension = aws_appautoscaling_target.video_task.scalable_dimension
  service_namespace  = "ecs"

  step_scaling_policy_configuration {
    adjustment_type         = "ExactCapacity"
    cooldown                = 600
    metric_aggregation_type = "Maximum"

    step_adjustment {
      scaling_adjustment          = 1
      metric_interval_lower_bound = 0
      metric_interval_upper_bound = 10
    }

    step_adjustment {
      scaling_adjustment          = 3
      metric_interval_lower_bound = 10
    }
  }
}

resource "aws_appautoscaling_policy" "video_task_scale_in" {
  name               = "yuno-video-task-scale-in-${var.env}"
  policy_type        = "StepScaling"
  resource_id        = aws_appautoscaling_target.video_task.resource_id
  scalable_dimension = aws_appautoscaling_target.video_task.scalable_dimension
  service_namespace  = "ecs"

  step_scaling_policy_configuration {
    adjustment_type         = "ExactCapacity"
    cooldown                = 600
    metric_aggregation_type = "Maximum"

    step_adjustment {
      scaling_adjustment          = 0
      metric_interval_upper_bound = 0
    }
  }
}