package main

import (
	"context"
	"crypto/tls"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/opensearch-project/opensearch-go/v2"
	"github.com/opensearch-project/opensearch-go/v2/opensearchapi"
)

func Provider() *schema.Provider {

	provider := &schema.Provider{
		Schema: map[string]*schema.Schema{
			backendAddressField: {
				Type:     schema.TypeString,
				Optional: true,
				Default:  backendAddress,
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"buildonaws_character": resourceCharacter(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"buildonaws_character": datasourceCharacter(),
		},
		ConfigureContextFunc: providerConfigure,
	}

	return provider

}

func providerConfigure(ctx context.Context, data *schema.ResourceData) (interface{}, diag.Diagnostics) {

	backendAddress := data.Get(backendAddressField).(string)
	backendClient, _ := opensearch.NewClient(
		opensearch.Config{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
			Addresses: []string{backendAddress},
		},
	)

	pingRequest := opensearchapi.PingRequest{
		Pretty:     true,
		Human:      true,
		ErrorTrace: true,
	}
	_, err := pingRequest.Do(ctx, backendClient)
	if err != nil {
		var diags diag.Diagnostics
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failure connecting with the backend",
			Detail:   "Reason: " + err.Error(),
		})
		return nil, diags
	}

	return backendClient, nil

}
