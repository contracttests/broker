package model

import "time"

type CompatibilityMatrixRow struct {
	ID                       int64
	ParticipantID            int64
	Version                  string
	CounterpartParticipantID *int64
	CounterpartVersion       *string
	Deployable               bool
	CreatedAt                time.Time
}

func NewVacuousTrueRow(participantID int64, version string) *CompatibilityMatrixRow {
	return &CompatibilityMatrixRow{
		ParticipantID: participantID,
		Version:       version,
		Deployable:    true,
	}
}

func NewStrictFalseRow(participantID int64, version string, counterpartID int64) *CompatibilityMatrixRow {
	return &CompatibilityMatrixRow{
		ParticipantID:            participantID,
		Version:                  version,
		CounterpartParticipantID: &counterpartID,
		Deployable:               false,
	}
}

func NewPairCheckedRow(participantID int64, version string, counterpartID int64, counterpartVersion string, deployable bool) *CompatibilityMatrixRow {
	return &CompatibilityMatrixRow{
		ParticipantID:            participantID,
		Version:                  version,
		CounterpartParticipantID: &counterpartID,
		CounterpartVersion:       &counterpartVersion,
		Deployable:               deployable,
	}
}
