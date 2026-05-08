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