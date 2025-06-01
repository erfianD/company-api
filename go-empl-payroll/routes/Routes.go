package routes

import (
	"go-empl-payroll/handler"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(r *gin.Engine, db *gorm.DB) {

	r.POST("/login", handler.Login(db))

	admin := r.Group("/admin")
	{
		admin.POST("/attendance-period", func(c *gin.Context) {
			handler.CreateAttndancePeriod(c, db)
		})
		admin.POST("/payroll/:period_id", func(c *gin.Context) {
			handler.RunPayroll(c, db)
		})
		admin.GET("/payroll/summary/:period_id", func(c *gin.Context) {
			handler.GetPayrollSummary(c, db)
		})
	}

	employee := r.Group("/employee")
	{
		employee.POST("/attendance", func(c *gin.Context) {
			handler.SubmitAttendance(c, db)
		})
		employee.POST("/overtime", func(c *gin.Context) {
			handler.SubmitOvertime(c, db)
		})
		employee.POST("/reimbursement", func(c *gin.Context) {
			handler.SubmitReimbursement(c, db)
		})
		employee.GET("/payslip/:period_id", func(c *gin.Context) {
			handler.GetPayslip(c, db)
		})
	}
}
