package middleware

import (
	"encoding/json"
	"strings"

	"github.com/contracttesting/broker/internal/model"
	"github.com/gofiber/fiber/v3"
)

type Extractor func(ctx fiber.Ctx) string

func FromPath(name string) Extractor {
	return func(ctx fiber.Ctx) string {
		return strings.TrimSpace(ctx.Params(name))
	}
}

func FromQuery(name string) Extractor {
	return func(ctx fiber.Ctx) string {
		return strings.TrimSpace(ctx.Query(name))
	}
}

func FromBody(field string) Extractor {
	return func(ctx fiber.Ctx) string {
		body := map[string]any{}
		if err := json.Unmarshal(ctx.Body(), &body); err != nil {
			return ""
		}
		value, _ := body[field].(string)
		return strings.TrimSpace(value)
	}
}

const (
	ParticipantRequired = "participant is required"
	ParticipantNotFound = "participant not found"
	VersionRequired     = "version is required"
	VersionNotFound     = "version not found"
	EnvironmentRequired = "environment is required"
	EnvironmentNotFound = "environment not found"
)

type errorResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func reject(ctx fiber.Ctx, message string) error {
	return ctx.Status(fiber.StatusBadRequest).JSON(errorResponse{
		Success: false,
		Message: message,
	})
}

type localsKey string

const (
	participantKey localsKey = "middleware.participant"
	versionKey     localsKey = "middleware.version"
	environmentKey localsKey = "middleware.environment"
)

func ParticipantFrom(ctx fiber.Ctx) *model.Participant {
	participant, _ := ctx.Locals(participantKey).(*model.Participant)
	return participant
}

func EnvironmentFrom(ctx fiber.Ctx) *model.Environment {
	environment, _ := ctx.Locals(environmentKey).(*model.Environment)
	return environment
}

func VersionFrom(ctx fiber.Ctx) string {
	version, _ := ctx.Locals(versionKey).(string)
	return version
}
