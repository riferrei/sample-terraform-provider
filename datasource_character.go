package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/opensearch-project/opensearch-go/v2"
	"github.com/opensearch-project/opensearch-go/v2/opensearchapi"
)

func datasourceCharacter() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourceCharacterRead,
		Schema: map[string]*schema.Schema{
			fullNameField: {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			identityField: {
				Type:     schema.TypeString,
				Required: true,
			},
			knownasField: {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			typeField: {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func datasourceCharacterRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {

	var diags diag.Diagnostics
	backendClient := meta.(*opensearch.Client)
	identity := data.Get(identityField).(string)

	searchBody := &struct {
		Query struct {
			Match struct {
				Identity string `json:"identity,omitempty"`
			} `json:"match,omitempty"`
		} `json:"query,omitempty"`
	}{}
	searchBody.Query.Match.Identity = identity
	bodyContent, _ := json.Marshal(searchBody)
	bodyReader := bytes.NewReader(bodyContent)

	searchRequest := opensearchapi.SearchRequest{
		Index: []string{backendIndex},
		Body:  bodyReader,
	}
	searchResponse, err := searchRequest.Do(ctx, backendClient)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failure to retrieve character",
			Detail:   "Reason: " + err.Error(),
		})
		return diags
	}
	defer searchResponse.Body.Close()
	bodyContent, err = io.ReadAll(searchResponse.Body)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failure reading response",
			Detail:   "Reason: " + err.Error(),
		})
		return diags
	}

	backendSearchResponse := &BackendSearchResponse{}
	json.Unmarshal(bodyContent, backendSearchResponse)

	if backendSearchResponse.Hits.Total.Value > 0 {
		data.SetId(backendSearchResponse.Hits.Hits[0].ID)
		character := backendSearchResponse.Hits.Hits[0].Source
		data.Set(fullNameField, character.FullName)
		data.Set(identityField, character.Identity)
		data.Set(knownasField, character.KnownAs)
		data.Set(typeField, character.Type)
	} else {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "Datasource was not loaded",
			Detail:   "Reason: no character with the identity '" + identity + "'.",
		})
	}

	return diags

}
