resource "aws_sqs_queue" "video_processing_dlq" {
    name = "yuno-video-processing-dlq-${var.env}"
    tags = var.common_tags

    message_retention_seconds = 1209600 # 14일
}

resource "aws_sqs_queue" "video_processing" {
    name = "yuno-video-processing-${var.env}"
    tags = var.common_tags

    message_retention_seconds  = 86400 # 1일
    visibility_timeout_seconds = 900   # 15분

    redrive_policy = jsonencode({
        deadLetterTargetArn = aws_sqs_queue.video_processing_dlq.arn
        maxReceiveCount     = 2
    })
}

resource "aws_sqs_queue" "resize_dlq" {
    name = "yuno-resize-dlq-${var.env}"
    tags = var.common_tags

    message_retention_seconds = 1209600 # 14일
}

resource "aws_sqs_queue" "resize" {
    name = "yuno-resize-${var.env}"
    tags = var.common_tags

    message_retention_seconds  = 86400 # 1일
    visibility_timeout_seconds = 660   # Lambda timeout(600s) + 10%

    redrive_policy = jsonencode({
        deadLetterTargetArn = aws_sqs_queue.resize_dlq.arn
        maxReceiveCount     = 3
    })
}

# S3가 SQS에 메시지를 보낼 수 있도록 리소스 기반 큐 정책
resource "aws_sqs_queue_policy" "resize_s3_send" {
    queue_url = aws_sqs_queue.resize.id

    policy = jsonencode({
        Version = "2012-10-17"
        Statement = [{
            Effect    = "Allow"
            Principal = { Service = "s3.amazonaws.com" }
            Action    = "sqs:SendMessage"
            Resource  = aws_sqs_queue.resize.arn
            Condition = {
                ArnLike = { "aws:SourceArn" = var.media_bucket_arn }
            }
        }]
    })
}

resource "aws_sqs_queue" "face_recognition_dlq" {
    name = "yuno-face-recognition-dlq-${var.env}"
    tags = var.common_tags

    message_retention_seconds = 1209600 # 14일
}

resource "aws_sqs_queue" "face_recognition" {
    name = "yuno-face-recognition-${var.env}"
    tags = var.common_tags

    message_retention_seconds = 86400 # 1일

    visibility_timeout_seconds = 300 # 5분

    redrive_policy = jsonencode({
        deadLetterTargetArn = aws_sqs_queue.face_recognition_dlq.arn
        maxReceiveCount     = 3
    })
}