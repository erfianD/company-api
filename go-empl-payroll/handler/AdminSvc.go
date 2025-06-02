package handler

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"go-empl-payroll/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AttendancePeriodInput struct {
	PeriodId  string `json:"period_id" binding:"required"`
	StartDate string `json:"start_date" binding:"required"`
	EndDate   string `json:"end_date" binding:"required"`
}

func CreateAttndancePeriod(c *gin.Context, db *gorm.DB) {
	var input AttendancePeriodInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	start, err1 := time.Parse("2006-01-02", input.StartDate)
	end, err2 := time.Parse("2006-01-02", input.EndDate)

	if err1 != nil || err2 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid date format. Use YYYY-MM-DD"})
	}

	period := model.AttendancePeriod{
		PeriodId:  input.PeriodId,
		StartDate: start,
		EndDate:   end,
	}
	db.Create(&period)

	c.JSON(http.StatusOK, gin.H{"message": "Attendance period created", "data": period})
}

func GetPayrollSummary(c *gin.Context, db *gorm.DB) {
	periodId := c.Param("period_id")

	var payslip []model.Payslip
	if err := db.Where("period_id = ?", periodId).Find(&payslip).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get payroll summary"})
		return
	}

	type SummaryRow struct {
		EmployeeId    string
		TotalPay      float64
		Attendance    int
		Overtime      float64
		Reimbursement float64
	}

	summary := []SummaryRow{}
	var total float64 = 0

	for _, p := range payslip {
		summary = append(summary, SummaryRow{
			EmployeeId:    p.EmployeeID,
			TotalPay:      p.TotalPay,
			Attendance:    p.AttendanceDays,
			Overtime:      p.OvertimeHours,
			Reimbursement: p.Reimbursement,
		})
		total += p.TotalPay
	}

	c.JSON(http.StatusOK, gin.H{
		"summary":       summary,
		"take_home_pay": total,
	})
}

func businessDays(start, end time.Time) int {
	days := 0
	for d := start; !d.After(end); d = d.AddDate(0, 0, 1) {
		if d.Weekday() != time.Saturday && d.Weekday() != time.Sunday {
			days++
		}
	}
	return days
}

func RunPayroll(c *gin.Context, db *gorm.DB) {
	var period model.AttendancePeriod

	if err := db.Where("period_id = ?", c.Param("period_id")).First(&period).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Period not found"})
		return
	}

	var existingPayroll model.Payslip
	if err := db.Where("period_id = ?", period.PeriodId).First(&existingPayroll).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Payroll already run for this period"})
		return
	}

	var employees []model.Employee
	db.Find(&employees)

	for _, emp := range employees {
		var attendanceCount int64
		db.Model(&model.Overtime{}).
			Where("employee_id = ? AND date BETWEEN ? AND ?", emp.EmployeeId, period.StartDate, period.EndDate).
			Count(&attendanceCount)

		var totalOvertime float64
		db.Model(&model.Overtime{}).
			Where("employee_id = ? AND date BETWEEN ? AND ?", emp.EmployeeId, period.StartDate, period.EndDate).
			Select("COALESCE(SUM(hours),0)").Scan(&totalOvertime)

		var totalReimbursement float64
		db.Model(model.Reimbursement{}).
			Where("employee_id = ? AND created_at BETWEEN ? AND ?", emp.EmployeeId, period.StartDate, period.EndDate).
			Select("COALESCE(SUM(amount),0)").Scan(&totalReimbursement)

		daysInMonth := businessDays(period.StartDate, period.EndDate)
		dailyRate := emp.Salary / float64(daysInMonth)
		attendancePay := dailyRate * float64(attendanceCount)
		overtimePay := (dailyRate / 8.0) * totalOvertime * 2.0

		total := attendancePay + overtimePay + totalReimbursement

		fmt.Println("Total : ", total, "Total Reimburse : ", totalReimbursement, "Attendance Pay :", attendancePay)
		db.Create(&model.Payslip{
			EmployeeID:     emp.EmployeeId,
			PeriodID:       period.PeriodId,
			AttendanceDays: int(attendanceCount),
			OvertimeHours:  totalOvertime,
			Reimbursement:  totalReimbursement,
			TotalPay:       total,
			GeneratedAt:    time.Now(),
		})
	}

	c.JSON(http.StatusOK, gin.H{"message": "Payroll processed"})
}
