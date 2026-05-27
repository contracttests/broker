package integration_test

import (
	"context"
	"database/sql"
	"net/http"
	"time"
)

const (
	apiParticipantForCID    = `{"name":"api"}`
	frontParticipantBody    = `{"name":"front"}`
	billingParticipantBody  = `{"name":"billing"}`
	dbParticipantBody       = `{"name":"db"}`
	productionEnvForCID     = `{"name":"production"}`
)

const apiV1ProvidesThings = `{
  "provides": {
    "rest": {
      "/things": {
        "get": { "responses": { "200": "ThingV1" } }
      }
    }
  },
  "schemas": {
    "ThingV1": {
      "type": "object",
      "properties": { "id": { "type": "string" } }
    }
  }
}`

const apiV2ProvidesThings = `{
  "provides": {
    "rest": {
      "/things": {
        "get": { "responses": { "200": "ThingV2" } }
      }
    }
  },
  "schemas": {
    "ThingV2": {
      "type": "object",
      "properties": {
        "id": { "type": "string" },
        "name": { "type": "string" }
      }
    }
  }
}`

const frontV1ConsumesThings = `{
  "consumes": {
    "api": {
      "rest": {
        "/things": {
          "get": { "responses": { "200": "ThingV1" } }
        }
      }
    }
  },
  "schemas": {
    "ThingV1": {
      "type": "object",
      "properties": { "id": { "type": "string" } }
    }
  }
}`

const frontV2ConsumesThings = `{
  "consumes": {
    "api": {
      "rest": {
        "/things": {
          "get": { "responses": { "200": "ThingV2" } }
        }
      }
    }
  },
  "schemas": {
    "ThingV2": {
      "type": "object",
      "properties": {
        "id": { "type": "string" },
        "name": { "type": "string" }
      }
    }
  }
}`

const frontV1ConsumesGhost = `{
  "consumes": {
    "ghost": {
      "rest": {
        "/x": {
          "get": { "responses": { "200": "X" } }
        }
      }
    }
  },
  "schemas": {
    "X": {
      "type": "object",
      "properties": { "id": { "type": "string" } }
    }
  }
}`

const frontV2ProvidesOnly = `{
  "provides": {
    "rest": {
      "/internal": {
        "get": { "responses": { "200": "Internal" } }
      }
    }
  },
  "schemas": {
    "Internal": {
      "type": "object",
      "properties": { "id": { "type": "string" } }
    }
  }
}`

const frontV2ProvidesOnlyAlt = `{
  "provides": {
    "rest": {
      "/internal-v2": {
        "get": { "responses": { "200": "Internal" } }
      }
    }
  },
  "schemas": {
    "Internal": {
      "type": "object",
      "properties": { "id": { "type": "string" } }
    }
  }
}`

const apiV1ProvidesThingsAndWidgets = `{
  "provides": {
    "rest": {
      "/things": {
        "get": { "responses": { "200": "ThingV1" } }
      },
      "/widgets": {
        "get": { "responses": { "200": "WidgetV1" } }
      }
    }
  },
  "schemas": {
    "ThingV1": {
      "type": "object",
      "properties": { "id": { "type": "string" } }
    },
    "WidgetV1": {
      "type": "object",
      "properties": { "id": { "type": "string" } }
    }
  }
}`

const frontConsumesThingsAndWidgets = `{
  "consumes": {
    "api": {
      "rest": {
        "/things": {
          "get": { "responses": { "200": "ThingV1" } }
        },
        "/widgets": {
          "get": { "responses": { "200": "WidgetV1" } }
        }
      }
    }
  },
  "schemas": {
    "ThingV1": {
      "type": "object",
      "properties": { "id": { "type": "string" } }
    },
    "WidgetV1": {
      "type": "object",
      "properties": { "id": { "type": "string" } }
    }
  }
}`

const billingV1ProvidesInvoices = `{
  "provides": {
    "rest": {
      "/invoices": {
        "get": { "responses": { "200": "InvoiceV1" } }
      }
    }
  },
  "schemas": {
    "InvoiceV1": {
      "type": "object",
      "properties": { "id": { "type": "string" } }
    }
  }
}`

const dbV1ProvidesRecords = `{
  "provides": {
    "rest": {
      "/records": {
        "get": { "responses": { "200": "RecordV1" } }
      }
    }
  },
  "schemas": {
    "RecordV1": {
      "type": "object",
      "properties": { "id": { "type": "string" } }
    }
  }
}`

