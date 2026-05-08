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