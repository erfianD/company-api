package config

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDatabase() *gorm.DB {
	dsn := "host=localhost user=postgres password=P@ssw0rd dbname=company port=5432 sslmode=disable"
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database!")
	}
	// database.AutoMigrate(
	// 	&model.Admin{},
	// 	&model.AttendancePeriod{},
	// 	&model.Employee{},
	// 	&model.Payroll{},
	// 	&model.Overtime{},
	// 	&model.Reimbursement{})

	// DB = database
	return database
}
