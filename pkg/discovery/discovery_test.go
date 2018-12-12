package discovery

import (
	"errors"
	"io/ioutil"
	"testing"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// conditionalityCheckerMock - implements model.ConditionalityChecker interface for tests
type conditionalityCheckerMock struct {
	isPresent           bool
	isPresentErr        error
	missingMandatory    []model.Input
	missingMandatoryErr error
}

// IsOptional - not used in discovery test
func (c conditionalityCheckerMock) IsOptional(method, endpoint string, specification string) (bool, error) {
	return false, nil
}

// Returns IsMandatory true for POST /account-access-consents, false for all other endpoint/methods.
func (c conditionalityCheckerMock) IsMandatory(method, endpoint string, specification string) (bool, error) {
	if method == "POST" && endpoint == "/account-access-consents" {
		return true, nil
	}
	return false, nil
}

// IsConditional - not used in discovery test
func (c conditionalityCheckerMock) IsConditional(method, endpoint string, specification string) (bool, error) {
	return false, nil
}

// IsPresent - returns stubbed isPresent boolean value
func (c conditionalityCheckerMock) IsPresent(method, endpoint string, specification string) (bool, error) {
	return c.isPresent, c.isPresentErr
}

// MissingMandatory - returns stubbed array of missing endpoints
func (c conditionalityCheckerMock) MissingMandatory(endpoints []model.Input, specification string) ([]model.Input, error) {
	return c.missingMandatory, c.missingMandatoryErr
}

// unmarshalDiscoveryJSON - returns discovery model
func testUnmarshalDiscoveryJSON(t *testing.T, discoveryJSON string) *Model {
	t.Helper()

	discovery, err := unmarshalDiscoveryJSON(discoveryJSON)
	assert.NoError(t, err)
	return discovery
}

// discoveryStub - returns discovery JSON with given field stubbed with given value
func discoveryStub(field string, value string) string {
	name := "ob-v3.0-generic"
	description := "An Open Banking UK generic discovery template for v3.0 of Accounts and Payments."
	version := "v0.1.0"
	specName := "Account and Transaction API Specification"
	specURL := "https://openbanking.atlassian.net/wiki/spaces/DZ/pages/642090641/Account+and+Transaction+API+Specification+-+v3.0"
	specVersion := "v3.0"
	schemaVersion := "https://raw.githubusercontent.com/OpenBankingUK/read-write-api-specs/v3.0.0/dist/account-info-swagger.json"
	endpoints := `, "endpoints": [
		{
			"method": "POST",
			"path": "/account-access-consents"
		},
		{
			"method": "GET",
			"path": "/accounts/{AccountId}/balances"
		}
	]`

	switch field {
	case "version":
		version = value
	case "endpoints":
		if value == "" {
			endpoints = ""
		} else {
			endpoints = `, "endpoints": ` + value
		}
	case "specName":
		specName = value
	case "schemaVersion":
		schemaVersion = value
	case "specURL":
		specURL = value
	case "specVersion":
		specVersion = value
	case "name":
		if value == "" {
			name = ""
		}
	case "description":
		if value == "" {
			description = ""
		}
	}

	apiSpecification := `"apiSpecification": {
		"name": "` + specName + `",
		"url": "` + specURL + `",
		"version": "` + specVersion + `",
		"schemaVersion": "` + schemaVersion + `"
	},`
	if field == "apiSpecification" {
		if value == "" {
			apiSpecification = ""
		} else {
			apiSpecification = `"apiSpecification": ` + value + `,`
		}
	}

	discoveryItems := `, "discoveryItems": [
		{
			` + apiSpecification + `
			"openidConfigurationUri": "https://as.aspsp.ob.forgerock.financial/oauth2/.well-known/openid-configuration",
			"resourceBaseUri": "https://rs.aspsp.ob.forgerock.financial:443/"` + endpoints + `
		}
	]`
	if field == "discoveryItems" {
		if value == "" {
			discoveryItems = ""
		} else {
			discoveryItems = `, "discoveryItems": ` + value
		}
	}

	return `
		{
			"discoveryModel": {` +
		`"name": "` + name + `",` +
		`"description": "` + description + `",` +
		`"discoveryVersion": "` + version + `"` +
		discoveryItems + `
			}
	}`
}

// TestValidate - test Validate function
func TestValidate(t *testing.T) {

	// invalidTestCase
	// discoveryJSON - the discovery model JSON
	// failures - list of ValidationFailure structs
	// err - the expected err
	type invalidTest struct {
		discoveryJSON string
		success       bool
		failures      []ValidationFailure
		err           error
	}

	// testValidateFailures - run Validate, and test validation failure expectations
	testValidateFailures := func(t *testing.T, checker model.ConditionalityChecker, expected *invalidTest) {
		t.Helper()

		discovery := testUnmarshalDiscoveryJSON(t, expected.discoveryJSON)
		result, failures, err := Validate(checker, discovery)
		assert.Equal(t, expected.success, result)
		assert.Equal(t, expected.err, err)
		assert.Equal(t, expected.failures, failures)
	}

	t.Run("name missing returns failure", func(t *testing.T) {
		testValidateFailures(t, conditionalityCheckerMock{}, &invalidTest{
			discoveryJSON: discoveryStub("name", ""),
			failures: []ValidationFailure{
				ValidationFailure{
					Key:   "DiscoveryModel.Name",
					Error: "Field validation for 'Name' failed on the 'required' tag",
				},
			},
		})
	})

	t.Run("description missing returns failure", func(t *testing.T) {
		testValidateFailures(t, conditionalityCheckerMock{}, &invalidTest{
			discoveryJSON: discoveryStub("description", ""),
			failures: []ValidationFailure{
				ValidationFailure{
					Key:   "DiscoveryModel.Description",
					Error: "Field validation for 'Description' failed on the 'required' tag",
				},
			},
		})
	})

	t.Run("when discoveryVersion missing returns failure", func(t *testing.T) {
		testValidateFailures(t, conditionalityCheckerMock{}, &invalidTest{
			discoveryJSON: discoveryStub("version", ""),
			failures: []ValidationFailure{
				ValidationFailure{
					Key:   "DiscoveryModel.DiscoveryVersion",
					Error: "Field validation for 'DiscoveryVersion' failed on the 'required' tag",
				},
			}})
	})

	t.Run("when version not in discovery.SupportedVersions() returns failure", func(t *testing.T) {
		testValidateFailures(t, conditionalityCheckerMock{isPresent: true}, &invalidTest{
			discoveryJSON: discoveryStub("version", "v9.9.9"),
			failures: []ValidationFailure{
				ValidationFailure{
					Key:   "DiscoveryModel.DiscoveryVersion",
					Error: "DiscoveryVersion 'v9.9.9' not in list of supported versions",
				},
			}})
	})

	t.Run("when discoveryItems missing returns failure", func(t *testing.T) {
		testValidateFailures(t, conditionalityCheckerMock{}, &invalidTest{
			discoveryJSON: discoveryStub("discoveryItems", ""),
			failures: []ValidationFailure{
				ValidationFailure{
					Key:   "DiscoveryModel.DiscoveryItems",
					Error: "Field validation for 'DiscoveryItems' failed on the 'required' tag",
				},
			}})
	})

	t.Run("when discoveryItems is empty array returns failure", func(t *testing.T) {
		testValidateFailures(t, conditionalityCheckerMock{}, &invalidTest{
			discoveryJSON: discoveryStub("discoveryItems", "[]"),
			failures: []ValidationFailure{
				ValidationFailure{
					Key:   "DiscoveryModel.DiscoveryItems",
					Error: "Field validation for 'DiscoveryItems' failed on the 'gt' tag",
				},
			}})
	})

	t.Run("when discoveryItem missing apiSpecification returns failures", func(t *testing.T) {
		testValidateFailures(t, conditionalityCheckerMock{}, &invalidTest{
			discoveryJSON: discoveryStub("apiSpecification", ""),
			failures: []ValidationFailure{
				ValidationFailure{
					Key:   "DiscoveryModel.DiscoveryItems[0].APISpecification.Name",
					Error: "Field validation for 'Name' failed on the 'required' tag",
				},
				ValidationFailure{
					Key:   "DiscoveryModel.DiscoveryItems[0].APISpecification.URL",
					Error: "Field validation for 'URL' failed on the 'required' tag",
				},
				ValidationFailure{
					Key:   "DiscoveryModel.DiscoveryItems[0].APISpecification.Version",
					Error: "Field validation for 'Version' failed on the 'required' tag",
				},
				ValidationFailure{
					Key:   "DiscoveryModel.DiscoveryItems[0].APISpecification.SchemaVersion",
					Error: "Field validation for 'SchemaVersion' failed on the 'required' tag",
				},
			}})
	})

	t.Run("when discoveryItem apiSpecification schemaVersion not in suite config returns failure", func(t *testing.T) {
		testValidateFailures(t, conditionalityCheckerMock{}, &invalidTest{
			discoveryJSON: discoveryStub("schemaVersion", "http://example.com/bad-schema"),
			failures: []ValidationFailure{
				ValidationFailure{
					Key:   "DiscoveryModel.DiscoveryItems[0].APISpecification.SchemaVersion",
					Error: "'SchemaVersion' not supported by suite 'http://example.com/bad-schema'",
				},
			}})
	})

	t.Run("when discoveryItem apiSpecification schemaVersion in suite config but Name field not matching returns failure", func(t *testing.T) {
		testValidateFailures(t, conditionalityCheckerMock{isPresent: true}, &invalidTest{
			discoveryJSON: discoveryStub("specName", "Bad Spec Name"),
			failures: []ValidationFailure{
				ValidationFailure{
					Key:   "DiscoveryModel.DiscoveryItems[0].APISpecification.Name",
					Error: "'Name' should be 'Account and Transaction API Specification' when schemaVersion is 'https://raw.githubusercontent.com/OpenBankingUK/read-write-api-specs/v3.0.0/dist/account-info-swagger.json'",
				},
			}})
	})

	t.Run("when discoveryItem apiSpecification schemaVersion in suite config but Version field not matching returns failure", func(t *testing.T) {
		testValidateFailures(t, conditionalityCheckerMock{isPresent: true}, &invalidTest{
			discoveryJSON: discoveryStub("specVersion", "v9.9.9"),
			failures: []ValidationFailure{
				ValidationFailure{
					Key:   "DiscoveryModel.DiscoveryItems[0].APISpecification.Version",
					Error: "'Version' should be 'v3.0' when schemaVersion is 'https://raw.githubusercontent.com/OpenBankingUK/read-write-api-specs/v3.0.0/dist/account-info-swagger.json'",
				},
			}})
	})

	t.Run("when discoveryItem apiSpecification schemaVersion in suite config but URL field not matching returns failure", func(t *testing.T) {
		testValidateFailures(t, conditionalityCheckerMock{isPresent: true}, &invalidTest{
			discoveryJSON: discoveryStub("specURL", "http://example.com/bad-url"),
			failures: []ValidationFailure{
				ValidationFailure{
					Key:   "DiscoveryModel.DiscoveryItems[0].APISpecification.URL",
					Error: "'URL' should be 'https://openbanking.atlassian.net/wiki/spaces/DZ/pages/642090641/Account+and+Transaction+API+Specification+-+v3.0' when schemaVersion is 'https://raw.githubusercontent.com/OpenBankingUK/read-write-api-specs/v3.0.0/dist/account-info-swagger.json'",
				},
			}})
	})

	t.Run("when discoveryItems has empty endpoints array returns failure", func(t *testing.T) {
		testValidateFailures(t, conditionalityCheckerMock{}, &invalidTest{
			discoveryJSON: discoveryStub("endpoints", "[]"),
			failures: []ValidationFailure{
				ValidationFailure{
					Key:   "DiscoveryModel.DiscoveryItems[0].Endpoints",
					Error: "Field validation for 'Endpoints' failed on the 'gt' tag",
				},
			}})
	})

	t.Run("when conditionality checker isPresent throws error returns failures", func(t *testing.T) {
		stubIsPresentError := conditionalityCheckerMock{
			isPresent:    false,
			isPresentErr: errors.New("some error message"),
		}
		testValidateFailures(t, stubIsPresentError, &invalidTest{
			discoveryJSON: discoveryStub("", ""),
			failures: []ValidationFailure{
				ValidationFailure{
					Key:   "DiscoveryModel.DiscoveryItems[0].Endpoints[0]",
					Error: "some error message",
				},
				ValidationFailure{
					Key:   "DiscoveryModel.DiscoveryItems[0].Endpoints[1]",
					Error: "some error message",
				},
			},
		},
		)
	})

	t.Run("when conditionality checker reports endpoints not present returns failures", func(t *testing.T) {
		stubAllNotPresent := conditionalityCheckerMock{
			isPresent: false,
		}
		testValidateFailures(t, stubAllNotPresent, &invalidTest{
			discoveryJSON: discoveryStub("", ""),
			failures: []ValidationFailure{
				ValidationFailure{
					Key:   "DiscoveryModel.DiscoveryItems[0].Endpoints[0]",
					Error: "Invalid endpoint Method='POST', Path='/account-access-consents'",
				},
				ValidationFailure{
					Key:   "DiscoveryModel.DiscoveryItems[0].Endpoints[1]",
					Error: "Invalid endpoint Method='GET', Path='/accounts/{AccountId}/balances'",
				},
			},
		})
	})

	t.Run("when conditionality checker reports endpoints present returns no failures", func(t *testing.T) {
		stubAllPresent := conditionalityCheckerMock{
			isPresent: true,
		}
		testValidateFailures(t, stubAllPresent, &invalidTest{
			discoveryJSON: discoveryStub("", ""),
			success:       true,
			failures:      []ValidationFailure{},
		})
	})

	t.Run("when conditionality checker missingMandatory throws error returns failures", func(t *testing.T) {
		stubIsPresentError := conditionalityCheckerMock{
			isPresent:           true,
			missingMandatory:    []model.Input{},
			missingMandatoryErr: errors.New("the error message"),
		}
		testValidateFailures(t, stubIsPresentError, &invalidTest{
			discoveryJSON: discoveryStub("", ""),
			failures: []ValidationFailure{
				ValidationFailure{
					Key:   "DiscoveryModel.DiscoveryItems[0].Endpoints",
					Error: "the error message",
				},
			},
		},
		)
	})

	t.Run("when conditionality checker reports missing mandatory endpoints returns failures", func(t *testing.T) {
		stubMissingMandatory := conditionalityCheckerMock{
			isPresent: true,
			missingMandatory: []model.Input{
				model.Input{Method: "GET", Endpoint: "/account-access-consents/{ConsentId}"},
				model.Input{Method: "DELETE", Endpoint: "/account-access-consents/{ConsentId}"},
			},
		}
		testValidateFailures(t, stubMissingMandatory, &invalidTest{
			discoveryJSON: discoveryStub("", ""),
			failures: []ValidationFailure{
				ValidationFailure{
					Key:   "DiscoveryModel.DiscoveryItems[0].Endpoints",
					Error: "Missing mandatory endpoint Method='GET', Path='/account-access-consents/{ConsentId}'",
				},
				ValidationFailure{
					Key:   "DiscoveryModel.DiscoveryItems[0].Endpoints",
					Error: "Missing mandatory endpoint Method='DELETE', Path='/account-access-consents/{ConsentId}'",
				},
			},
		})
	})
}

func TestDiscovery_FromJSONString_Valid(t *testing.T) {
	discoveryExample, err := ioutil.ReadFile("./templates/ob-v3.0-ozone.json")
	require.NoError(t, err)
	require.NotNil(t, discoveryExample)
	config := string(discoveryExample)

	accountAPIDiscoveryItem := ModelDiscoveryItem{
		APISpecification: ModelAPISpecification{
			Name:          "Account and Transaction API Specification",
			URL:           "https://openbanking.atlassian.net/wiki/spaces/DZ/pages/642090641/Account+and+Transaction+API+Specification+-+v3.0",
			Version:       "v3.0",
			SchemaVersion: "https://raw.githubusercontent.com/OpenBankingUK/read-write-api-specs/v3.0.0/dist/account-info-swagger.json",
		},
		OpenidConfigurationURI: "https://modelobankauth2018.o3bank.co.uk:4101/.well-known/openid-configuration",
		ResourceBaseURI:        "https://modelobank2018.o3bank.co.uk:4501/open-banking/v3.0/",
		Endpoints: []ModelEndpoint{
			ModelEndpoint{
				Method:                "POST",
				Path:                  "/account-access-consents",
				ConditionalProperties: []ModelConditionalProperties(nil),
			},
			ModelEndpoint{
				Method:                "GET",
				Path:                  "/account-access-consents/{ConsentId}",
				ConditionalProperties: []ModelConditionalProperties(nil),
			},
			ModelEndpoint{Method: "DELETE",
				Path:                  "/account-access-consents/{ConsentId}",
				ConditionalProperties: []ModelConditionalProperties(nil),
			},
			ModelEndpoint{Method: "GET",
				Path:                  "/accounts/{AccountId}/product",
				ConditionalProperties: []ModelConditionalProperties(nil),
			},
			ModelEndpoint{Method: "GET",
				Path: "/accounts/{AccountId}/transactions",
				ConditionalProperties: []ModelConditionalProperties{
					ModelConditionalProperties{
						Schema:   "OBTransaction3Detail",
						Property: "Balance",
						Path:     "Data.Transaction.*.Balance",
					},
					ModelConditionalProperties{
						Schema:   "OBTransaction3Detail",
						Property: "MerchantDetails",
						Path:     "Data.Transaction.*.MerchantDetails",
					},
					ModelConditionalProperties{
						Schema:   "OBTransaction3Basic",
						Property: "TransactionReference",
						Path:     "Data.Transaction.*.TransactionReference",
					},
					ModelConditionalProperties{
						Schema:   "OBTransaction3Detail",
						Property: "TransactionReference",
						Path:     "Data.Transaction.*.TransactionReference",
					},
				},
			},
			ModelEndpoint{
				Method:                "GET",
				Path:                  "/accounts",
				ConditionalProperties: []ModelConditionalProperties(nil),
			},
			ModelEndpoint{
				Method:                "GET",
				Path:                  "/accounts/{AccountId}",
				ConditionalProperties: []ModelConditionalProperties(nil),
			},
			ModelEndpoint{
				Method:                "GET",
				Path:                  "/accounts/{AccountId}/balances",
				ConditionalProperties: []ModelConditionalProperties(nil),
			},
		},
	}

	modelActual, err := unmarshalDiscoveryJSON(config)
	assert.NoError(t, err)
	assert.NotNil(t, modelActual.DiscoveryModel)
	discoveryModel := modelActual.DiscoveryModel

	t.Run("model has a version", func(t *testing.T) {
		assert := assert.New(t)
		assert.Equal(discoveryModel.DiscoveryVersion, "v0.1.0")
	})

	t.Run("model has correct number of discovery items", func(t *testing.T) {
		assert := assert.New(t)
		assert.Equal(len(discoveryModel.DiscoveryItems), 2)
	})

	t.Run("model has correct discovery item contents", func(t *testing.T) {
		assert := assert.New(t)
		assert.Equal(accountAPIDiscoveryItem, discoveryModel.DiscoveryItems[0])
	})
}

func TestDiscovery_Version(t *testing.T) {
	assert := assert.New(t)

	assert.Equal(Version(), "v0.1.0")
}
