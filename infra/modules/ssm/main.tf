
resource "aws_ssm_parameter" "app_ssm_params" {
    for_each = var.app_ssm_params
    name = "/yuno/${var.env}/${each.value}"
    type = "String"
    value = each.value

    lifecycle {
        ignore_changes = [value]
    }
}

resource "aws_ssm_parameter" "app_ssm_secret_params" {
    for_each = var.app_ssm_secret_params
    name = "/yuno/${var.env}/${each.value}"
    type = "SecureString"
    value = "CHANGE_ME"

    lifecycle {
        ignore_changes = [value]
    }
}
