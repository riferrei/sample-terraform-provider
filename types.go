package main

import "net/http"

type SessionKey struct {
	Key int
}

type Session struct {
	Endpoint   string
	HttpClient *http.Client
}

type MarvelCharacter struct {
	ID       string `json:"_id,omitempty"`
	FullName string `json:"fullname,omitempty"`
	Identity string `json:"identity,omitempty"`
	KnownAs  string `json:"knownas,omitempty"`
	Type     string `json:"type,omitempty"`
}

const (
	tokenField    = "token"
	timeoutField  = "timeout"
	fullNameField = "fullname"
	identityField = "identity"
	knownasField  = "knownas"
	typeField     = "type"
)

var (
	sessionKey     = SessionKey{Key: 1}
	characterTypes = []string{"hero", "super-hero", "anti-hero", "villain"}
)
