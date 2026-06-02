locals {
  env        = "production"
  aws_region = "ap-northeast-1"

  common_tags = {
    Project     = "yuno"
    Environment = local.env
    ManagedBy   = "terraform"
  }

  app_ssm_params = toset([
    "APP_ENV",
    "APP_URL",
    "API_URL",
    "MEDIA_URL",
    "MEDIA_BUCKET_NAME",
    "LINE_CHANNEL_ID",
    "KAKAO_CLIENT_ID",
  ])

  app_ssm_secret_params = toset([
    "DATABASE_URL",
    "AUTH_SECRET_KEY",
    "LINE_CHANNEL_SECRET",
    "KAKAO_CLIENT_SECRET",
  ])

  # Lambda / ECS 공통 환경변수
  # ECS 모듈 outputs(cluster_arn 등)은 순환 참조 방지를 위해 제외 → lambda 모듈 호출부에서 merge
  shared_app_env = {
    APP_ENV             = local.env
    APP_URL             = "https://${var.domain_name}"
    API_URL             = "https://api.${var.domain_name}"
    MEDIA_URL           = "https://media.${var.domain_name}"
    MEDIA_BUCKET_NAME   = module.s3.media_bucket_id
    SQS_QUEUE_URL       = module.sqs.queue_url
    VIDEO_SQS_QUEUE_URL = module.sqs.video_queue_url
    STORAGE_TYPE        = "s3"
    DATABASE_URL        = var.database_url
    AUTH_SECRET_KEY     = var.auth_secret_key
    LINE_CHANNEL_ID     = var.line_channel_id
    LINE_CHANNEL_SECRET = var.line_channel_secret
    KAKAO_CLIENT_ID     = var.kakao_client_id
    KAKAO_CLIENT_SECRET = var.kakao_client_secret
    CLOUDFRONT_KEY_PAIR_ID = module.cloudfront.media_signing_key_pair_id
    CLOUDFRONT_PRIVATE_KEY = var.cloudfront_private_key
    CLOUDFRONT_MEDIA_DOMAIN = "media.${var.domain_name}"

  }
}

