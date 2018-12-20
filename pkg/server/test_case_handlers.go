package server

import (
	"bitbucket.org/openbankingteam/conformance-suite/pkg/web"
	"github.com/labstack/echo"
	"net/http"
)

type testCaseHandlers struct {
	webJourney web.Journey
}

func newTestCaseHandlers(webJourney web.Journey) testCaseHandlers {
	return testCaseHandlers{webJourney}
}

func (d testCaseHandlers) testCasesHandler(c echo.Context) error {
	testCases := d.webJourney.TestCases()
	return c.JSON(http.StatusOK, testCases)
}
