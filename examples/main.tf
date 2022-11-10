provider "buildonaws" {
}

resource "buildonaws_character" "daredevil" {
  fullname = "Daredevil"
  identity = "Matt Murdock"
  knownas = "The man without fear"
  type = "super-hero"
}

output "daredevil_secret_identity" {
  value = "The secret identity of ${buildonaws_character.daredevil.fullname} is '${buildonaws_character.daredevil.identity}'"
}