const frontConsumesThreeCounterparts = `{
  "consumes": {
    "api": {
      "rest": {
        "/things": {
          "get": { "responses": { "200": "ThingV1" } }
        }
      }
    },
    "billing": {
      "rest": {
        "/invoices": {
          "get": { "responses": { "200": "InvoiceV1" } }
        }
      }
    },
    "db": {
      "rest": {
        "/records": {
          "get": { "responses": { "200": "RecordV1" } }
        }
      }
    }
  },
  "schemas": {
    "ThingV1": {
      "type": "object",
      "properties": { "id": { "type": "string" } }
    },
    "InvoiceV1": {
      "type": "object",
      "properties": { "id": { "type": "string" } }
    },
    "RecordV1": {
      "type": "object",
      "properties": { "id": { "type": "string" } }
    }
  }
}`

const billingV1ProvidesInvoicesBreaking = `{
  "provides": {
    "rest": {
      "/invoices": {
        "get": { "responses": { "200": "InvoiceV1" } }
      }
    }
  },
  "schemas": {
    "InvoiceV1": {
      "type": "object",
      "properties": { "id": { "type": "integer" } }
    }
  }
}`

const frontConsumesApiAndBilling = `{
  "consumes": {
    "api": {
      "rest": {
        "/things": {
          "get": { "responses": { "200": "ThingV1" } }
        }
      }
    },
    "billing": {
      "rest": {
        "/invoices": {
          "get": { "responses": { "200": "InvoiceV1" } }
        }
      }
    }
  },
  "schemas": {
    "ThingV1": {
      "type": "object",
      "properties": { "id": { "type": "string" } }
    },
    "InvoiceV1": {
      "type": "object",
      "properties": { "id": { "type": "string" } }
    }
  }
}`

type matrixRow struct {
	ParticipantID            int64
	Version                  string
	CounterpartParticipantID sql.NullInt64
	CounterpartVersion       sql.NullString
	Deployable               bool
}

func (s *IntegrationSuite) loadAllMatrixRows() []matrixRow {
	rows, err := s.Pool.Query(context.Background(),
		`SELECT participant_id, version, counterpart_participant_id, counterpart_version, deployable
		 FROM compatibility_matrix
		 ORDER BY id`,
	)
	s.Require().NoError(err)
	defer rows.Close()

	var out []matrixRow
	for rows.Next() {
		var r matrixRow
		s.Require().NoError(rows.Scan(&r.ParticipantID, &r.Version, &r.CounterpartParticipantID, &r.CounterpartVersion, &r.Deployable))
		out = append(out, r)
	}
	return out
}

func (s *IntegrationSuite) seedParticipant(body string) {
	status, _ := s.post("/api/participants", body)
	s.Require().Equal(http.StatusOK, status)
}

func (s *IntegrationSuite) seedEnvironment(body string) {
	status, _ := s.post("/api/environments", body)
	s.Require().Equal(http.StatusOK, status)
}

func (s *IntegrationSuite) seedContract(participant, version, body string) {
	status, _ := s.post("/api/"+participant+"/contracts/"+version, body)
	s.Require().Equal(http.StatusOK, status)
}

func (s *IntegrationSuite) seedDeployment(participant, version, environment string) {
	status, _ := s.post(
		"/api/"+participant+"/deployments",
		`{"version":"`+version+`","environment":"`+environment+`"}`,
	)
	s.Require().Equal(http.StatusOK, status)
}

// 10.2
func (s *IntegrationSuite) TestCanIDeploy_VacuousTrueForPureProviderWithNoConsumers() {
	s.seedParticipant(apiParticipantForCID)
	s.seedEnvironment(productionEnvForCID)
	s.seedContract("api", "v1", apiV1ProvidesThings)

	status, body := s.get("/api/api/can-i-deploy?version=v1&environment=production")
	s.Equal(http.StatusOK, status)
	s.JSONEq(`{"success":true,"deployable":true}`, body)

	rows := s.loadAllMatrixRows()
	s.Require().Len(rows, 1)
	s.Equal(s.lookupParticipantID("api"), rows[0].ParticipantID)
	s.Equal("v1", rows[0].Version)
	s.False(rows[0].CounterpartParticipantID.Valid)
	s.False(rows[0].CounterpartVersion.Valid)
	s.True(rows[0].Deployable)
}

