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

var (
	types = []string{"hero", "super-hero", "anti-hero", "villain"}
)

func sampleCharacter() *schema.Resource {
	return &schema.Resource{
		CreateContext: characterCreate,
		ReadContext:   characterRead,
		UpdateContext: characterUpdate,
		DeleteContext: characterDelete,
		Schema: map[string]*schema.Schema{
			"fullname": {
				Type:     schema.TypeString,
				Required: true,
			},
			"identity": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: false,
			},
			"knownas": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
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

type CharacterData struct {
	ID       string `json:"_id,omitempty"`
	FullName string `json:"fullname,omitempty"`
	Identity string `json:"identity,omitempty"`
	KnownAs  string `json:"knownas,omitempty"`
	Type     string `json:"type,omitempty"`
}

func characterCreate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {

	endpoint := meta.(string)
	character := &CharacterData{
		FullName: data.Get("fullname").(string),
		Identity: data.Get("identity").(string),
		KnownAs:  data.Get("knownas").(string),
		Type:     data.Get("type").(string),
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

	createdCharacter := &CharacterData{}
	json.Unmarshal(bodyBytes, createdCharacter)
	data.SetId(createdCharacter.ID)

	return characterRead(ctx, data, meta)

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

	character := &CharacterData{}
	json.Unmarshal(bodyBytes, character)
	data.Set("fullname", character.FullName)
	data.Set("identity", character.Identity)
	data.Set("knownas", character.KnownAs)
	data.Set("type", character.Type)

	return diags

}

func characterUpdate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {

	endpoint := meta.(string) + "/" + data.Id()
	character := &CharacterData{
		FullName: data.Get("fullname").(string),
		Identity: data.Get("identity").(string),
		KnownAs:  data.Get("knownas").(string),
		Type:     data.Get("type").(string),
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
	character := &CharacterData{
		FullName: data.Get("fullname").(string),
		Identity: data.Get("identity").(string),
		KnownAs:  data.Get("knownas").(string),
		Type:     data.Get("type").(string),
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
