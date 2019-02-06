package model

import (
	"errors"
	"net/url"
)

// Specification - Represents OB API specification.
// Fields are from the APIReference JSON-LD schema, see: https://schema.org/APIReference
type Specification struct {
	Identifier string
	Name       string
	// URL of confluence specifications file.
	URL *url.URL
	// Version of the specifications
	Version string
	// URL of OpenAPI/Swagger specifications file.
	SchemaVersion *url.URL
}

var (
	specifications = []Specification{
		{
			Identifier:    "account-transaction-v3.1",
			Name:          "Account and Transaction API Specification",
			URL:           mustParseURL("https://openbanking.atlassian.net/wiki/spaces/DZ/pages/937820271/Account+and+Transaction+API+Specification+-+v3.1"),
			Version:       "v3.1",
			SchemaVersion: mustParseURL("https://raw.githubusercontent.com/OpenBankingUK/read-write-api-specs/v3.1.0/dist/account-info-swagger.json"),
		},
		{
			Identifier:    "payment-initiation-v3.1",
			Name:          "Payment Initiation API",
			URL:           mustParseURL("https://openbanking.atlassian.net/wiki/spaces/DZ/pages/937754701/Payment+Initiation+API+Specification+-+v3.1"),
			Version:       "v3.1",
			SchemaVersion: mustParseURL("https://raw.githubusercontent.com/OpenBankingUK/read-write-api-specs/v3.1.0/dist/payment-initiation-swagger.json"),
		},
		{
			Identifier:    "confirmation-funds-v3.1",
			Name:          "Confirmation of Funds API Specification",
			URL:           mustParseURL("https://openbanking.atlassian.net/wiki/spaces/DZ/pages/937951380/Confirmation+of+Funds+API+Specification+-+v3.1"),
			Version:       "v3.1",
			SchemaVersion: mustParseURL("https://raw.githubusercontent.com/OpenBankingUK/read-write-api-specs/v3.1.0/dist/confirmation-funds-swagger.json"),
		},
		{
			Identifier:    "event-notification-aspsp-v3.1",
			Name:          "Event Notification API Specification - ASPSP Endpoints",
			URL:           mustParseURL("https://openbanking.atlassian.net/wiki/spaces/DZ/pages/937951397/Event+Notification+API+Specification+-+v3.1"),
			Version:       "v3.1",
			SchemaVersion: mustParseURL("https://raw.githubusercontent.com/OpenBankingUK/read-write-api-specs/v3.1.0/dist/callback-urls-swagger.yaml"),
		},
		{
			Identifier:    "event-notification-tpp-v3.1",
			Name:          "Event Notification API Specification - TPP Endpoints",
			URL:           mustParseURL("https://openbanking.atlassian.net/wiki/spaces/DZ/pages/937951397/Event+Notification+API+Specification+-+v3.1"),
			Version:       "v3.1",
			SchemaVersion: mustParseURL("https://raw.githubusercontent.com/OpenBankingUK/read-write-api-specs/v3.1.0/dist/event-notifications-swagger.json"),
		},
	}
)

// Specifications - get a clone of the `specifications` array.
func Specifications() []Specification {
	clone := make([]Specification, len(specifications))
	copy(clone, specifications)
	return clone
}

// SpecificationFromSchemaVersion - returns specification struct
// for given schema version URL, or nil when there is no match.
func SpecificationFromSchemaVersion(schemaVersion string) (Specification, error) {
	var spec Specification
	for _, specification := range specifications {
		if specification.SchemaVersion.String() == schemaVersion {
			return specification, nil
		}
	}
	return spec, errors.New("no specifications found for schema version: " + schemaVersion)
}

func mustParseURL(rawurl string) *url.URL {
	parsedUrl, err := url.Parse(rawurl)
	if err != nil {
		panic(rawurl)
	}
	return parsedUrl
}