// 10.3
func (s *IntegrationSuite) TestCanIDeploy_CompatiblePairReturnsTrue() {
	s.seedParticipant(apiParticipantForCID)
	s.seedParticipant(frontParticipantBody)
	s.seedEnvironment(productionEnvForCID)
	s.seedContract("api", "v1", apiV1ProvidesThings)
	s.seedContract("front", "v1", frontV1ConsumesThings)
	s.seedDeployment("api", "v1", "production")

	status, body := s.get("/api/front/can-i-deploy?version=v1&environment=production")
	s.Equal(http.StatusOK, status)
	s.JSONEq(`{"success":true,"deployable":true}`, body)

	rows := s.loadAllMatrixRows()
	s.Require().Len(rows, 1)
	s.Equal(s.lookupParticipantID("front"), rows[0].ParticipantID)
	s.Equal("v1", rows[0].Version)
	s.Equal(s.lookupParticipantID("api"), rows[0].CounterpartParticipantID.Int64)
	s.Equal("v1", rows[0].CounterpartVersion.String)
	s.True(rows[0].Deployable)
}

// 10.4
func (s *IntegrationSuite) TestCanIDeploy_IncompatiblePairReturnsFalse() {
	s.seedParticipant(apiParticipantForCID)
	s.seedParticipant(frontParticipantBody)
	s.seedEnvironment(productionEnvForCID)
	s.seedContract("api", "v1", apiV1ProvidesThings)
	s.seedContract("front", "v2", frontV2ConsumesThings)
	s.seedDeployment("api", "v1", "production")

	status, body := s.get("/api/front/can-i-deploy?version=v2&environment=production")
	s.Equal(http.StatusOK, status)
	s.JSONEq(`{"success":true,"deployable":false}`, body)

	rows := s.loadAllMatrixRows()
	s.Require().Len(rows, 1)
	s.Equal(s.lookupParticipantID("front"), rows[0].ParticipantID)
	s.Equal("v2", rows[0].Version)
	s.Equal(s.lookupParticipantID("api"), rows[0].CounterpartParticipantID.Int64)
	s.Equal("v1", rows[0].CounterpartVersion.String)
	s.False(rows[0].Deployable)
}

// 10.5
func (s *IntegrationSuite) TestCanIDeploy_StrictFalseWhenCounterpartParticipantPresentButNotDeployed() {
	s.seedParticipant(apiParticipantForCID)
	s.seedParticipant(frontParticipantBody)
	s.seedEnvironment(productionEnvForCID)
	s.seedContract("api", "v1", apiV1ProvidesThings)
	s.seedContract("front", "v1", frontV1ConsumesThings)

	status, body := s.get("/api/front/can-i-deploy?version=v1&environment=production")
	s.Equal(http.StatusOK, status)
	s.JSONEq(`{"success":true,"deployable":false}`, body)

	rows := s.loadAllMatrixRows()
	s.Require().Len(rows, 1)
	s.Equal(s.lookupParticipantID("front"), rows[0].ParticipantID)
	s.Equal("v1", rows[0].Version)
	s.Equal(s.lookupParticipantID("api"), rows[0].CounterpartParticipantID.Int64)
	s.False(rows[0].CounterpartVersion.Valid)
	s.False(rows[0].Deployable)
}

// 10.6
func (s *IntegrationSuite) TestCanIDeploy_FalseWithoutRowWhenCounterpartParticipantUnknown() {
	s.seedParticipant(frontParticipantBody)
	s.seedEnvironment(productionEnvForCID)
	s.seedContract("front", "v1", frontV1ConsumesGhost)

	status, body := s.get("/api/front/can-i-deploy?version=v1&environment=production")
	s.Equal(http.StatusOK, status)
	s.JSONEq(`{"success":true,"deployable":false}`, body)

	rows := s.loadAllMatrixRows()
	s.Require().Len(rows, 1)
	s.Equal(s.lookupParticipantID("front"), rows[0].ParticipantID)
	s.False(rows[0].CounterpartParticipantID.Valid)
	s.False(rows[0].CounterpartVersion.Valid)
	s.True(rows[0].Deployable)
}

