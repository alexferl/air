package handlers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestAssetBadRequests(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/:id")
	c.SetParamNames("id")
	c.SetParamValues("123")
	h := &Handler{}

	type params struct {
		key   string
		value string
	}

	requests := []struct {
		params []params
		code   int
	}{
		{[]params{{"size", "120x"}}, http.StatusBadRequest},
		{[]params{{"size", "-1"}}, http.StatusBadRequest},
		{[]params{{"size", "9001"}}, http.StatusBadRequest},
		{[]params{{"width", "a"}}, http.StatusBadRequest},
		{[]params{{"width", "-1"}}, http.StatusBadRequest},
		{[]params{{"width", "9001"}}, http.StatusBadRequest},
		{[]params{{"height", "-1"}}, http.StatusBadRequest},
		{[]params{{"height", "a"}}, http.StatusBadRequest},
		{[]params{{"height", "9001"}}, http.StatusBadRequest},
		{[]params{{"quality", "-1"}}, http.StatusBadRequest},
		{[]params{{"quality", "a"}}, http.StatusBadRequest},
		{[]params{{"quality", "9001"}}, http.StatusBadRequest},
		{[]params{{"format", "1"}}, http.StatusBadRequest},
		{[]params{{"format", "a"}}, http.StatusBadRequest},
		{[]params{{"width", "100"}, {"height", "-1"}}, http.StatusBadRequest},
		{[]params{{"width", "-1"}, {"height", "100"}}, http.StatusBadRequest},
		{[]params{{"width", "100"}, {"height", "a"}}, http.StatusBadRequest},
		{[]params{{"width", "a"}, {"height", "100"}}, http.StatusBadRequest},
		{[]params{{"width", "100"}, {"height", "100"}, {"quality", "-1"}}, http.StatusBadRequest},
		{[]params{{"width", "100"}, {"height", "100"}, {"format", "-1"}}, http.StatusBadRequest},
		{[]params{{"width", "100"}, {"height", "100"}, {"quality", "80"}, {"format", "a"}}, http.StatusBadRequest},
	}

	for _, request := range requests {
		req.URL.RawQuery = ""
		for _, param := range request.params {
			q := req.URL.Query()
			q.Add(param.key, param.value)
			req.URL.RawQuery = q.Encode()
		}
		fmt.Printf("Query params: %s\n", req.URL.RawQuery)
		if assert.NoError(t, h.Asset(c)) {
			assert.Equal(t, request.code, rec.Code)
		}
	}
}
