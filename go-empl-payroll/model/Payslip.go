package model

import (
	"time"

	"gorm.io/gorm"
)

type Payslip struct {
	gorm.Model
	PeriodID           string
	EmployeeID         string
	AttendanceDays     int
	OvertimeHours      float64
	Reimbursement      float64
	BaseSalary         float64
	TotalOvertime      int
	TotalReimbursement float64
	TotalPay           float64
	GeneratedAt        time.Time
}
