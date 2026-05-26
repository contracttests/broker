package integration_test

import (
	"net/http"
)

const petsParticipantBody = `{"name":"pets-service"}`

func (s *IntegrationSuite) TestHappyPath_CreateParticipant() {
	status, body := s.post("/api/participants", petsParticipantBody)
	s.Equal(http.StatusOK, status)
	s.JSONEq(`{"success":true,"message":"participant created"}`, body)

	s.Equal(1, s.countRows("participants"))
}

func (s *IntegrationSuite) TestUnhappyPath_DuplicateParticipantName() {
	status, _ := s.post("/api/participants", petsParticipantBody)
	s.Equal(http.StatusOK, status)

	status, body := s.post("/api/participants", petsParticipantBody)
	s.Equal(http.StatusBadRequest, status)
	s.JSONEq(`{"success":false,"message":"participant already exists"}`, body)

	s.Equal(1, s.countRows("participants"))
}
