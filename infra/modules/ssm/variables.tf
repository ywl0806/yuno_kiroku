variable "env" {
    type = string
    description = "배포 환경"
}

variable "app_ssm_params" {
    type = set(string)
    description = "애플리케이션 환경변수"
    default = []
}

variable "app_ssm_secret_params" {
    type = set(string)
    description = "애플리케이션 비밀 환경변수"
    default = []
}