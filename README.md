Sample Terraform Provider
=========================

This repository contains a sample code implementation for a Terraform provider. It is meant to be used as a way to teach, educate, and show the internals of a provider. Even if you are not looking to learn how to build custom providers, you may benefit from learning how one works behind the scenes. For more information about how to build custom providers, please visit the [HashiCorp Learn platform section about this](https://learn.hashicorp.com/tutorials/terraform/provider-use?in=terraform/providers).

**Important**: to play with this provider, you are going to create a new CRUD REST API which is used as the upstream backend service this provider communicates with. For simplicity, this projects leverages the REST APIs from the service [CrudCrud](https://crudcrud.com). It gives you a free API available for 24 hours, where you can send up to 100 requests.

Requirements
------------

- [Terraform](https://www.terraform.io/downloads.html) 0.12+
- [Go](https://golang.org/doc/install) 1.19 (to build the provider plugin)

Building the Provider
----------------------

To play with this provider, first you have to build it. Then you must install the native executable generated into your local plugins folder as explained [here](https://www.terraform.io/docs/plugins/basics.html#installing-a-plugin). The easiest way to do wall of this is running the following command:

```console
make install
```

Once you have done this, you can start writing your `.tf` files and run the command `terraform init` to initialize the plugins of your project.

Using the Provider
----------------------

The provider allows you to create, read, update, and delete characters from Marvel. Here is an example that you can use to play with the provider:

```
provider "sample" {
  token = <TOKEN_FROM_CRUDCRUD_WEBSITE>
}

resource "sample_marvel_character" "deadpool" {
  fullname = "Deadpool"
  identity = "Wade Wilson"
  knownas = "Merc with a Mouth"
  type = "anti-hero"
}
```

# License

This project is licensed under the [Apache 2.0 License](./LICENSE).
