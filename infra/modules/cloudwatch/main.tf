# ── SNS 알림 인프라 ───────────────────────────────────────────

resource "aws_sns_topic" "alerts" {
  name = "yuno-alerts-${var.env}"
  tags = var.common_tags
}

resource "aws_sns_topic_subscription" "email" {
  topic_arn = aws_sns_topic.alerts.arn
  protocol  = "email"
  endpoint  = var.alert_email
}

# ── DLQ 알람 ─────────────────────────────────────────────────

resource "aws_cloudwatch_metric_alarm" "resize_dlq" {
  alarm_name          = "yuno-resize-dlq-alarm-${var.env}"
  comparison_operator = "GreaterThanThreshold"
  evaluation_periods  = 1
  metric_name         = "ApproximateNumberOfMessagesVisible"
  namespace           = "AWS/SQS"
  period              = 60
  statistic           = "Sum"
  threshold           = 0
  alarm_actions       = [aws_sns_topic.alerts.arn]

  dimensions = {
    QueueName = var.resize_dlq_name
  }
}

resource "aws_cloudwatch_metric_alarm" "face_recognition_dlq" {
  alarm_name          = "yuno-face-recognition-dlq-alarm-${var.env}"
  comparison_operator = "GreaterThanThreshold"
  evaluation_periods  = 1
  metric_name         = "ApproximateNumberOfMessagesVisible"
  namespace           = "AWS/SQS"
  period              = 60
  statistic           = "Sum"
  threshold           = 0
  alarm_actions       = [aws_sns_topic.alerts.arn]

  dimensions = {
    QueueName = var.face_recognition_dlq_name
  }
}

resource "aws_cloudwatch_metric_alarm" "video_dlq" {
  alarm_name          = "yuno-video-processing-dlq-alarm-${var.env}"
  comparison_operator = "GreaterThanThreshold"
  evaluation_periods  = 1
  metric_name         = "ApproximateNumberOfMessagesVisible"
  namespace           = "AWS/SQS"
  period              = 60
  statistic           = "Sum"
  threshold           = 0
  alarm_actions       = [aws_sns_topic.alerts.arn]

  dimensions = {
    QueueName = var.video_dlq_name
  }
}

# ── API 에러 로그 감시 ────────────────────────────────────────

resource "aws_cloudwatch_log_metric_filter" "api_error" {
  name           = "yuno-api-error-filter-${var.env}"
  log_group_name = var.api_log_group_name
  pattern        = "{ $.level = \"ERROR\" }"

  metric_transformation {
    name      = "ApiErrorCount-${var.env}"
    namespace = "YunoApp"
    value     = "1"
  }
}

resource "aws_cloudwatch_metric_alarm" "api_error" {
  alarm_name          = "yuno-api-error-alarm-${var.env}"
  comparison_operator = "GreaterThanOrEqualToThreshold"
  evaluation_periods  = 1
  metric_name         = "ApiErrorCount-${var.env}"
  namespace           = "YunoApp"
  period              = 300
  statistic           = "Sum"
  threshold           = 10
  alarm_actions       = [aws_sns_topic.alerts.arn]
}

# ── ECS 스케일링 알람 ─────────────────────────────────────────

resource "aws_cloudwatch_metric_alarm" "ai_task_scale_out" {
  alarm_name = "yuno-ai-task-scale-out-${var.env}"
  comparison_operator = "GreaterThanOrEqualToThreshold"
  evaluation_periods  = 1
  metric_name         = "ApproximateNumberOfMessagesVisible"
  namespace           = "AWS/SQS"
  period              = 60
  statistic           = "Maximum"
  threshold           = 1

  dimensions = {
    QueueName = var.face_recognition_queue_name
  }

  alarm_actions = [var.ai_task_scale_out_policy_arn]
}

# Scale-IN 알람: 대기 메시지가 0개
resource "aws_cloudwatch_metric_alarm" "sqs_scale_in" {
  alarm_name          = "yuno-ai-task-scale-in-${var.env}"
  comparison_operator = "LessThanThreshold"
  evaluation_periods  = 3             # 3분 연속 0개일 때 축소
  metric_name         = "ApproximateNumberOfMessagesVisible"
  namespace           = "AWS/SQS"
  period              = 60
  statistic           = "Maximum"
  threshold           = 1

  dimensions = {
    QueueName = var.face_recognition_queue_name
  }

  alarm_actions = [var.ai_task_scale_in_policy_arn]
}

resource "aws_cloudwatch_metric_alarm" "video_task_scale_out" {
  alarm_name          = "yuno-video-task-scale-out-${var.env}"
  comparison_operator = "GreaterThanOrEqualToThreshold"
  evaluation_periods  = 1
  metric_name         = "ApproximateNumberOfMessagesVisible"
  namespace           = "AWS/SQS"
  period              = 60
  statistic           = "Maximum"
  threshold           = 1

  dimensions = {
    QueueName = var.video_queue_name
  }

  alarm_actions = [var.video_task_scale_out_policy_arn]
}

resource "aws_cloudwatch_metric_alarm" "video_task_scale_in" {
  alarm_name          = "yuno-video-task-scale-in-${var.env}"
  comparison_operator = "LessThanThreshold"
  evaluation_periods  = 3
  metric_name         = "ApproximateNumberOfMessagesVisible"
  namespace           = "AWS/SQS"
  period              = 60
  statistic           = "Maximum"
  threshold           = 1

  dimensions = {
    QueueName = var.video_queue_name
  }

  alarm_actions = [var.video_task_scale_in_policy_arn]
}