// 10.7
func (s *IntegrationSuite) TestCanIDeploy_AskerExcludedFromCounterparts() {
	s.seedParticipant(frontParticipantBody)
	s.seedEnvironment(productionEnvForCID)
	s.seedContract("front", "v1", frontV2ProvidesOnly)
	s.seedDeployment("front", "v1", "production")
	s.seedContract("front", "v2", frontV2ProvidesOnlyAlt)

	status, _ := s.get("/api/front/can-i-deploy?version=v2&environment=production")
	s.Equal(http.StatusOK, status)

	frontID := s.lookupParticipantID("front")
	for _, r := range s.loadAllMatrixRows() {
		if r.CounterpartParticipantID.Valid && r.CounterpartParticipantID.Int64 == frontID {
			s.Failf("self-counterpart row present", "row: %+v", r)
		}
	}
}

// 10.8
func (s *IntegrationSuite) TestCanIDeploy_AskerNeedNotBeDeployed() {
	s.seedParticipant(apiParticipantForCID)
	s.seedParticipant(frontParticipantBody)
	s.seedEnvironment(productionEnvForCID)
	s.seedContract("api", "v1", apiV1ProvidesThings)
	s.seedContract("front", "v1", frontV1ConsumesThings)
	s.seedDeployment("api", "v1", "production")

	status, body := s.get("/api/front/can-i-deploy?version=v1&environment=production")
	s.Equal(http.StatusOK, status)
	s.JSONEq(`{"success":true,"deployable":true}`, body)
}

// 11.1
func (s *IntegrationSuite) TestCanIDeploy_OneRowPerCounterpartRegardlessOfSharedResources() {
	s.seedParticipant(apiParticipantForCID)
	s.seedParticipant(frontParticipantBody)
	s.seedEnvironment(productionEnvForCID)
	s.seedContract("api", "v1", apiV1ProvidesThingsAndWidgets)
	s.seedContract("front", "v2", frontConsumesThingsAndWidgets)
	s.seedDeployment("api", "v1", "production")

	status, _ := s.get("/api/front/can-i-deploy?version=v2&environment=production")
	s.Equal(http.StatusOK, status)

	rows := s.loadAllMatrixRows()
	s.Require().Len(rows, 1)
	s.Equal(s.lookupParticipantID("api"), rows[0].CounterpartParticipantID.Int64)
	s.Equal("v1", rows[0].CounterpartVersion.String)
	s.True(rows[0].Deployable)
}

// 11.2
func (s *IntegrationSuite) TestCanIDeploy_MultipleCounterpartsEachGetOwnRow() {
	s.seedParticipant(apiParticipantForCID)
	s.seedParticipant(billingParticipantBody)
	s.seedParticipant(dbParticipantBody)
	s.seedParticipant(frontParticipantBody)
	s.seedEnvironment(productionEnvForCID)
	s.seedContract("api", "v1", apiV1ProvidesThings)
	s.seedContract("billing", "v1", billingV1ProvidesInvoices)
	s.seedContract("db", "v1", dbV1ProvidesRecords)
	s.seedContract("front", "v1", frontConsumesThreeCounterparts)
	s.seedDeployment("api", "v1", "production")
	s.seedDeployment("billing", "v1", "production")
	s.seedDeployment("db", "v1", "production")

	status, body := s.get("/api/front/can-i-deploy?version=v1&environment=production")
	s.Equal(http.StatusOK, status)
	s.JSONEq(`{"success":true,"deployable":true}`, body)

	rows := s.loadAllMatrixRows()
	s.Require().Len(rows, 3)

	counterpartIDs := map[int64]bool{}
	for _, r := range rows {
		s.True(r.Deployable)
		s.True(r.CounterpartParticipantID.Valid)
		counterpartIDs[r.CounterpartParticipantID.Int64] = true
	}
	s.True(counterpartIDs[s.lookupParticipantID("api")])
	s.True(counterpartIDs[s.lookupParticipantID("billing")])
	s.True(counterpartIDs[s.lookupParticipantID("db")])
}

// 11.3
func (s *IntegrationSuite) TestCanIDeploy_AnyFalseMakesOverallFalse() {
	s.seedParticipant(apiParticipantForCID)
	s.seedParticipant(billingParticipantBody)
	s.seedParticipant(frontParticipantBody)
	s.seedEnvironment(productionEnvForCID)
	s.seedContract("api", "v1", apiV1ProvidesThings)
	s.seedContract("billing", "v1", billingV1ProvidesInvoicesBreaking)
	s.seedContract("front", "v1", frontConsumesApiAndBilling)
	s.seedDeployment("api", "v1", "production")
	s.seedDeployment("billing", "v1", "production")

	status, body := s.get("/api/front/can-i-deploy?version=v1&environment=production")
	s.Equal(http.StatusOK, status)
	s.JSONEq(`{"success":true,"deployable":false}`, body)

	rows := s.loadAllMatrixRows()
	s.Require().Len(rows, 2)

	apiID := s.lookupParticipantID("api")
	billingID := s.lookupParticipantID("billing")

	results := map[int64]bool{}
	for _, r := range rows {
		results[r.CounterpartParticipantID.Int64] = r.Deployable
	}
	s.True(results[apiID])
	s.False(results[billingID])
}

