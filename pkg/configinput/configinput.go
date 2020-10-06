package configinput

import (
	"strings"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/discovery"
)

// Field describes what input rules apply to a configuration field
type Field struct {
	FieldNameJSON string `json:"field_name,omitempty"`
	Rule          Rule   `json:"validation,omitempty"`
}

// Rule is a suggestion for validation mechanisms about what checks to perform
// if FE validation is to be performed on the field. The Required, Match, In fields are
// evaluated as AND, e.g. if multiple are set then all set conditions are expeceted to
// be verified.
type Rule struct {
	// Required - Does not mean that the value is required.
	// It means that the field is required to be displayed on the input form.
	Required bool `json:"required,omitempty"`
	// Match regular expression
	Match string `json:"match,omitempty"`
	// In one of the provided values
	In []interface{} `json:"in,omitempty"`
	// Custom validation function describes custom validation mechanisms (lt, gt, if, etc.)
	Condition *CustomCondition `json:"condition,omitempty"`
}

// CustomCondition describes special cases
type CustomCondition struct {
	Op     string   `json:"op"`
	Params []string `json:"params,omitempty"`
}

// RuleSet holds all the rules which are required to render an input form for a GlobalConfguration struct
// type RuleSet struct {
// 	// keys are the FieldNameJSON values
// 	Rules map[string][]validation.Rule
// }

// FieldRules takes a populated discovery struct and returns a list
// of Field items which describe the input fields required by the scenarios
// defined in the discovery model.
func FieldRules(discovery *discovery.ModelDiscovery) []Field {
	configFields := append(clientRequirements, wellKnownRequirements...)
	for _, spec := range discovery.DiscoveryItems {
		for _, endpoint := range spec.Endpoints {
			key := strings.Join([]string{endpoint.Method, endpoint.Path}, " ")
			endpointFields, isSet := endpointRequirements[key]
			if !isSet {
				continue
			}

			configFields = mergeConfigFields(configFields, endpointFields...)
		}
	}
	return configFields
}

func mergeConfigFields(existingFields []Field, newFields ...Field) []Field {
NextNew:
	for _, newField := range newFields {
		for _, ef := range existingFields {
			if ef.FieldNameJSON == newField.FieldNameJSON {
				// TODO: update the existing if the new is stricter (not a use case now, but the right thing to do)
				continue NextNew
			}
		}

		// append new item
		existingFields = append(existingFields, newField)
	}
	return existingFields
}

