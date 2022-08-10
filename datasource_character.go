package main

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"golang.org/x/exp/slices"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func datasourceMarvelCharacter() *schema.Resource {
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
	session := meta.(*Session)
	identity := data.Get(identityField).(string)

	request, err := http.NewRequest(http.MethodGet, session.Endpoint, nil)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failure to create HTTP request",
			Detail:   err.Error(),
		})
		return diags
	}
	response, err := session.HttpClient.Do(request)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failure to execute HTTP request",
			Detail:   err.Error(),
		})
		return diags
	}
	defer response.Body.Close()
	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failure to read response",
			Detail:   err.Error(),
		})
		return diags
	}

	var characterList []MarvelCharacter
	json.Unmarshal(bodyBytes, &characterList)

	if err == nil && len(characterList) > 0 {
		idx := slices.IndexFunc(characterList,
			func(character MarvelCharacter) bool {
				return strings.EqualFold(identity, character.Identity)
			},
		)
		if idx >= 0 {
			character := characterList[idx]
			data.SetId(character.ID)
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
	}

	return diags

}
