package server

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

type testTableItem struct {
	label                string
	endpoint             string
	httpStatusExpected   int
	requestBody          string
	responseBodyExpected string
}

func TestRedirectHandlersFragmentOK(t *testing.T) {
	require := require.New(t)

	server := NewServer(nullLogger(), conditionalityCheckerMock{}, mockVersionChecker())
	defer func() {
		require.NoError(server.Shutdown(context.TODO()))
	}()
	require.NotNil(server)

	// Valid code and c_hash combination
	testItemOK := testTableItem{
		label:                "fragment ok",
		endpoint:             `/api/redirect/fragment/ok`,
		httpStatusExpected:   http.StatusOK,
		responseBodyExpected: "null",
		requestBody: `
{
    "code": "a052c795-742d-415a-843f-8a4939d740d1",
    "scope": "openid accounts",
    "id_token": "eyJ0eXAiOiJKV1QiLCJraWQiOiJGb2w3SXBkS2VMWm16S3RDRWdpMUxEaFNJek09IiwiYWxnIjoiRVMyNTYifQ.eyJzdWIiOiJtYmFuYSIsImF1ZGl0VHJhY2tpbmdJZCI6IjY5YzZkZmUzLWM4MDEtNGRkMi05Mjc1LTRjNWVhNzdjZWY1NS0xMDMzMDgyIiwiaXNzIjoiaHR0cHM6Ly9tYXRscy5hcy5hc3BzcC5vYi5mb3JnZXJvY2suZmluYW5jaWFsL29hdXRoMi9vcGVuYmFua2luZyIsInRva2VuTmFtZSI6ImlkX3Rva2VuIiwibm9uY2UiOiI1YTZiMGQ3ODMyYTlmYjRmODBmMTE3MGEiLCJhY3IiOiJ1cm46b3BlbmJhbmtpbmc6cHNkMjpzY2EiLCJhdWQiOiI1NGY2NDMwOS00MzNkLTQ2MTAtOTVkMi02M2QyZjUyNTM0MTIiLCJjX2hhc2giOiIxbGt1SEFuaVJDZlZNS2xEc0pxTTNBIiwib3BlbmJhbmtpbmdfaW50ZW50X2lkIjoiQTY5MDA3Nzc1LTcwZGQtNGIyMi1iZmM1LTlkNTI0YTkxZjk4MCIsInNfaGFzaCI6ImZ0OWRrQTdTWXdlb2hlZXpjOGFHeEEiLCJhenAiOiI1NGY2NDMwOS00MzNkLTQ2MTAtOTVkMi02M2QyZjUyNTM0MTIiLCJhdXRoX3RpbWUiOjE1Mzk5NDM3NzUsInJlYWxtIjoiL29wZW5iYW5raW5nIiwiZXhwIjoxNTQwMDMwMTgxLCJ0b2tlblR5cGUiOiJKV1RUb2tlbiIsImlhdCI6MTUzOTk0Mzc4MX0.8bm69KPVQIuvcTlC-p0FGcplTV1LnmtacHybV2PTb2uEgMgrL3JNA0jpT2OYO73r3zPC41mNQlMDvVOUn78osQ",
    "state": "5a6b0d7832a9fb4f80f1170a"
}
	`,
	}

	// c_hash will not match calculated c_hash due to invalid `code`
	testItemInvalidCode := testTableItem{
		label:                "fragment invalid code",
		endpoint:             "/api/redirect/fragment/ok",
		httpStatusExpected:   http.StatusBadRequest,
		responseBodyExpected: `{"error":"c_hash invalid"}`,
		requestBody: `
{
    "code": "---invalid code---",
    "scope": "openid accounts",
    "id_token": "eyJ0eXAiOiJKV1QiLCJraWQiOiJGb2w3SXBkS2VMWm16S3RDRWdpMUxEaFNJek09IiwiYWxnIjoiRVMyNTYifQ.eyJzdWIiOiJtYmFuYSIsImF1ZGl0VHJhY2tpbmdJZCI6IjY5YzZkZmUzLWM4MDEtNGRkMi05Mjc1LTRjNWVhNzdjZWY1NS0xMDMzMDgyIiwiaXNzIjoiaHR0cHM6Ly9tYXRscy5hcy5hc3BzcC5vYi5mb3JnZXJvY2suZmluYW5jaWFsL29hdXRoMi9vcGVuYmFua2luZyIsInRva2VuTmFtZSI6ImlkX3Rva2VuIiwibm9uY2UiOiI1YTZiMGQ3ODMyYTlmYjRmODBmMTE3MGEiLCJhY3IiOiJ1cm46b3BlbmJhbmtpbmc6cHNkMjpzY2EiLCJhdWQiOiI1NGY2NDMwOS00MzNkLTQ2MTAtOTVkMi02M2QyZjUyNTM0MTIiLCJjX2hhc2giOiIxbGt1SEFuaVJDZlZNS2xEc0pxTTNBIiwib3BlbmJhbmtpbmdfaW50ZW50X2lkIjoiQTY5MDA3Nzc1LTcwZGQtNGIyMi1iZmM1LTlkNTI0YTkxZjk4MCIsInNfaGFzaCI6ImZ0OWRrQTdTWXdlb2hlZXpjOGFHeEEiLCJhenAiOiI1NGY2NDMwOS00MzNkLTQ2MTAtOTVkMi02M2QyZjUyNTM0MTIiLCJhdXRoX3RpbWUiOjE1Mzk5NDM3NzUsInJlYWxtIjoiL29wZW5iYW5raW5nIiwiZXhwIjoxNTQwMDMwMTgxLCJ0b2tlblR5cGUiOiJKV1RUb2tlbiIsImlhdCI6MTUzOTk0Mzc4MX0.8bm69KPVQIuvcTlC-p0FGcplTV1LnmtacHybV2PTb2uEgMgrL3JNA0jpT2OYO73r3zPC41mNQlMDvVOUn78osQ",
    "state": "5a6b0d7832a9fb4f80f1170a"
}
	`,
	}

	// Note c_hash is manipulated (invalid) in JWT, meaning invalid signature as a result
	// if signature validation is implemented this test shall fail.
	testItemInvalidCHash := testTableItem{
		label:                "fragment invalid c_hash",
		endpoint:             "/api/redirect/fragment/ok",
		httpStatusExpected:   http.StatusBadRequest,
		responseBodyExpected: `{"error":"c_hash invalid"}`,
		requestBody: `
{
    "code": "80bf17a3-e617-4983-9d62-b50bd8e6fce4",
    "scope": "openid accounts",
    "id_token": "eyJhbGciOiJQUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJodHRwczovL2p3dC1pZHAuZXhhbXBsZS5jb20iLCJzdWIiOiJtYWlsdG86bWlrZUBleGFtcGxlLmNvbSIsIm5iZiI6MTU0OTU1NjY0MiwiZXhwIjoxNTQ5NTYwMjQyLCJpYXQiOjE1NDk1NTY2NDIsImp0aSI6ImlkMTIzNDU2IiwidHlwIjoiaHR0cHM6Ly9leGFtcGxlLmNvbS9yZWdpc3RlciIsICJjX2hhc2giOiAiaW52YWxpZC1jX2hhc2gifQ.CG8bb_wT7EetLsE3SB8W30K3z_be14ZNXsjWQiklXkqImE-aWFvCqruwh-3aAG5xvaQ_u6T5mj7jaK-ZX93v591FsMPmX1MWyUYNfJp5MsPsWUfzZX69Us5UAqOgZ2zxu662prcE8fVqsL-GB-boVR_0e1SUj4NjKhiEHCNVYe-SGclRSZPvjRf0ymQBacmFz84kLqFVZYTXJFkufXd09FUopnNVK-aK2aCc39TaxzxFwwLaAW_iOtJnzHUtnNdF1OUW5MLTeJYd7hPg0Oq5hPUtz2h2XLVl76ERdJYNWNa1yws4gaWE9PaDgNu-mYbfEZVIHnB28XkB7d6BaCW-GQ",
    "state": "5a6b0d7832a9fb4f80f1170a"
}
	`,
	}

	ttData := []testTableItem{
		testItemOK,
		testItemInvalidCode,
		testItemInvalidCHash,
	}

	for _, ttItem := range ttData {
		// read the file we expect to be served.
		code, body, headers := request(
			http.MethodPost,
			ttItem.endpoint,
			strings.NewReader(ttItem.requestBody),
			server,
		)

		// do assertions.
		require.Equal(ttItem.httpStatusExpected, code, ttItem.label)
		require.Len(headers, 2, ttItem.label)
		require.Equal("application/json; charset=UTF-8", headers["Content-Type"][0], ttItem.label)
		require.NotNil(body, ttItem.label)

		bodyActual := body.String()
		require.Equal(ttItem.responseBodyExpected, bodyActual, ttItem.label)
	}
}

