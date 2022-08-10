package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceMarvelCharacter() *schema.Resource {
	return &schema.Resource{
		CreateContext: characterCreate,
		ReadContext:   characterRead,
		UpdateContext: characterUpdate,
		DeleteContext: characterDelete,
		Schema: map[string]*schema.Schema{
			fullNameField: {
				Type:     schema.TypeString,
				Required: true,
			},
			identityField: {
				Type:     schema.TypeString,
				Required: true,
			},
			knownasField: {
				Type:     schema.TypeString,
				Required: true,
			},
			typeField: {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: func(v interface{}, k string) (wrs []string, ers []error) {
					var errors []error
					var warns []string
					value, _ := v.(string)
					validType := false
					for _, typeValue := range types {
						if value == typeValue {
							validType = true
							break
						}
					}
					if !validType {
						errors = append(errors, fmt.Errorf("Invalid value for "+
							"type. Valid values are: "+
							strings.Join(types, ", ")))
						return warns, errors
					}
					return warns, errors
				},
			},
		},
	}
}

func characterCreate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {

	var diags diag.Diagnostics
	endpoint := meta.(string)
	character := &MarvelCharacter{
		FullName: data.Get(fullNameField).(string),
		Identity: data.Get(identityField).(string),
		KnownAs:  data.Get(knownasField).(string),
		Type:     data.Get(typeField).(string),
	}

	jsonBody, _ := json.Marshal(character)
	bodyReader := bytes.NewReader(jsonBody)
	request, err := http.NewRequest(http.MethodPost, endpoint, bodyReader)
	if err != nil {
		return diag.FromErr(err)
	}

	request.Header.Set("Content-Type", "application/json; charset=utf-8")
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return diag.FromErr(err)
	}
	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return diag.FromErr(err)
	}

	createdCharacter := &MarvelCharacter{}
	json.Unmarshal(bodyBytes, createdCharacter)
	data.SetId(createdCharacter.ID)
	data.Set(fullNameField, createdCharacter.FullName)
	data.Set(identityField, createdCharacter.Identity)
	data.Set(knownasField, createdCharacter.KnownAs)
	data.Set(typeField, createdCharacter.Type)

	return diags

}

func characterRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {

	var diags diag.Diagnostics
	endpoint := meta.(string) + "/" + data.Id()
	response, err := http.Get(endpoint)
	if err != nil {
		log.Fatal(err)
	}

	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	character := &MarvelCharacter{}
	json.Unmarshal(bodyBytes, character)
	data.Set(fullNameField, character.FullName)
	data.Set(identityField, character.Identity)
	data.Set(knownasField, character.KnownAs)
	data.Set(typeField, character.Type)

	return diags

}

func characterUpdate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {

	endpoint := meta.(string) + "/" + data.Id()
	character := &MarvelCharacter{
		FullName: data.Get(fullNameField).(string),
		Identity: data.Get(identityField).(string),
		KnownAs:  data.Get(knownasField).(string),
		Type:     data.Get(typeField).(string),
	}
	jsonBody, _ := json.Marshal(character)
	bodyReader := bytes.NewReader(jsonBody)
	request, err := http.NewRequest(http.MethodPut, endpoint, bodyReader)
	if err != nil {
		return diag.FromErr(err)
	}

	request.Header.Set("Content-Type", "application/json; charset=utf-8")
	_, err = http.DefaultClient.Do(request)
	if err != nil {
		return diag.FromErr(err)
	}

	return characterRead(ctx, data, meta)

}

func characterDelete(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {

	var diags diag.Diagnostics
	endpoint := meta.(string) + "/" + data.Id()
	character := &MarvelCharacter{
		FullName: data.Get(fullNameField).(string),
		Identity: data.Get(identityField).(string),
		KnownAs:  data.Get(knownasField).(string),
		Type:     data.Get(typeField).(string),
	}

	jsonBody, _ := json.Marshal(character)
	bodyReader := bytes.NewReader(jsonBody)
	request, err := http.NewRequest(http.MethodDelete, endpoint, bodyReader)
	if err != nil {
		return diag.FromErr(err)
	}

	request.Header.Set("Content-Type", "application/json; charset=utf-8")
	_, err = http.DefaultClient.Do(request)
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId("")
	return diags

}
