package model

import (
	"time"

	"gorm.io/gorm"
)

type Overtime struct {
	gorm.Model
	EmployeeId string
	Date       time.Time `gorm:"index"`
	Hours      int
	PeriodId   uint
}
