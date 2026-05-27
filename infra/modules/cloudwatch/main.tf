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