var (
	// collected configuration fields on the client level
	clientRequirements = []Field{
		{"signing_private", Rule{Required: true}},
		{"signing_public", Rule{Required: true}},
		{"transport_private", Rule{Required: true}},
		{"transport_public", Rule{Required: true}},
		{"use_eidas_cert", Rule{Required: true}},
		{"client_id", Rule{Required: true}},
		{"x_fapi_financial_id", Rule{Required: true}},
		{"x_fapi_customer_ip_address", Rule{Required: true}}, // Optional, but needs to be rendered.
		{"transaction_from_date", Rule{Required: true}},      // Only when headless is used.. tbd: what cases use headless?
		{"transaction_to_date", Rule{Required: true}},        // ditto
		{"use_eidas_cert", Rule{Required: true}},
		{"eidas_signing_kid", Rule{Condition: &CustomCondition{"if", []string{"use_eidas_cert"}}}},
		{"eidas_issuer", Rule{Condition: &CustomCondition{"if", []string{"use_eidas_cert"}}}},
	}

	// requirements for interacting with the well known endpoints
	wellKnownRequirements = []Field{
		{"token_endpoint", Rule{Required: true}},
		{"response_type", Rule{Required: true}},
		{"token_endpoint_auth_method", Rule{Required: true}},
		{"request_object_signing_alg", Rule{Required: true}},
		{"authorization_endpoint", Rule{Required: true}},
		{"resource_base_url", Rule{Required: true}},
		{"issuer", Rule{Required: true}},
	}

	// endpointRequirements contains additional configuration requirements
	// of the Open Banking RW API endpoints. The key is constructed as METHOD PATH.
	// The PATH folloews the syntax in the discovery file.
	endpointRequirements = map[string][]Field{

		//
		// Account and Transaction API endpoints:
		//

		"POST /account-access-consents":               {},
		"GET /account-access-consents/{ConsentId}":    {},
		"DELETE /account-access-consents/{ConsentId}": {},
		"GET /accounts": {},

		"GET /accounts/{AccountId}": {
			{"resource_ids.account_ids", Rule{Required: true}},
		},
		"GET /accounts/{AccountId}/balances": {
			{"resource_ids.account_ids", Rule{Required: true}},
		},
		"GET /accounts/{AccountId}/beneficiaries": {
			{"resource_ids.account_ids", Rule{Required: true}},
		},
		"GET /accounts/{AccountId}/direct-debits": {
			{"resource_ids.account_ids", Rule{Required: true}},
		},
		"GET /accounts/{AccountId}/offers": {
			{"resource_ids.account_ids", Rule{Required: true}},
		},
		"GET /accounts/{AccountId}/party": {
			{"resource_ids.account_ids", Rule{Required: true}},
		},
		"GET /accounts/{AccountId}/product": {
			{"resource_ids.account_ids", Rule{Required: true}},
		},
		"GET /accounts/{AccountId}/scheduled-payments": {
			{"resource_ids.account_ids", Rule{Required: true}},
		},
		"GET /accounts/{AccountId}/standing-orders": {
			{"resource_ids.account_ids", Rule{Required: true}},
		},
		"GET /accounts/{AccountId}/statements": {
			{"resource_ids.account_ids", Rule{Required: true}},
		},
		"GET /accounts/{AccountId}/statements/{StatementId}": {
			{"resource_ids.account_ids", Rule{Required: true}},
			{"resource_ids.statement_ids", Rule{Required: true}},
		},
		"GET /accounts/{AccountId}/statements/{StatementId}/file": {
			{"resource_ids.account_ids", Rule{Required: true}},
			{"resource_ids.statement_ids", Rule{Required: true}},
		},
		"GET /accounts/{AccountId}/statements/{StatementId}/transactions": {
			{"resource_ids.account_ids", Rule{Required: true}},
			{"resource_ids.statement_ids", Rule{Required: true}},
		},
		"GET /accounts/{AccountId}/transactions": {
			{"resource_ids", Rule{Required: true}},
		},
		"GET /balances":           {},
		"GET /beneficiaries":      {},
		"GET /direct-debits":      {},
		"GET /offers":             {},
		"GET /party":              {},
		"GET /products":           {},
		"GET /scheduled-payments": {},
		"GET /standing-orders":    {},
		"GET /statements":         {},
		"GET /transactions":       {},

		//
		// Payment Initiation API endpoints
		//

		"POST /domestic-payment-consents": {
			{"creditor_account", Rule{Required: true}},
			{"instructed_amount", Rule{Required: true}},
		},
		"GET /domestic-payment-consents/{ConsentId}":                    {},
		"GET /domestic-payment-consents/{ConsentId}/funds-confirmation": {},
		"POST /domestic-payments": {
			{"creditor_account", Rule{Required: true}},
			{"instructed_amount", Rule{Required: true}},
		},
		"GET /domestic-payments/{DomesticPaymentId}": {},
		"POST /domestic-scheduled-payment-consents": {
			{"creditor_account", Rule{Required: true}},
			{"instructed_amount", Rule{Required: true}},
			{"requested_execution_date_time", Rule{Required: true}},
		},
		"GET /domestic-scheduled-payment-consents/{ConsentId}": {},
		"POST /domestic-scheduled-payments": {
			{"creditor_account", Rule{Required: true}},
			{"instructed_amount", Rule{Required: true}},
			{"requested_execution_date_time", Rule{Required: true}},
		},
		"GET /domestic-scheduled-payments/{DomesticScheduledPaymentId}": {},
		"POST /domestic-standing-order-consents": {
			{"creditor_account", Rule{Required: true}},
			{"instructed_amount", Rule{Required: true}},
			{"payment_frequency", Rule{Required: true}},
			{"first_payment_date_time", Rule{Required: true}},
		},
		"GET /domestic-standing-order-consents/{ConsentId}": {},
		"POST /domestic-standing-orders": {
			{"creditor_account", Rule{Required: true}},
			{"instructed_amount", Rule{Required: true}},
			{"payment_frequency", Rule{Required: true}},
			{"first_payment_date_time", Rule{Required: true}},
		},
		"GET /domestic-standing-orders/{DomesticStandingOrderId}": {},
		"POST /international-payment-consents": {
			{"international_creditor_account", Rule{Required: true}},
			{"currency_of_transfer", Rule{Required: true}},
			{"instructed_amount", Rule{Required: true}},
		},
		"GET /international-payment-consents/{ConsentId}":                    {},
		"GET /international-payment-consents/{ConsentId}/funds-confirmation": {},
		"POST /international-payments": {
			{"international_creditor_account", Rule{Required: true}},
			{"currency_of_transfer", Rule{Required: true}},
			{"instructed_amount", Rule{Required: true}},
		},
		"GET /international-payments/{InternationalPaymentId}": {},
		"POST /international-scheduled-payment-consents": {
			{"international_creditor_account", Rule{Required: true}},
			{"currency_of_transfer", Rule{Required: true}},
			{"instructed_amount", Rule{Required: true}},
			{"requested_execution_date_time", Rule{Required: true}},
		},
		"GET /international-scheduled-payment-consents/{ConsentId}":                    {},
		"GET /international-scheduled-payment-consents/{ConsentId}/funds-confirmation": {},
		"POST /international-scheduled-payments": {
			{"international_creditor_account", Rule{Required: true}},
			{"currency_of_transfer", Rule{Required: true}},
			{"instructed_amount", Rule{Required: true}},
			{"requested_execution_date_time", Rule{Required: true}},
		},
		"GET /international-scheduled-payments/{InternationalScheduledPaymentId}": {},
		"POST /international-standing-order-consents": {
			{"international_creditor_account", Rule{Required: true}},
			{"currency_of_transfer", Rule{Required: true}},
			{"instructed_amount", Rule{Required: true}},
			{"payment_frequency", Rule{Required: true}},
			{"first_payment_date_time", Rule{Required: true}},
		},
		"GET /international-standing-order-consents/{ConsentId}": {},
		"POST /international-standing-orders": {
			{"international_creditor_account", Rule{Required: true}},
			{"currency_of_transfer", Rule{Required: true}},
			{"instructed_amount", Rule{Required: true}},
			{"payment_frequency", Rule{Required: true}},
			{"first_payment_date_time", Rule{Required: true}},
		},
		"GET /international-standing-orders/{InternationalStandingOrderPaymentId}": {},
		"POST /file-payment-consents":                    {},
		"GET /file-payment-consents/{ConsentId}":         {},
		"POST /file-payment-consents/{ConsentId}/file":   {},
		"GET /file-payment-consents/{ConsentId}/file":    {},
		"POST /file-payments":                            {},
		"GET /file-payments/{FilePaymentId}":             {},
		"GET /file-payments/{FilePaymentId}/report-file": {},

		//
		// Confirmation Of Funds API endpoints
		//

		"POST /funds-confirmation-consents": {
			{"cbpii_debtor_account", Rule{Required: true}},
		},
		"GET /funds-confirmation-consents/{ConsentId}":    {},
		"DELETE /funds-confirmation-consents/{ConsentId}": {},
		"POST /funds-confirmations": {
			{"instructed_amount", Rule{Required: true}},
		},
	}
)

func newField(fieldName string, rule Rule) Field {
	return Field{
		FieldNameJSON: fieldName,
		Rule:          rule,
	}
}
