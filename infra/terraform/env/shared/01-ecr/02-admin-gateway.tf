module "ecr_admin_gateawy" {
  source = "./../../../modules/ecr"

  #####################################################################
  # ECR Repository (Public)
  #####################################################################
  repository_name = "furumaru/admin-gateway"
  description     = ""

  operating_systems = ["Linux"]
  architectures     = ["ARM 64", "x64-64"]

  about_text = ""
  usage_text = ""
}
