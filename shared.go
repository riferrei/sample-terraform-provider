package main

const (
	providerName          = "buildonaws"
	dataSourceName        = "_character"
	resourceName          = dataSourceName
	backendAddressField   = "backend_address"
	backendAddressDefault = "http://localhost:9200"
	backendIndex          = "buildonaws"
	idField               = "id"
	fullNameField         = "fullname"
	identityField         = "identity"
	knownasField          = "knownas"
	typeField             = "type"
	lastUpdatedField      = "last_updated"
)

var (
	characterTypes = []string{"hero", "super-hero", "anti-hero", "villain"}
)