// 11.4
func (s *IntegrationSuite) TestCanIDeploy_RollbackResolvesLatestDeployedAt() {
	s.seedParticipant(apiParticipantForCID)
	s.seedParticipant(frontParticipantBody)
	s.seedEnvironment(productionEnvForCID)
	s.seedContract("api", "v1", apiV1ProvidesThings)
	s.seedContract("api", "v2", apiV2ProvidesThings)
	s.seedContract("front", "v1", frontV1ConsumesThings)

	apiID := s.lookupParticipantID("api")
	prodID := s.lookupEnvironmentID("production")
	s.insertDeploymentAt(apiID, "v1", prodID, "2026-05-01T00:00:00Z")
	s.insertDeploymentAt(apiID, "v2", prodID, "2026-05-10T00:00:00Z")
	s.insertDeploymentAt(apiID, "v1", prodID, "2026-05-15T00:00:00Z")

	status, body := s.get("/api/front/can-i-deploy?version=v1&environment=production")
	s.Equal(http.StatusOK, status)
	s.JSONEq(`{"success":true,"deployable":true}`, body)

	rows := s.loadAllMatrixRows()
	s.Require().Len(rows, 1)
	s.Equal(apiID, rows[0].CounterpartParticipantID.Int64)
	s.Equal("v1", rows[0].CounterpartVersion.String)
}

// 12.1
func (s *IntegrationSuite) TestCanIDeploy_DiscoveryWalkthrough() {
	s.seedParticipant(apiParticipantForCID)
	s.seedParticipant(frontParticipantBody)
	s.seedEnvironment(productionEnvForCID)
	s.seedContract("api", "v1", apiV1ProvidesThings)
	s.seedContract("front", "v1", frontV1ConsumesThings)
	s.seedDeployment("api", "v1", "production")
	s.seedDeployment("front", "v1", "production")

	s.seedContract("front", "v2", frontV2ConsumesThings)

	status, body := s.get("/api/front/can-i-deploy?version=v2&environment=production")
	s.Equal(http.StatusOK, status)
	s.JSONEq(`{"success":true,"deployable":false}`, body)

	frontID := s.lookupParticipantID("front")
	apiID := s.lookupParticipantID("api")

	rows := s.loadAllMatrixRows()
	s.Require().Len(rows, 1)
	s.Equal(frontID, rows[0].ParticipantID)
	s.Equal("v2", rows[0].Version)
	s.Equal(apiID, rows[0].CounterpartParticipantID.Int64)
	s.Equal("v1", rows[0].CounterpartVersion.String)
	s.False(rows[0].Deployable)

	s.seedContract("api", "v2", apiV2ProvidesThings)

	status, body = s.get("/api/api/can-i-deploy?version=v2&environment=production")
	s.Equal(http.StatusOK, status)
	s.JSONEq(`{"success":true,"deployable":true}`, body)

	rows = s.loadAllMatrixRows()
	s.Require().Len(rows, 2)
	s.Equal(apiID, rows[1].ParticipantID)
	s.Equal("v2", rows[1].Version)
	s.Equal(frontID, rows[1].CounterpartParticipantID.Int64)
	s.Equal("v1", rows[1].CounterpartVersion.String)
	s.True(rows[1].Deployable)

	s.seedDeployment("api", "v2", "production")

	status, body = s.get("/api/front/can-i-deploy?version=v2&environment=production")
	s.Equal(http.StatusOK, status)
	s.JSONEq(`{"success":true,"deployable":true}`, body)

	rows = s.loadAllMatrixRows()
	s.Require().Len(rows, 3)
	s.Equal(frontID, rows[2].ParticipantID)
	s.Equal("v2", rows[2].Version)
	s.Equal(apiID, rows[2].CounterpartParticipantID.Int64)
	s.Equal("v2", rows[2].CounterpartVersion.String)
	s.True(rows[2].Deployable)
}