func TestRedirectHandlersQueryOK(t *testing.T) {
	require := require.New(t)

	server := NewServer(nullLogger(), conditionalityCheckerMock{}, mockVersionChecker())
	defer func() {
		require.NoError(server.Shutdown(context.TODO()))
	}()
	require.NotNil(server)

	// valid code and c_hash combination
	testItemOK := testTableItem{
		label:                "query ok",
		endpoint:             `/api/redirect/query/ok`,
		httpStatusExpected:   http.StatusOK,
		responseBodyExpected: "null",
		requestBody: `
{
    "code": "a052c795-742d-415a-843f-8a4939d740d1",
    "scope": "openid accounts",
    "id_token": "eyJ0eXAiOiJKV1QiLCJraWQiOiJGb2w3SXBkS2VMWm16S3RDRWdpMUxEaFNJek09IiwiYWxnIjoiRVMyNTYifQ.eyJzdWIiOiJtYmFuYSIsImF1ZGl0VHJhY2tpbmdJZCI6IjY5YzZkZmUzLWM4MDEtNGRkMi05Mjc1LTRjNWVhNzdjZWY1NS0xMDMzMDgyIiwiaXNzIjoiaHR0cHM6Ly9tYXRscy5hcy5hc3BzcC5vYi5mb3JnZXJvY2suZmluYW5jaWFsL29hdXRoMi9vcGVuYmFua2luZyIsInRva2VuTmFtZSI6ImlkX3Rva2VuIiwibm9uY2UiOiI1YTZiMGQ3ODMyYTlmYjRmODBmMTE3MGEiLCJhY3IiOiJ1cm46b3BlbmJhbmtpbmc6cHNkMjpzY2EiLCJhdWQiOiI1NGY2NDMwOS00MzNkLTQ2MTAtOTVkMi02M2QyZjUyNTM0MTIiLCJjX2hhc2giOiIxbGt1SEFuaVJDZlZNS2xEc0pxTTNBIiwib3BlbmJhbmtpbmdfaW50ZW50X2lkIjoiQTY5MDA3Nzc1LTcwZGQtNGIyMi1iZmM1LTlkNTI0YTkxZjk4MCIsInNfaGFzaCI6ImZ0OWRrQTdTWXdlb2hlZXpjOGFHeEEiLCJhenAiOiI1NGY2NDMwOS00MzNkLTQ2MTAtOTVkMi02M2QyZjUyNTM0MTIiLCJhdXRoX3RpbWUiOjE1Mzk5NDM3NzUsInJlYWxtIjoiL29wZW5iYW5raW5nIiwiZXhwIjoxNTQwMDMwMTgxLCJ0b2tlblR5cGUiOiJKV1RUb2tlbiIsImlhdCI6MTUzOTk0Mzc4MX0.8bm69KPVQIuvcTlC-p0FGcplTV1LnmtacHybV2PTb2uEgMgrL3JNA0jpT2OYO73r3zPC41mNQlMDvVOUn78osQ",
    "state": "5a6b0d7832a9fb4f80f1170a"
}
	`,
	}

	// invalid value for `code`
	testItemInvalidCode := testTableItem{
		label:                "query invalid code",
		endpoint:             `/api/redirect/query/ok`,
		httpStatusExpected:   http.StatusBadRequest,
		responseBodyExpected: `{"error":"c_hash invalid"}`,
		requestBody: `
{
    "code": "---invalid code---",
    "scope": "openid accounts",
    "id_token": "eyJ0eXAiOiJKV1QiLCJraWQiOiJGb2w3SXBkS2VMWm16S3RDRWdpMUxEaFNJek09IiwiYWxnIjoiRVMyNTYifQ.eyJzdWIiOiJtYmFuYSIsImF1ZGl0VHJhY2tpbmdJZCI6IjY5YzZkZmUzLWM4MDEtNGRkMi05Mjc1LTRjNWVhNzdjZWY1NS0xMDMzMDgyIiwiaXNzIjoiaHR0cHM6Ly9tYXRscy5hcy5hc3BzcC5vYi5mb3JnZXJvY2suZmluYW5jaWFsL29hdXRoMi9vcGVuYmFua2luZyIsInRva2VuTmFtZSI6ImlkX3Rva2VuIiwibm9uY2UiOiI1YTZiMGQ3ODMyYTlmYjRmODBmMTE3MGEiLCJhY3IiOiJ1cm46b3BlbmJhbmtpbmc6cHNkMjpzY2EiLCJhdWQiOiI1NGY2NDMwOS00MzNkLTQ2MTAtOTVkMi02M2QyZjUyNTM0MTIiLCJjX2hhc2giOiIxbGt1SEFuaVJDZlZNS2xEc0pxTTNBIiwib3BlbmJhbmtpbmdfaW50ZW50X2lkIjoiQTY5MDA3Nzc1LTcwZGQtNGIyMi1iZmM1LTlkNTI0YTkxZjk4MCIsInNfaGFzaCI6ImZ0OWRrQTdTWXdlb2hlZXpjOGFHeEEiLCJhenAiOiI1NGY2NDMwOS00MzNkLTQ2MTAtOTVkMi02M2QyZjUyNTM0MTIiLCJhdXRoX3RpbWUiOjE1Mzk5NDM3NzUsInJlYWxtIjoiL29wZW5iYW5raW5nIiwiZXhwIjoxNTQwMDMwMTgxLCJ0b2tlblR5cGUiOiJKV1RUb2tlbiIsImlhdCI6MTUzOTk0Mzc4MX0.8bm69KPVQIuvcTlC-p0FGcplTV1LnmtacHybV2PTb2uEgMgrL3JNA0jpT2OYO73r3zPC41mNQlMDvVOUn78osQ",
    "state": "5a6b0d7832a9fb4f80f1170a"
}
	`,
	}

	// invalid value for `c_hash` inside the `id_token` field
	testItemInvalidCHash := testTableItem{
		label:                "query invalid c_hash",
		endpoint:             `/api/redirect/query/ok`,
		httpStatusExpected:   http.StatusBadRequest,
		responseBodyExpected: `{"error":"c_hash invalid"}`,
		requestBody: `
{
    "code": "a052c795-742d-415a-843f-8a4939d740d1",
    "scope": "openid accounts",
    "id_token": "eyJhbGciOiJQUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJodHRwczovL2p3dC1pZHAuZXhhbXBsZS5jb20iLCJzdWIiOiJtYWlsdG86bWlrZUBleGFtcGxlLmNvbSIsIm5iZiI6MTU0OTU1NjY0MiwiZXhwIjoxNTQ5NTYwMjQyLCJpYXQiOjE1NDk1NTY2NDIsImp0aSI6ImlkMTIzNDU2IiwidHlwIjoiaHR0cHM6Ly9leGFtcGxlLmNvbS9yZWdpc3RlciIsICJjX2hhc2giOiAiaW52YWxpZC1jX2hhc2gifQ.CG8bb_wT7EetLsE3SB8W30K3z_be14ZNXsjWQiklXkqImE-aWFvCqruwh-3aAG5xvaQ_u6T5mj7jaK-ZX93v591FsMPmX1MWyUYNfJp5MsPsWUfzZX69Us5UAqOgZ2zxu662prcE8fVqsL-GB-boVR_0e1SUj4NjKhiEHCNVYe-SGclRSZPvjRf0ymQBacmFz84kLqFVZYTXJFkufXd09FUopnNVK-aK2aCc39TaxzxFwwLaAW_iOtJnzHUtnNdF1OUW5MLTeJYd7hPg0Oq5hPUtz2h2XLVl76ERdJYNWNa1yws4gaWE9PaDgNu-mYbfEZVIHnB28XkB7d6BaCW-GQ",
    "state": "5a6b0d7832a9fb4f80f1170a"
}
	`,
	}

	ttData := []testTableItem{
		testItemOK,
		testItemInvalidCode,
		testItemInvalidCHash,
	}

	for _, ttItem := range ttData {
		// read the file we expect to be served.
		code, body, headers := request(
			http.MethodPost,
			ttItem.endpoint,
			strings.NewReader(ttItem.requestBody),
			server,
		)

		// do assertions.
		require.Equal(ttItem.httpStatusExpected, code, ttItem.label)
		require.Len(headers, 2, ttItem.label)
		require.Equal("application/json; charset=UTF-8", headers["Content-Type"][0], ttItem.label)
		require.NotNil(body, ttItem.label)

		bodyActual := body.String()
		require.JSONEq(ttItem.responseBodyExpected, bodyActual, ttItem.label)
	}
}

