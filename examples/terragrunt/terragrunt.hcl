locals {
  terraform_version  = get_env("TERRAGRUNT_CFG_BINARIES_TERRAFORM_VERSION", "1.5.1")
  terragrunt_version = get_env("TERRAGRUNT_CFG_BINARIES_TERRAGRUNT_VERSION", "0.42.8")
}

terraform {
  extra_arguments "optional_vars" {
    commands = [
      "apply",
      "destroy",
      "plan",
    ]
  }

}
  generate "terraform_version" {
    path              = ".terraform-version"
    if_exists         = "overwrite"
    disable_signature = true

    contents = <<-EOF
    ${local.terraform_version}
  EOF
  }

  generate "terragrunt_version" {
    path              = ".terragrunt-version"
    if_exists         = "overwrite"
    disable_signature = true

    contents = <<-EOF
    ${local.terragrunt_version}
  EOF
  }

  generate "providers" {
    path      = "providers.tf"
    if_exists = "overwrite_terragrunt"

    contents = templatefile("${get_repo_root()}/examples/terragrunt/providers.tf.tmpl", {
      aws_region_passed_from_env = "us-east-1"
    })
  }
