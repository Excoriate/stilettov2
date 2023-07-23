resource "random_string" "password" {
  length  = 16
  special = true
}

output "password" {
  value = random_string.password.result
}
