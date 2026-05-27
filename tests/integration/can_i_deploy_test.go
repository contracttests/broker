package integration_test

import (
	"context"
	"database/sql"
	"net/http"
)

const (
	canIDeployApiParticipant   = `{"name":"api"}`
	canIDeployFrontParticipant = `{"name":"front"}`
	canIDeployEnvironment      = `{"name":"production"}`
)

const apiV1Contract = `{
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

const apiV2Contract = `{
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
        "id":   { "type": "string" },
        "name": { "type": "string" }
      }
    }
  }
}`

const frontV1Contract = `{
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

const frontV2Contract = `{
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
        "id":   { "type": "string" },
        "name": { "type": "string" }
      }
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
		var row matrixRow
		s.Require().NoError(rows.Scan(&row.ParticipantID, &row.Version, &row.CounterpartParticipantID, &row.CounterpartVersion, &row.Deployable))
		out = append(out, row)
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

func (s *IntegrationSuite) seedApiAndFront() {
	s.seedParticipant(canIDeployApiParticipant)
	s.seedParticipant(canIDeployFrontParticipant)
	s.seedEnvironment(canIDeployEnvironment)
}

// Happy: a pure provider with no consumers in the environment is always
// deployable; the call records a vacuous-true row.
func (s *IntegrationSuite) TestCanIDeploy_VacuousTrue_PureProviderHasNoConsumers() {
	s.seedParticipant(canIDeployApiParticipant)
	s.seedEnvironment(canIDeployEnvironment)
	s.seedContract("api", "v1", apiV1Contract)

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

// Happy: front@v1 consumes api@v1, api@v1 deployed → compatible pair.
func (s *IntegrationSuite) TestCanIDeploy_CompatiblePair_DeployableTrue() {
	s.seedApiAndFront()
	s.seedContract("api", "v1", apiV1Contract)
	s.seedContract("front", "v1", frontV1Contract)
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

// Unhappy: front@v2 expects a field api@v1 does not provide → breaking change.
func (s *IntegrationSuite) TestCanIDeploy_IncompatiblePair_DeployableFalse() {
	s.seedApiAndFront()
	s.seedContract("api", "v1", apiV1Contract)
	s.seedContract("front", "v2", frontV2Contract)
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

// Unhappy: api is a known participant but has no deployment in env → strict-false.
func (s *IntegrationSuite) TestCanIDeploy_StrictFalse_CounterpartParticipantNeverDeployed() {
	s.seedApiAndFront()
	s.seedContract("api", "v1", apiV1Contract)
	s.seedContract("front", "v1", frontV1Contract)

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

// Workflow: front@v2 is initially blocked by api@v1; deploying api@v2
// unblocks it.
func (s *IntegrationSuite) TestCanIDeploy_CounterpartUpgradeUnblocksAsker() {
	s.seedApiAndFront()
	s.seedContract("api", "v1", apiV1Contract)
	s.seedContract("api", "v2", apiV2Contract)
	s.seedContract("front", "v2", frontV2Contract)
	s.seedDeployment("api", "v1", "production")

	status, body := s.get("/api/front/can-i-deploy?version=v2&environment=production")
	s.Equal(http.StatusOK, status)
	s.JSONEq(`{"success":true,"deployable":false}`, body)

	s.seedDeployment("api", "v2", "production")

	status, body = s.get("/api/front/can-i-deploy?version=v2&environment=production")
	s.Equal(http.StatusOK, status)
	s.JSONEq(`{"success":true,"deployable":true}`, body)

	rows := s.loadAllMatrixRows()
	s.Require().Len(rows, 2)
	s.Equal("v1", rows[0].CounterpartVersion.String)
	s.False(rows[0].Deployable)
	s.Equal("v2", rows[1].CounterpartVersion.String)
	s.True(rows[1].Deployable)
}

func (s *IntegrationSuite) TestCanIDeploy_MissingVersion_Returns400() {
	s.seedParticipant(canIDeployFrontParticipant)

	status, body := s.get("/api/front/can-i-deploy?environment=production")
	s.Equal(http.StatusBadRequest, status)
	s.JSONEq(`{"success":false,"message":"can-i-deploy invalid input"}`, body)
	s.Equal(0, s.countRows("compatibility_matrix"))
}

func (s *IntegrationSuite) TestCanIDeploy_MissingEnvironment_Returns400() {
	s.seedParticipant(canIDeployFrontParticipant)

	status, body := s.get("/api/front/can-i-deploy?version=v1")
	s.Equal(http.StatusBadRequest, status)
	s.JSONEq(`{"success":false,"message":"can-i-deploy invalid input"}`, body)
	s.Equal(0, s.countRows("compatibility_matrix"))
}

func (s *IntegrationSuite) TestCanIDeploy_UnknownParticipant_Returns404() {
	status, body := s.get("/api/unknown/can-i-deploy?version=v1&environment=production")
	s.Equal(http.StatusNotFound, status)
	s.JSONEq(`{"success":false,"message":"participant not found"}`, body)
	s.Equal(0, s.countRows("compatibility_matrix"))
}

func (s *IntegrationSuite) TestCanIDeploy_UnpublishedVersion_Returns422() {
	s.seedParticipant(canIDeployFrontParticipant)
	s.seedEnvironment(canIDeployEnvironment)

	status, body := s.get("/api/front/can-i-deploy?version=v99&environment=production")
	s.Equal(http.StatusUnprocessableEntity, status)
	s.JSONEq(`{"success":false,"message":"version not published"}`, body)
	s.Equal(0, s.countRows("compatibility_matrix"))
}

func (s *IntegrationSuite) TestCanIDeploy_UnknownEnvironment_Returns422() {
	s.seedParticipant(canIDeployFrontParticipant)
	s.seedContract("front", "v1", frontV1Contract)

	status, body := s.get("/api/front/can-i-deploy?version=v1&environment=ghost")
	s.Equal(http.StatusUnprocessableEntity, status)
	s.JSONEq(`{"success":false,"message":"environment not found"}`, body)
	s.Equal(0, s.countRows("compatibility_matrix"))
}