func TestRedirectHandlersError(t *testing.T) {
	require := require.New(t)

	server := NewServer(nullLogger(), conditionalityCheckerMock{}, mockVersionChecker())
	defer func() {
		require.NoError(server.Shutdown(context.TODO()))
	}()
	require.NotNil(server)

	bodyExpected := `
{
    "error_description": "JWT invalid. Expiration time incorrect.",
    "error": "invalid_request",
    "state": "5a6b0d7832a9fb4f80f1170a"
}
	`
	// read the file we expect to be served.
	code, body, headers := request(
		http.MethodPost,
		`/api/redirect/error`,
		strings.NewReader(bodyExpected),
		server,
	)

	// do assertions.
	require.Equal(http.StatusOK, code)
	require.Len(headers, 2)
	require.Equal("application/json; charset=UTF-8", headers["Content-Type"][0])
	require.NotNil(body)

	bodyActual := body.String()
	require.JSONEq(bodyExpected, bodyActual)
}

func TestCalculateCHash(t *testing.T) {
	require := require.New(t)

	tt := []struct {
		label         string
		code          string
		alg           string
		expectedHash  string
		expectedError error
	}{
		{
			label:        "ES256 empty code",
			code:         "",
			alg:          "ES256",
			expectedHash: "47DEQpj8HBSa-_TImW-5JA",
		},
		{
			label:        "ES256 code valid",
			code:         "80bf17a3-e617-4983-9d62-b50bd8e6fce4",
			alg:          "ES256",
			expectedHash: "EE_Bf-grXWv5GGhs5FZ0ug",
		},
		{
			label:        "PS256 code valid",
			code:         "80bf17a3-e617-4983-9d62-b50bd8e6fce4",
			alg:          "PS256",
			expectedHash: "EE_Bf-grXWv5GGhs5FZ0ug",
		},
		{
			label:         "algorithm not supported",
			code:          "80bf17a3-e617-4983-9d62-b50bd8e6fce4",
			alg:           "bad-algorithm",
			expectedHash:  "",
			expectedError: fmt.Errorf("bad-algorithm algorithm not supported"),
		},
	}

	for _, tti := range tt {
		cHash, err := calculateCHash(tti.alg, tti.code)
		require.Equal(tti.expectedHash, cHash, tti.label)

		if tti.expectedError != nil {
			require.Equal(tti.expectedError, err, tti.label)
		}

	}
}
