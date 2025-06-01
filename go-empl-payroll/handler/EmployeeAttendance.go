package handler

import (
	"net/http"
	"time"

	"go-empl-payroll/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AttendanceRequest struct {
	EmployeeID string `json:"employee_id" binding:"required"`
	Date       string `json:"date"`
	Hours      int    `json:"hours"`
}

func GetEmployeeID(c *gin.Context) string {
	employeeID, exists := c.Get("employee_id")
	if !exists {
		panic("employee_id not found in context")
	}
	return employeeID.(string)
}

func SubmitAttendance(c *gin.Context, db *gorm.DB) {
	var req AttendanceRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
		return
	}

	attendanceDate, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid date format"})
		return
	}

	weekday := attendanceDate.Weekday()
	if weekday == time.Saturday || weekday == time.Sunday {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Cannot submit on weekends"})
		return
	}

	var period model.AttendancePeriod
	err = db.Where("start_date <= ? AND end_date >= ? AND is_processed = false", attendanceDate, attendanceDate).
		First(&period).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "No Active attendance period for this date"})
		return
	}

	employeeId := req.EmployeeID
	var existing model.Attendance
	err = db.Where("employee_id = ? AND date = ?", employeeId, attendanceDate).
		First(&existing).Error
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Attendance already submitted for this date"})
		return
	}

	attendance := model.Attendance{
		EmployeeId: employeeId,
		Date:       attendanceDate,
		PeriodId:   period.PeriodId,
	}

	if err := db.Create(&attendance).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to save attendance."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Attendance submitted successfully"})
}

func SubmitOvertime(c *gin.Context, db *gorm.DB) {

	var input AttendanceRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request body."})
		return
	}

	date, err := time.Parse("2006-01-02", input.Date)
	if err != nil || input.Hours <= 0 || input.Hours > 3 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid overtime date or hours request."})
		return
	}

	employeeId := input.EmployeeID
	var exist model.Overtime
	if err := db.Where("employee_id = ? AND date = ?", employeeId, date).
		First(&exist).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Overtime already submitted for this date"})
		return
	}

	db.Create(&model.Overtime{
		EmployeeId: employeeId,
		Date:       date,
		Hours:      input.Hours,
	})

	c.JSON(http.StatusOK, gin.H{"message": "Overtime submitted"})
}

func SubmitReimbursement(c *gin.Context, db *gorm.DB) {

	type Input struct {
		EmployeeId  string  `json:"employee_id"`
		Amount      float64 `json:"amount"`
		Description string  `json:"description"`
	}

	var input Input
	if err := c.ShouldBindJSON(&input); err != nil || input.Amount <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request body."})
		return
	}

	db.Create(&model.Reimbursement{
		EmployeeId:  input.EmployeeId,
		Amount:      input.Amount,
		Description: input.Description,
	})

	c.JSON(http.StatusOK, gin.H{"message": "Reimbursement submitted."})
}

func GetPayslip(c *gin.Context, db *gorm.DB) {

	periodId := c.Param("period_id")
	employeeId := c.Param("employee_id")

	var payslip model.Payslip
	if err := db.Where("employee_id = ? AND period_id = ?", employeeId, periodId).First(&payslip).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Payslip not found."})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"attendance_days": payslip.AttendanceDays,
		"overtime_hours":  payslip.OvertimeHours,
		"reimbursement":   payslip.Reimbursement,
		"take_home_pay":   payslip.TotalPay,
	})
}
