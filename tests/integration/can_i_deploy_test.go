package integration_test

import (
	"context"
	"encoding/json"
	"net/http"
)

// api@v1 provides Thing{id}.
const apiV1ProviderContract = `
{
  "provides": {
    "rest": {
      "/things": {
        "get": {
          "responses": {
            "200": "Thing"
          }
        }
      }
    }
  },
  "schemas": {
    "Thing": {
      "type": "object",
      "properties": {
        "id": { "type": "string" }
      }
    }
  }
}`

// front@v1 consumes only the "id" field that api@v1 provides.
const frontV1ConsumerContract = `
{
  "consumes": {
    "api": {
      "rest": {
        "/things": {
          "get": {
            "responses": {
              "200": "Thing"
            }
          }
        }
      }
    }
  },
  "schemas": {
    "Thing": {
      "type": "object",
      "properties": {
        "id": { "type": "string" }
      }
    }
  }
}`

// front@v2 expects "id" as an integer (api@v1 provides a string) and adds "name"
// (api@v1 does not provide it at all): two breaking changes.
const frontV2ConsumerContract = `
{
  "consumes": {
    "api": {
      "rest": {
        "/things": {
          "get": {
            "responses": {
              "200": "Thing"
            }
          }
        }
      }
    }
  },
  "schemas": {
    "Thing": {
      "type": "object",
      "properties": {
        "id": { "type": "integer" },
        "name": { "type": "string" }
      }
    }
  }
}`

func (s *IntegrationSuite) TestCanIDeploy_HappyPath() {
	// api@v1 provides Thing{id}; publish it and deploy it to production so the
	// compatibility check can resolve it as the provider in that environment.
	status, _ := s.post("/api/participants", `{"name":"api"}`)
	s.Require().Equal(http.StatusOK, status)

	status, _ = s.post("/api/contracts", `{"name":"api","version":"v1","contract":`+apiV1ProviderContract+`}`)
	s.Require().Equal(http.StatusOK, status)

	status, _ = s.post("/api/environments", `{"name":"production"}`)
	s.Require().Equal(http.StatusOK, status)

	status, _ = s.post("/api/deployments", `{"name":"api","version":"v1","environment":"production"}`)
	s.Require().Equal(http.StatusOK, status)

	status, _ = s.post("/api/participants", `{"name":"front"}`)
	s.Require().Equal(http.StatusOK, status)

	// front@v1 consumes only "id", which api@v1 provides: it is deployable.
	status, _ = s.post("/api/contracts", `{"name":"front","version":"v1","contract":`+frontV1ConsumerContract+`}`)
	s.Require().Equal(http.StatusOK, status)

	status, body := s.post("/api/can-i-deploy", `{"name":"front","version":"v1","environment":"production"}`)
	s.Equal(http.StatusOK, status)
	s.JSONEq(`{"success":true,"deployable":true}`, body)

	// A compatible decision is persisted as a deployable row.
	s.Equal(1, s.countRows("compatibility_matrix"))
	var v1Deployable bool
	s.Require().NoError(s.Pool.QueryRow(context.Background(),
		`SELECT deployable FROM compatibility_matrix WHERE version = 'v1'`).Scan(&v1Deployable))
	s.True(v1Deployable)

	status, _ = s.post("/api/deployments", `{"name":"front","version":"v1","environment":"production"}`)
	s.Require().Equal(http.StatusOK, status)

	// front@v2 is incompatible with api@v1 on two counts, so it is not deployable
	// and the response carries the breaking changes.
	status, _ = s.post("/api/contracts", `{"name":"front","version":"v2","contract":`+frontV2ConsumerContract+`}`)
	s.Require().Equal(http.StatusOK, status)

	status, body = s.post("/api/can-i-deploy", `{"name":"front","version":"v2","environment":"production"}`)
	s.Equal(http.StatusOK, status)

	type brokenResource struct {
		Direction  string `json:"direction"`
		Kind       string `json:"kind"`
		Provider   string `json:"provider"`
		Endpoint   string `json:"endpoint"`
		Method     string `json:"method"`
		StatusCode string `json:"status_code"`
	}

	type breakItem struct {
		LeftResource  brokenResource `json:"left_resource"`
		RightResource brokenResource `json:"right_resource"`
		Reason        string         `json:"reason"`
		Property      string         `json:"property"`
		HumanReadable string         `json:"human_readable"`
	}

	var got struct {
		Success    bool                   `json:"success"`
		Deployable bool                   `json:"deployable"`
		Breaks     map[string][]breakItem `json:"breaks"`
	}

	s.Require().NoError(json.Unmarshal([]byte(body), &got))
	s.True(got.Success)
	s.False(got.Deployable)

	consumer := brokenResource{
		Direction:  "consumes",
		Kind:       "rest_response",
		Provider:   "api",
		Endpoint:   "/things",
		Method:     "get",
		StatusCode: "200",
	}

	provider := brokenResource{
		Direction:  "provides",
		Kind:       "rest_response",
		Provider:   "",
		Endpoint:   "/things",
		Method:     "get",
		StatusCode: "200",
	}

	// The break order is not deterministic (the checker ranges over a property map).
	s.ElementsMatch([]breakItem{
		{
			LeftResource:  consumer,
			RightResource: provider,
			Reason:        "type_mismatch",
			Property:      "root.id",
			HumanReadable: "Property root.id type mismatch, provider api expects string but consumer front expects integer",
		},
		{
			LeftResource:  consumer,
			RightResource: provider,
			Reason:        "missing_in_provider",
			Property:      "root.name",
			HumanReadable: "Property root.name is missing in provider api",
		},
	}, got.Breaks["front"])

	// The incompatible decision is persisted too, as a non-deployable row.
	s.Equal(2, s.countRows("compatibility_matrix"))
	var v2Deployable bool
	s.Require().NoError(s.Pool.QueryRow(context.Background(),
		`SELECT deployable FROM compatibility_matrix WHERE version = 'v2'`).Scan(&v2Deployable))
	s.False(v2Deployable)
}
