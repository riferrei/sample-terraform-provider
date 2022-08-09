package main

const (
	fullNameField = "fullname"
	identityField = "identity"
	knownasField  = "knownas"
	typeField     = "type"
)

var (
	types = []string{"hero", "super-hero", "anti-hero", "villain"}
)

type MarvelCharacter struct {
	ID       string `json:"_id,omitempty"`
	FullName string `json:"fullname,omitempty"`
	Identity string `json:"identity,omitempty"`
	KnownAs  string `json:"knownas,omitempty"`
	Type     string `json:"type,omitempty"`
}
