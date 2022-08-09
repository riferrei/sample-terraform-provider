variable "crudCrudToken" {
  type = string
  default = "<TOKEN_FROM_CRUDCRUD_WEBSITE>"
}

provider "sample" {
  token = var.crudCrudToken
}

resource "sample_marvel_character" "daredevil" {
  fullname = "DareDevil"
  identity = "Matt Murdock"
  knownas = "The man without fear"
  type = "super-hero"
}

output "daredevil_secret_identity" {
  value = "The secret identity of ${sample_marvel_character.daredevil.fullname} is '${sample_marvel_character.daredevil.identity}'"
}
