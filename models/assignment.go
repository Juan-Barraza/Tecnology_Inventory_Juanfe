package models

import "time"

type AssignmentStatus string

const (
	AssignmentStatusActive     AssignmentStatus = "active"
	AssignmentStatusReleased   AssignmentStatus = "released"
	AssignmentStatusWrittenOff AssignmentStatus = "written_off"
)

type Assignment struct {
	ID                  string
	AssetID             string
	ResponsibleName     *string
	ResponsiblePosition *string
	AssignedAt          time.Time
	DeactivatedAt       *time.Time
	DeactivationReason  *string
	Status              AssignmentStatus
	CreatedBy           string
	CreatedAt           time.Time
}

type AssignmentDetail struct {
	Assignment
	AssetCode        string
	AssetDescription string
	CreatedByName    string
}
