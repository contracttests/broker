package integration_test

import (
	"net/http"
)

const (
	petsParticipantBody = `{"name":"pets-service"}`

	uploadContractBody = `{
  "participant": "pets-service",
  "version": "1",
  "contract": {
    "provides": {
      "rest": {
        "/pets": {
          "get": {
            "responses": {
              "200": "Pet"
            }
          }
        }
      }
    },
    "schemas": {
      "Pet": {
        "type": "object",
        "properties": {
          "id": { "type": "string" },
          "name": { "type": "string" }
        }
      }
    }
  }
}`
)

func (s *IntegrationSuite) TestHappyPath_CreateParticipantThenUploadContract() {
	status, body := s.post("/api/participants", petsParticipantBody)
	s.Equal(http.StatusOK, status)
	s.JSONEq(`{"success":true,"message":"participant created"}`, body)

	status, body = s.post("/api/contracts", uploadContractBody)
	s.Equal(http.StatusOK, status)
	s.JSONEq(`{"success":true,"message":"contract upload successful"}`, body)

	s.Equal(1, s.countRows("participants"))
	s.Equal(1, s.countRows("contracts"))
	s.Equal(1, s.countRows("resources"))
	s.GreaterOrEqual(s.countRows("properties"), 1)
}

func (s *IntegrationSuite) TestUnhappyPath_DuplicateParticipantName() {
	status, _ := s.post("/api/participants", petsParticipantBody)
	s.Equal(http.StatusOK, status)

	status, body := s.post("/api/participants", petsParticipantBody)
	s.Equal(http.StatusBadRequest, status)
	s.JSONEq(`{"success":false,"message":"participant already exists"}`, body)

	s.Equal(1, s.countRows("participants"))
}
