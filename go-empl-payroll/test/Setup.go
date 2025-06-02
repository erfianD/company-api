package test

import (
	"go-empl-payroll/model"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var TestDB *gorm.DB

func InitTestDB() {
	dsn := "host=localhost user=postgres password=P@ssw0rd dbname=company port=5432 sslmode=disable"
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database!")
	}
	database.AutoMigrate(
		&model.Admin{},
		&model.Attendance{},
		&model.AttendancePeriod{},
		&model.Employee{},
		&model.Payslip{},
		&model.Overtime{},
		&model.Reimbursement{})

	TestDB = database
}
