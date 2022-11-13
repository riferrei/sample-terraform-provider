provider "buildonaws" {
}

resource "buildonaws_character" "daredevil" {
  fullname = "Daredevil"
  identity = "Matt Murdock"
  knownas = "The man without fear"
  type = "super-hero"
}

// Execute the script `deadpool.sh` to create
// the Deadpool character in the backend.
data "buildonaws_character" "deadpool" {
  identity = "Wade"
}

output "daredevil_secret_identity" {
  value = "The secret identity of ${buildonaws_character.daredevil.fullname} is '${buildonaws_character.daredevil.identity}'"
}

output "deadpool_is_knownas" {
  value = "${data.buildonaws_character.deadpool.fullname} is also known as '${data.buildonaws_character.deadpool.knownas}'"
}
