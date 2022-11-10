Sample Terraform Provider
=========================

This repository contains a sample code implementation for a Terraform provider. It is meant to be used as a way to teach, educate, and show the internals of a provider. Even if you are not looking to learn how to build custom providers, you may benefit from learning how one works behind the scenes. For more information about how to build custom providers, please visit the [HashiCorp Learn platform section about this](https://learn.hashicorp.com/tutorials/terraform/provider-use?in=terraform/providers).

**Important**: this provider uses [OpenSearch](https://opensearch.org) as the upstream backend service. Before playing with the provider, you first need to get OpenSearch up-and-running. For simplicity, you can use the [docker-compose.yml](./docker-compose.yml) available in this repository.

Requirements
------------

- [Terraform](https://www.terraform.io/downloads.html) 0.12+
- [Go](https://golang.org/doc/install) 1.19 (to build the provider plugin)
- [Docker](https://www.docker.com/get-started) 4.11+ (to execute the backend)

Building the Provider
---------------------

To play with this provider, first you have to build it. Then you must install the native executable generated into your local plugins folder as explained [here](https://www.terraform.io/docs/plugins/basics.html#installing-a-plugin). The easiest way to do wall of this is running the following command:

```bash
make install
```

Once you have done this, you can start writing your `.tf` files and run the command `terraform init` to initialize the plugins of your project.

Using the Provider
------------------

The provider allows you to create, read, update, and delete comic characters. Here is an example that you can use to play with the provider:

```tcl
provider "buildonaws" {
}

resource "buildonaws_character" "deadpool" {
  fullname = "Deadpool"
  identity = "Wade Wilson"
  knownas = "Merc with a Mouth"
  type = "anti-hero"
}
```

Debugging the Provider
----------------------

This provider was created using the [Terraform Plugin SDKv2](https://www.terraform.io/plugin/sdkv2), which allows developers to debug the code using tools like [delve](https://github.com/go-delve/delve). The code and build have been written to allow debugging, so all you have to do is starting a debugging session. You can do this using Visual Studio Code or via command-line.

### Visual Studio Code

1. Create a file named `launch.json` in the `.vscode` folder with the following content:

```json
{
    "version": "0.2.0",
    "configurations": [
        {
            "name": "Debug Provider",
            "type": "go",
            "request": "launch",
            "mode": "debug",
            "program": "${workspaceFolder}",
            "env": {},
            "args": [
                "-debug",
            ]
        }
    ]
}
```

2. Run the editor debugger. This will output the reattach configuration to the debug console.

### Command-line

1. Start a delve debugging session:

```bash
dlv exec --accept-multiclient --continue --headless ./terraform-provider-sample -- -debug
```

This will output the reattach configuration to the output.

### Debugging

Before issuing commands like `plan`, `apply`, and `destroy` to your Terraform code, you need to export the `TF_REATTACH_PROVIDERS` environment variable to reattach the CLI session to the started debugging server. After starting the debugging server session using one method above, the value to this variable will be provided. Here is an example of how to do this.

```bash
export TF_REATTACH_PROVIDERS='{"riferrei.com/terraform/sample":{"Protocol":"grpc","Pid":3382870,"Test":true,"Addr":{"Network":"unix","String":"/tmp/plugin713096927"}}}'
```

# License

This project is licensed under the [Apache 2.0 License](./LICENSE).
