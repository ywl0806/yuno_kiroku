

# ── AssumeRole 정책 ──────────────────────────────────────────

data "aws_iam_policy_document" "lambda_assume" {
    statement {
        effect = "Allow"
        actions = ["sts:AssumeRole"]
        principals {
            type = "Service"
            identifiers = ["lambda.amazonaws.com"]
        }
    }
}

data "aws_iam_policy_document" "ecs_task_assume" {
    statement {
        effect = "Allow"
        actions = ["sts:AssumeRole"]
        principals {
            type = "Service"
            identifiers = ["ecs-tasks.amazonaws.com"]
        }
    }
}

# ── API Lambda Role ──────────────────────────────────────────

resource "aws_iam_role" "api_lambda" {
    name = "yuno-api-lambda-role-${var.env}"
    assume_role_policy = data.aws_iam_policy_document.lambda_assume.json
    tags = var.common_tags
}

resource "aws_iam_role_policy" "api_lambda_s3" {
    name = "s3-media-access"
    role = aws_iam_role.api_lambda.id
    policy = jsonencode({
        Version = "2012-10-17"
        Statement = [
            {
                Effect = "Allow"
                Action = ["s3:GetObject", "s3:PutObject", "s3:DeleteObject"]
                Resource = "${var.media_bucket_arn}/*"
            },
            {
                Effect = "Allow"
                Action = ["s3:ListBucket"]
                Resource = var.media_bucket_arn
            }
        ]
    })
}

resource "aws_iam_role_policy_attachment" "api_lambda_logs" {
    role = aws_iam_role.api_lambda.name
    policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

# ── Resize Lambda Role ──────────────────────────────────────────

resource "aws_iam_role" "resize_lambda" {
    name = "yuno-resize-lambda-role-${var.env}"
    assume_role_policy = data.aws_iam_policy_document.lambda_assume.json
    tags = var.common_tags
}

resource "aws_iam_role_policy" "resize_lambda_s3" {
    name = "s3-media-access"
    role = aws_iam_role.resize_lambda.id
    policy = jsonencode({
        Version = "2012-10-17"
        Statement = [
            {
                Effect = "Allow"
                Action = ["s3:GetObject", "s3:PutObject", "s3:DeleteObject"]
                Resource = "${var.media_bucket_arn}/*"
            },
            {
                Effect = "Allow"
                Action = ["s3:ListBucket"]
                Resource = var.media_bucket_arn

            }
        ]
    })
}
resource "aws_iam_role_policy" "resize_lambda_ecs" {
  name = "sqs-face-recognition-access"
  role = aws_iam_role.resize_lambda.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect   = "Allow"
        Action   = ["sqs:SendMessage"]
        Resource = var.face_recognition_queue_arn
      }
    ]
  })
}

resource "aws_iam_role_policy_attachment" "resize_lambda_logs" {
  role       = aws_iam_role.resize_lambda.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

# ── github actions ci/cd ──────────────────────────────────────────

resource "aws_iam_role" "github_actions_ci_cd" {
    name = "yuno-github-actions-ci-cd-role-${var.env}"
    assume_role_policy = data.aws_iam_policy_document.lambda_assume.json
    tags = var.common_tags
}

resource "aws_iam_role_policy" "github_actions_ci_cd_lambda" {

    name = "github-actions-ci-cd-lambda"
    role = aws_iam_role.github_actions_ci_cd.id
    policy = jsonencode({
        Version = "2012-10-17"
        Statement = [
            {
                Effect = "Allow"
                Action = [
                    "lambda:UpdateFunctionCode",
                    "lambda:GetFunction"
                ]
                Resource = [
                    "arn:aws:lambda:${var.aws_region}:${var.aws_account_id}:function:yuno-api-${var.env}",
                    "arn:aws:lambda:${var.aws_region}:${var.aws_account_id}:function:yuno-resize-${var.env}",
                ]
            },
            {
                Effect = "Allow"
                Action = [
                    "s3:GetObject",
                    "s3:GetObjectVersion",
                    "s3:PutObject"
                ],
                Resource = [
                    "${var.lambda_zip_bucket_arn}/*"
                ]
            }
        ]
    })
}

# ── ECS Task Role (AI 컨테이너가 사용) ───────────────────────

resource "aws_iam_role" "ecs_task" {
  name               = "yuno-ecs-task-role-${var.env}"
  assume_role_policy = data.aws_iam_policy_document.ecs_task_assume.json
  tags               = var.common_tags
}

resource "aws_iam_role_policy" "ecs_task_s3" {
  name = "s3-media-access"
  role = aws_iam_role.ecs_task.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect   = "Allow"
        Action   = ["s3:GetObject", "s3:PutObject"]
        Resource = "${var.media_bucket_arn}/*"
      },
      {
        Effect   = "Allow"
        Action   = ["s3:ListBucket"]
        Resource = var.media_bucket_arn
      },
      {
        Effect   = "Allow"
        Action   = [
          "sqs:ReceiveMessage",
          "sqs:DeleteMessage",
          "sqs:GetQueueAttributes"
        ]
        Resource = [
          var.face_recognition_queue_arn,
          var.video_queue_arn,
        ]
      }
    ]
  })
}

resource "aws_iam_role_policy" "resize_lambda_video_sqs" {
  name = "sqs-video-send"
  role = aws_iam_role.resize_lambda.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect   = "Allow"
        Action   = ["sqs:SendMessage"]
        Resource = var.video_queue_arn
      }
    ]
  })
}

resource "aws_iam_role_policy" "resize_lambda_sqs_consume" {
  name = "sqs-resize-consume"
  role = aws_iam_role.resize_lambda.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "sqs:ReceiveMessage",
          "sqs:DeleteMessage",
          "sqs:GetQueueAttributes"
        ]
        Resource = var.resize_queue_arn
      }
    ]
  })
}

# ── ECS Task Execution Role (ECS 에이전트가 사용) ────────────

resource "aws_iam_role" "ecs_task_execution" {
  name               = "yuno-ecs-task-execution-role-${var.env}"
  assume_role_policy = data.aws_iam_policy_document.ecs_task_assume.json
  tags               = var.common_tags
}

resource "aws_iam_role_policy_attachment" "ecs_task_execution_base" {
  role       = aws_iam_role.ecs_task_execution.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy"
}

resource "aws_iam_role_policy" "ecs_task_execution_ecr" {
  name = "ecr-pull"
  role = aws_iam_role.ecs_task_execution.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "ecr:GetDownloadUrlForLayer",
          "ecr:BatchGetImage",
          "ecr:BatchCheckLayerAvailability",
        ]
        Resource = var.ecr_repository_arns
      },
      {
        Effect   = "Allow"
        Action   = ["ecr:GetAuthorizationToken"]
        Resource = "*"
      }
    ]
  })
}

resource "aws_iam_role_policy" "ecs_task_execution_ssm" {
  name = "ssm-read"
  role = aws_iam_role.ecs_task_execution.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect   = "Allow"
        Action   = ["ssm:GetParameters"]
        Resource = "arn:aws:ssm:${var.aws_region}:${var.aws_account_id}:parameter/yuno/*"
      }
    ]
  })
}
