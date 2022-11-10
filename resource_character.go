package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/opensearch-project/opensearch-go/v2"
	"github.com/opensearch-project/opensearch-go/v2/opensearchapi"
)

func resourceCharacter() *schema.Resource {
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
	backendClient := meta.(*opensearch.Client)

	comicCharacter := &ComicCharacter{
		FullName: data.Get(fullNameField).(string),
		Identity: data.Get(identityField).(string),
		KnownAs:  data.Get(knownasField).(string),
		Type:     data.Get(typeField).(string),
	}
	bodyContent, _ := json.Marshal(comicCharacter)
	bodyReader := bytes.NewReader(bodyContent)

	indexRequest := opensearchapi.IndexRequest{
		Index: backendIndex,
		Body:  bodyReader,
	}
	indexResponse, err := indexRequest.Do(ctx, backendClient)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failure to create character",
			Detail:   "Reason: " + err.Error(),
		})
		return diags
	}

	defer indexResponse.Body.Close()
	bodyContent, err = io.ReadAll(indexResponse.Body)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failure reading response",
			Detail:   "Reason: " + err.Error(),
		})
		return diags
	}

	backendResponse := &BackendResponse{}
	json.Unmarshal(bodyContent, backendResponse)
	data.SetId(backendResponse.ID)

	return diags

}

func characterRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {

	var diags diag.Diagnostics
	backendClient := meta.(*opensearch.Client)
	documentID := data.Id()

	getRequest := opensearchapi.GetRequest{
		Index:      backendIndex,
		DocumentID: documentID,
	}
	getResponse, err := getRequest.Do(ctx, backendClient)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failure to retrieve character",
			Detail:   "Reason: " + err.Error(),
		})
		return diags
	}

	defer getResponse.Body.Close()
	bodyContent, err := io.ReadAll(getResponse.Body)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failure reading response",
			Detail:   "Reason: " + err.Error(),
		})
		return diags
	}

	backendResponse := &BackendResponse{}
	json.Unmarshal(bodyContent, backendResponse)
	data.Set(fullNameField, backendResponse.Source.FullName)
	data.Set(identityField, backendResponse.Source.Identity)
	data.Set(knownasField, backendResponse.Source.KnownAs)
	data.Set(typeField, backendResponse.Source.Type)

	return diags

}

func characterUpdate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {

	var diags diag.Diagnostics
	backendClient := meta.(*opensearch.Client)
	documentID := data.Id()

	updateBody := &struct {
		Doc ComicCharacter `json:"doc,omitempty"`
	}{
		Doc: ComicCharacter{
			FullName: data.Get(fullNameField).(string),
			Identity: data.Get(identityField).(string),
			KnownAs:  data.Get(knownasField).(string),
			Type:     data.Get(typeField).(string),
		},
	}
	bodyContent, _ := json.Marshal(updateBody)
	bodyReader := bytes.NewReader(bodyContent)

	updateRequest := opensearchapi.UpdateRequest{
		Index:      backendIndex,
		DocumentID: documentID,
		Body:       bodyReader,
	}
	_, err := updateRequest.Do(ctx, backendClient)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failure to update character",
			Detail:   "Reason: " + err.Error(),
		})
		return diags
	}

	return diags

}

func characterDelete(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {

	var diags diag.Diagnostics
	backendClient := meta.(*opensearch.Client)
	documentID := data.Id()

	deleteRequest := opensearchapi.DeleteRequest{
		Index:      backendIndex,
		DocumentID: documentID,
	}
	_, err := deleteRequest.Do(ctx, backendClient)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failure to delete character",
			Detail:   "Reason: " + err.Error(),
		})
		return diags
	}

	data.SetId("")
	return diags

}
