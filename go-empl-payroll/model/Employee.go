package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Employee struct {
	gorm.Model
	ID         uuid.UUID
	EmployeeId string
	Username   string
	Password   string
	Role       string
	IsAdmin    bool
	Token      string
	Salary     float64
}
