package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Attendance struct {
	gorm.Model
	EmployeeId uuid.UUID
	PeriodId   uuid.UUID
	Date       time.Time `gorm:"index"`
}
