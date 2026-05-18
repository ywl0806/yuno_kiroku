output "resize_queue_arn" {
    value = aws_sqs_queue.resize.arn
}

output "resize_queue_url" {
    value = aws_sqs_queue.resize.url
}

output "resize_dlq_arn" {
    value = aws_sqs_queue.resize_dlq.arn
}

output "queue_arn" {
    value = aws_sqs_queue.face_recognition.arn
}

output "queue_url" {
    value = aws_sqs_queue.face_recognition.url
}

output "dlq_arn" {
    value = aws_sqs_queue.face_recognition_dlq.arn
}

output "dlq_url" {
    value = aws_sqs_queue.face_recognition_dlq.url
}

output "queue_name" {
    value = aws_sqs_queue.face_recognition.name
}

output "video_queue_arn" {
    value = aws_sqs_queue.video_processing.arn
}

output "video_queue_url" {
    value = aws_sqs_queue.video_processing.url
}

output "video_queue_name" {
    value = aws_sqs_queue.video_processing.name
}