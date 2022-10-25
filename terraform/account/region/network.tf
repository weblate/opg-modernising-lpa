module "network" {
  source                         = "github.com/ministryofjustice/opg-terraform-aws-network?ref=v1.2.0-MLPAB-344.1"
  cidr                           = var.network_cidr_block
  default_security_group_ingress = []
  default_security_group_egress  = []
  providers = {
    aws = aws.region
  }
}
