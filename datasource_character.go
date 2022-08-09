package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
	endpoint := meta.(string)
	identity := data.Get(identityField).(string)
	response, err := http.Get(endpoint)
	if err != nil {
		log.Fatal(err)
	}

	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	var characterList []CharacterData
	json.Unmarshal(bodyBytes, &characterList)
	if err == nil && len(characterList) > 0 {
		for _, character := range characterList {
			if strings.EqualFold(identity, character.Identity) {
				data.SetId(character.ID)
				data.Set(fullNameField, character.FullName)
				data.Set(identityField, character.Identity)
				data.Set(knownasField, character.KnownAs)
				data.Set(typeField, character.Type)
				break
			}
		}
	}

	return diags

}
