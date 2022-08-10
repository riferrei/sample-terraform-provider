package main

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	backendAPI = "https://crudcrud.com/api/${token}/sample"
)

func Provider() *schema.Provider {

	provider := &schema.Provider{
		Schema: map[string]*schema.Schema{
			tokenField: {
				Type:     schema.TypeString,
				Required: true,
			},
			timeoutField: {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  300,
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"sample_marvel_character": resourceMarvelCharacter(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"sample_marvel_character": datasourceMarvelCharacter(),
		},
		ConfigureContextFunc: providerConfigure,
	}

	return provider

}

func providerConfigure(ctx context.Context, data *schema.ResourceData) (interface{}, diag.Diagnostics) {

	token := data.Get(tokenField).(string)
	timeout := data.Get(timeoutField).(int)
	endpoint := strings.ReplaceAll(backendAPI, "${token}", token)

	session := &Session{
		Endpoint: endpoint,
		HttpClient: &http.Client{
			Timeout: time.Duration(timeout) * time.Second,
		},
	}

	return session, nil

}
