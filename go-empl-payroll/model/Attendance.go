package model

import (
	"time"

	"gorm.io/gorm"
)

type Attendance struct {
	gorm.Model
	EmployeeId string
	PeriodId   string
	Date       time.Time `gorm:"index"`
}
