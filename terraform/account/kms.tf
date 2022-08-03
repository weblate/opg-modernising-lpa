resource "aws_kms_key" "cloudwatch" {
  description             = "${local.mandatory_moj_tags.application} cloudwatch application logs encryption key for ${data.aws_region.eu_west_1.name}"
  deletion_window_in_days = 10
  enable_key_rotation     = true
  policy                  = local.account.account_name == "development" ? data.aws_iam_policy_document.cloudwatch_kms_merged.json : data.aws_iam_policy_document.cloudwatch_kms.json
  multi_region            = true
  provider                = aws.eu_west_1
}

resource "aws_kms_alias" "cloudwatch_alias" {
  name          = "alias/${local.mandatory_moj_tags.application}_cloudwatch_application_logs_encryption_${data.aws_region.eu_west_1.name}"
  target_key_id = aws_kms_key.cloudwatch.key_id
  provider      = aws.eu_west_1
}

# See the following link for further information
# https://docs.aws.amazon.com/kms/latest/developerguide/key-policies.html
data "aws_iam_policy_document" "cloudwatch_kms_merged" {
  provider = aws.global
  source_policy_documents = [
    data.aws_iam_policy_document.cloudwatch_kms.json,
    data.aws_iam_policy_document.cloudwatch_kms_development_account_operator_admin.json
  ]
}

data "aws_iam_policy_document" "cloudwatch_kms" {
  provider = aws.global
  statement {
    sid       = "Allow Key to be used for Encryption"
    effect    = "Allow"
    resources = ["*"]
    actions = [
      "kms:Encrypt",
      "kms:Decrypt",
      "kms:ReEncrypt*",
      "kms:GenerateDataKey*",
      "kms:DescribeKey",
    ]

    principals {
      type = "Service"
      identifiers = [
        "logs.${data.aws_region.eu_west_1.name}.amazonaws.com",
        "logs.${data.aws_region.eu_west_2.name}.amazonaws.com",
        "events.amazonaws.com"
      ]
    }
  }

  statement {
    sid       = "Key Administrator"
    effect    = "Allow"
    resources = ["*"]
    actions = [
      "kms:Create*",
      "kms:Describe*",
      "kms:Enable*",
      "kms:List*",
      "kms:Put*",
      "kms:Update*",
      "kms:Revoke*",
      "kms:Disable*",
      "kms:Get*",
      "kms:Delete*",
      "kms:TagResource",
      "kms:UntagResource",
      "kms:ScheduleKeyDeletion",
      "kms:CancelKeyDeletion"
    ]

    principals {
      type = "AWS"
      identifiers = [
        "arn:aws:iam::${data.aws_caller_identity.global.account_id}:role/breakglass",
        "arn:aws:iam::${data.aws_caller_identity.global.account_id}:role/${var.default_role}",
      ]
    }
  }
}

data "aws_iam_policy_document" "cloudwatch_kms_development_account_operator_admin" {
  provider = aws.global
  statement {
    sid       = "Dev Account Key Administrator"
    effect    = "Allow"
    resources = ["*"]
    actions = [
      "kms:Create*",
      "kms:Describe*",
      "kms:Enable*",
      "kms:List*",
      "kms:Put*",
      "kms:Update*",
      "kms:Revoke*",
      "kms:Disable*",
      "kms:Get*",
      "kms:Delete*",
      "kms:TagResource",
      "kms:UntagResource",
      "kms:ScheduleKeyDeletion",
      "kms:CancelKeyDeletion"
    ]

    principals {
      type = "AWS"
      identifiers = [
        "arn:aws:iam::${data.aws_caller_identity.global.account_id}:role/operator"
      ]
    }
  }
}
