package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Reimbursement struct {
	gorm.Model
	EmployeeId  string
	Amount      float64
	Date        time.Time
	Description string
	PeriodId    uuid.UUID
}