// 13.1
func (s *IntegrationSuite) TestCanIDeploy_MissingVersionReturns400() {
	s.seedParticipant(frontParticipantBody)

	status, body := s.get("/api/front/can-i-deploy?environment=production")
	s.Equal(http.StatusBadRequest, status)
	s.JSONEq(`{"success":false,"message":"can-i-deploy invalid input"}`, body)
	s.Equal(0, s.countRows("compatibility_matrix"))
}

// 13.2
func (s *IntegrationSuite) TestCanIDeploy_MissingEnvironmentReturns400() {
	s.seedParticipant(frontParticipantBody)

	status, body := s.get("/api/front/can-i-deploy?version=v1")
	s.Equal(http.StatusBadRequest, status)
	s.JSONEq(`{"success":false,"message":"can-i-deploy invalid input"}`, body)
	s.Equal(0, s.countRows("compatibility_matrix"))
}

// 13.3
func (s *IntegrationSuite) TestCanIDeploy_UnknownParticipantReturns404() {
	status, body := s.get("/api/unknown/can-i-deploy?version=v1&environment=production")
	s.Equal(http.StatusNotFound, status)
	s.JSONEq(`{"success":false,"message":"participant not found"}`, body)
	s.Equal(0, s.countRows("compatibility_matrix"))
}

// 13.4
func (s *IntegrationSuite) TestCanIDeploy_UnpublishedVersionReturns422() {
	s.seedParticipant(frontParticipantBody)
	s.seedEnvironment(productionEnvForCID)

	status, body := s.get("/api/front/can-i-deploy?version=v99&environment=production")
	s.Equal(http.StatusUnprocessableEntity, status)
	s.JSONEq(`{"success":false,"message":"version not published"}`, body)
	s.Equal(0, s.countRows("compatibility_matrix"))
}

// 13.5
func (s *IntegrationSuite) TestCanIDeploy_UnknownEnvironmentReturns422() {
	s.seedParticipant(frontParticipantBody)
	s.seedContract("front", "v1", frontV2ProvidesOnly)

	status, body := s.get("/api/front/can-i-deploy?version=v1&environment=ghost")
	s.Equal(http.StatusUnprocessableEntity, status)
	s.JSONEq(`{"success":false,"message":"environment not found"}`, body)
	s.Equal(0, s.countRows("compatibility_matrix"))
}

// 13.6
func (s *IntegrationSuite) TestCanIDeploy_CheckConstraintRejectsNonsensicalRow() {
	s.seedParticipant(frontParticipantBody)
	frontID := s.lookupParticipantID("front")

	_, err := s.Pool.Exec(context.Background(),
		`INSERT INTO compatibility_matrix
		   (participant_id, version, counterpart_participant_id, counterpart_version, deployable)
		 VALUES
		   ($1, 'v1', NULL, 'something', true)`,
		frontID,
	)
	s.Require().Error(err)
}

// 14.1
func (s *IntegrationSuite) TestCanIDeploy_TwoIdenticalCallsProduceTwoRowSets() {
	s.seedParticipant(apiParticipantForCID)
	s.seedParticipant(frontParticipantBody)
	s.seedEnvironment(productionEnvForCID)
	s.seedContract("api", "v1", apiV1ProvidesThings)
	s.seedContract("front", "v1", frontV1ConsumesThings)
	s.seedDeployment("api", "v1", "production")

	status1, body1 := s.get("/api/front/can-i-deploy?version=v1&environment=production")
	time.Sleep(10 * time.Millisecond)
	status2, body2 := s.get("/api/front/can-i-deploy?version=v1&environment=production")
	s.Equal(http.StatusOK, status1)
	s.Equal(http.StatusOK, status2)
	s.Equal(body1, body2)

	rows, err := s.Pool.Query(context.Background(),
		`SELECT id, created_at FROM compatibility_matrix
		 WHERE participant_id = $1 AND version = 'v1'
		 ORDER BY id`,
		s.lookupParticipantID("front"),
	)
	s.Require().NoError(err)
	defer rows.Close()

	var ids []int64
	var createdAts []time.Time
	for rows.Next() {
		var id int64
		var at time.Time
		s.Require().NoError(rows.Scan(&id, &at))
		ids = append(ids, id)
		createdAts = append(createdAts, at)
	}

	s.Require().Len(ids, 2)
	s.NotEqual(ids[0], ids[1])
	s.NotEqual(createdAts[0], createdAts[1])
}
