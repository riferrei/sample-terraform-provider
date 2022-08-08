variable "token" {
  type = string
  default = "<TOKEN_FROM_CRUDCRUD_WEBSITE>"
}

provider "sample" {
  token = var.token
}

resource "sample_character" "daredevil" {
  fullname = "DareDevil"
  identity = "Matt Murdock"
  knownas = "The man without fear"
  type = "super-hero"
}

output "reveal_secret_identity" {
  value = "The real identity of ${sample_character.daredevil.fullname} is '${sample_character.daredevil.identity}'"
}
