data "aws_iam_role" "ecs_autoscaling_service_role" {
  name     = "AWSServiceRoleForApplicationAutoScaling_ECSService"
  provider = aws.global
}

data "aws_sns_topic" "custom_cloudwatch_alarms" {
  name     = "ecs_autoscaling_alarms"
  provider = aws.region
}

module "app_ecs_autoscaling" {
  source                           = "./modules/ecs_autoscaling"
  environment                      = local.name_prefix
  aws_ecs_cluster_name             = aws_ecs_cluster.main.name
  aws_ecs_service_name             = module.app.ecs_service.name
  ecs_autoscaling_service_role_arn = data.aws_iam_role.ecs_autoscaling_service_role.arn
  ecs_task_autoscaling_minimum     = var.ecs_task_autoscaling.minimum
  ecs_task_autoscaling_maximum     = var.ecs_task_autoscaling.maximum
  max_scaling_alarm_actions        = [data.aws_sns_topic.custom_cloudwatch_alarms.arn]
  providers = {
    aws.region = aws.region
  }
}
