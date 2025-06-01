package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AttendancePeriod struct {
	gorm.Model
	StartDate   time.Time
	EndDate     time.Time
	IsProcessed bool
	PeriodId    uuid.UUID
}
