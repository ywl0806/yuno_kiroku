locals {
  env = "production"
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

}

