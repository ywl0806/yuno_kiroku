# ── ECR (공용 상태 파일에서 output 참조) ─────────────────────

data "terraform_remote_state" "route53" {
  backend = "s3"
  config = {
    bucket = "terraform-state-bucket-yuno"
    key    = "shared/route53/terraform.tfstate"
    region = "ap-northeast-1"
  }
}

data "terraform_remote_state" "acm" {
  backend = "s3"
  config = {
    bucket = "terraform-state-bucket-yuno"
    key    = "shared/acm/terraform.tfstate"
    region = "ap-northeast-1"
  }
}

data "terraform_remote_state" "ecr" {
  backend = "s3"
  config = {
    bucket = "terraform-state-bucket-yuno"
    key    = "shared/ecr/terraform.tfstate"
    region = "ap-northeast-1"
  }
}

# ── AWS 계정 정보 ─────────────────────────────────────────────

data "aws_caller_identity" "current" {}

# ── 기본 VPC / Subnet (Phase 1: VPC 없이 default VPC 사용) ───

data "aws_vpc" "default" {
  default = true
}

data "aws_subnets" "public" {
  filter {
    name   = "vpc-id"
    values = [data.aws_vpc.default.id]
  }
}
# ── SSM ────────────────────────────────────────────────────────
module "ssm" {
  source = "../../modules/ssm"
  env = local.env
  app_ssm_params = local.app_ssm_params
  app_ssm_secret_params = local.app_ssm_secret_params
}

# ── IAM ──────────────────────────────────────────────────────

module "iam" {
  source = "../../modules/iam"
  env = local.env
  media_bucket_arn = module.s3.media_bucket_arn
  lambda_zip_bucket_arn = module.s3.lambda_zip_bucket_arn
  common_tags = local.common_tags
  ecr_repository_arns = values(data.terraform_remote_state.ecr.outputs.repository_arns)
  aws_region = "ap-northeast-1"
  aws_account_id = data.aws_caller_identity.current.account_id
}

# ── S3 ────────────────────────────────────────────────────────

module "s3" {
  source = "../../modules/s3"
  env    = local.env
  media_cors_origins = ["https://${var.domain_name}"]
  frontend_distribution_arn = module.cloudfront.frontend_distribution_arn
  media_distribution_arn = module.cloudfront.media_distribution_arn
  tags = local.common_tags
}

# ── CloudFront ────────────────────────────────────────────────

module "cloudfront" {
  source = "../../modules/cloudfront"
  env = local.env
  media_bucket_id = module.s3.media_bucket_id
  media_bucket_domain = module.s3.media_bucket_domain
  frontend_bucket_id = module.s3.frontend_bucket_id
  frontend_bucket_domain = module.s3.frontend_bucket_domain
  frontend_domain =  var.domain_name 
  media_domain    =  "media.${var.domain_name}" 
  acm_certificate_arn =  data.terraform_remote_state.acm.outputs.cloudfront_certificate_arn
  common_tags = local.common_tags
}

# ── Lambda ────────────────────────────────────────────────────

module "lambda" {
  source = "../../modules/lambda"
  env = local.env
  api_lambda_role_arn = module.iam.api_lambda_role_arn
  resize_lambda_role_arn = module.iam.resize_lambda_role_arn
  resize_ecr_image_uri = "${data.terraform_remote_state.ecr.outputs.repository_urls["yuno-resize"]}:${local.env}"
  media_bucket_id = module.s3.media_bucket_id
  media_bucket_arn = module.s3.media_bucket_arn
  common_tags = local.common_tags
  api_lambda_zip_bucket_arn = module.s3.lambda_zip_bucket_arn
  api_lambda_zip_bucket_name = module.s3.lambda_zip_bucket_name
  api_lambda_zip_bucket_key = aws_s3_object.api_lambda_zip.key
  app_env_vars = {
    APP_ENV             = local.env
    APP_URL             = "https://${var.domain_name}"
    API_URL             = "https://api.${var.domain_name}"
    MEDIA_URL           = "https://media.${var.domain_name}"
    MEDIA_BUCKET_NAME   = module.s3.media_bucket_id
    STORAGE_TYPE        = "s3"
    DATABASE_URL        = var.database_url
    AUTH_SECRET_KEY     = var.auth_secret_key
    LINE_CHANNEL_ID     = var.line_channel_id
    LINE_CHANNEL_SECRET = var.line_channel_secret
    KAKAO_CLIENT_ID     = var.kakao_client_id
    KAKAO_CLIENT_SECRET = var.kakao_client_secret
  }
}



# ── S3 Object Lambda ────────────────────────────────────────────

resource "aws_s3_object" "api_lambda_zip" {
  bucket = module.s3.lambda_zip_bucket_name
  key = "function.zip"
  source = "${path.module}/files/function.zip"
}

# ── ECS ────────────────────────────────────────────────────────

module "ecs" {
  source = "../../modules/ecs"
  env = local.env
  common_tags = local.common_tags
  aws_region = "ap-northeast-1"
  ecs_task_role_arn = module.iam.ecs_task_role_arn
  ecs_task_execution_role_arn = module.iam.ecs_task_execution_role_arn
  vpc_id = data.aws_vpc.default.id
  ai_image_uri = "${data.terraform_remote_state.ecr.outputs.repository_urls["yuno-ai"]}:${local.env}"
}

# ── API Gateway ────────────────────────────────────────────────

module "api_gateway" {
  source = "../../modules/api_gateway"
  env = local.env
  aws_region = "ap-northeast-1"
  aws_account_id = data.aws_caller_identity.current.account_id
  api_lambda_function_name = module.lambda.api_lambda_function_name
  api_lambda_invoke_arn = module.lambda.api_lambda_function_invoke_arn
  api_domain          = "api.${var.domain_name}"
  acm_certificate_arn = data.terraform_remote_state.acm.outputs.api_gw_certificate_arn
}

# ── Route53 DNS 레코드 ─────────────────────────────────────────

module "route53" {
  source = "../../modules/route53"

  hosted_zone_id        = data.terraform_remote_state.route53.outputs.hosted_zone_id
  domain_name           = var.domain_name
  frontend_cf_domain    = module.cloudfront.frontend_distribution_domain
  media_cf_domain       = module.cloudfront.media_distribution_domain
  api_gw_domain_target  = module.api_gateway.custom_domain_target
  api_gw_domain_zone_id = module.api_gateway.custom_domain_zone_id
}