package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Payslip struct {
	gorm.Model
	PeriodID           uuid.UUID
	EmployeeID         uuid.UUID
	AttendanceDays     int
	OvertimeHours      float64
	Reimbursement      float64
	BaseSalary         float64
	TotalOvertime      int
	TotalReimbursement float64
	TotalPay           float64
	GeneratedAt        time.Time
}
