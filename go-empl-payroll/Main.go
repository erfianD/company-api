package main

import (
	"go-empl-payroll/config"
	"go-empl-payroll/dummy"
	"go-empl-payroll/model"
	"go-empl-payroll/routes"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	db := config.ConnectDatabase()
	db.AutoMigrate(
		&model.Admin{},
		&model.Attendance{},
		&model.AttendancePeriod{},
		&model.Employee{},
		&model.Payslip{},
		&model.Overtime{},
		&model.Reimbursement{})

	r := gin.Default()
	dummy.InputData(db)
	routes.SetupRouter(r, db)
	r.Run(":8090")
	log.Println("Server is running.")
}
