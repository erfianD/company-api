package model

import (
	"time"

	"gorm.io/gorm"
)

type Reimbursement struct {
	gorm.Model
	EmployeeId  string
	Amount      float64
	Date        time.Time
	Description string
	PeriodId    string
}
