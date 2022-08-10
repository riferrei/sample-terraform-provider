package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
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
					for _, characterType := range characterTypes {
						if value == characterType {
							validType = true
							break
						}
					}
					if !validType {
						errors = append(errors, fmt.Errorf("Invalid value for "+
							"type. Valid values are: "+
							strings.Join(characterTypes, ", ")))
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
	session := meta.(*Session)
	character := &MarvelCharacter{
		FullName: data.Get(fullNameField).(string),
		Identity: data.Get(identityField).(string),
		KnownAs:  data.Get(knownasField).(string),
		Type:     data.Get(typeField).(string),
	}

	jsonBody, _ := json.Marshal(character)
	bodyReader := bytes.NewReader(jsonBody)
	request, err := http.NewRequest(http.MethodPost, session.Endpoint, bodyReader)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failure to create HTTP request",
			Detail:   err.Error(),
		})
		return diags
	}
	request.Header.Set("Content-Type", "application/json; charset=utf-8")
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
	session := meta.(*Session)

	request, err := http.NewRequest(http.MethodGet, session.Endpoint+"/"+data.Id(), nil)
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

	character := &MarvelCharacter{}
	json.Unmarshal(bodyBytes, character)
	data.Set(fullNameField, character.FullName)
	data.Set(identityField, character.Identity)
	data.Set(knownasField, character.KnownAs)
	data.Set(typeField, character.Type)

	return diags

}

func characterUpdate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {

	var diags diag.Diagnostics
	session := meta.(*Session)
	character := &MarvelCharacter{
		FullName: data.Get(fullNameField).(string),
		Identity: data.Get(identityField).(string),
		KnownAs:  data.Get(knownasField).(string),
		Type:     data.Get(typeField).(string),
	}

	jsonBody, _ := json.Marshal(character)
	bodyReader := bytes.NewReader(jsonBody)
	request, err := http.NewRequest(http.MethodPut, session.Endpoint+"/"+data.Id(), bodyReader)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failure to create HTTP request",
			Detail:   err.Error(),
		})
		return diags
	}
	request.Header.Set("Content-Type", "application/json; charset=utf-8")
	_, err = session.HttpClient.Do(request)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failure to execute HTTP request",
			Detail:   err.Error(),
		})
		return diags
	}

	return diags

}

func characterDelete(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {

	var diags diag.Diagnostics
	session := meta.(*Session)
	character := &MarvelCharacter{
		FullName: data.Get(fullNameField).(string),
		Identity: data.Get(identityField).(string),
		KnownAs:  data.Get(knownasField).(string),
		Type:     data.Get(typeField).(string),
	}

	jsonBody, _ := json.Marshal(character)
	bodyReader := bytes.NewReader(jsonBody)
	request, err := http.NewRequest(http.MethodDelete, session.Endpoint+"/"+data.Id(), bodyReader)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failure to create HTTP request",
			Detail:   err.Error(),
		})
		return diags
	}
	_, err = session.HttpClient.Do(request)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failure to execute HTTP request",
			Detail:   err.Error(),
		})
		return diags
	}

	data.SetId("")
	return diags

}
