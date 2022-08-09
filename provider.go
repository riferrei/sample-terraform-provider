package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	backendAPI = "https://crudcrud.com/api/${token}/sample"
)

func Provider() *schema.Provider {

	provider := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"token": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"sample_character": resourceCharacter(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"sample_character": datasourceCharacter(),
		},
		ConfigureContextFunc: providerConfigure,
	}

	return provider

}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {

	token := d.Get("token").(string)
	endpoint := strings.ReplaceAll(backendAPI, "${token}", token)
	response, err := http.Get(endpoint)
	if err != nil {
		log.Fatal(err)
	}

	var diags diag.Diagnostics
	if response.StatusCode == 400 {
		detail, err := io.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to invoke the backend API",
			Detail:   "Reason: " + string(detail),
		})
		return nil, diags
	}

	return endpoint, nil

}
