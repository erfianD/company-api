package model

import (
	"time"

	"gorm.io/gorm"
)

type AttendancePeriod struct {
	gorm.Model
	StartDate   time.Time
	EndDate     time.Time
	IsProcessed bool
	PeriodId    string
}
