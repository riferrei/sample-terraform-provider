package main

type ComicCharacter struct {
	ID       string `json:"_id,omitempty"`
	FullName string `json:"fullname,omitempty"`
	Identity string `json:"identity,omitempty"`
	KnownAs  string `json:"knownas,omitempty"`
	Type     string `json:"type,omitempty"`
}

type BackendResponse struct {
	Index   string          `json:"_index"`
	ID      string          `json:"_id"`
	Version int             `json:"_version"`
	Source  *ComicCharacter `json:"_source"`
}

type BackendSearchResponse struct {
	Hits struct {
		Total struct {
			Value int64 `json:"value"`
		} `json:"total"`
		Hits []*struct {
			ID     string          `json:"_id"`
			Source *ComicCharacter `json:"_source"`
		} `json:"hits"`
	} `json:"hits"`
}
