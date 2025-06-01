package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Overtime struct {
	gorm.Model
	EmployeeId uuid.UUID
	Date       time.Time `gorm:"index"`
	Hours      int
	PeriodId   uint
}
