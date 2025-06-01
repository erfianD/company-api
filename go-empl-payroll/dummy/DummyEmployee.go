package dummy

import (
	"go-empl-payroll/model"
	"math/rand"
	"strconv"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func InputData(db *gorm.DB) {

	db.Exec("DELETE FROM employees")

	for i := 1; i <= 100; i++ {
		password := "password_" + strconv.Itoa(i)
		hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

		db.Create(&model.Employee{
			ID:         uuid.New(),
			Username:   "employee_" + strconv.Itoa(i),
			EmployeeId: "employee_" + strconv.Itoa(i),
			Password:   string(hashed),
			Salary:     float64(rand.Intn(5_000_000) + 3_000_000),
			IsAdmin:    false,
		})
	}

	hashedAdmin, _ := bcrypt.GenerateFromPassword([]byte("adminPsswd"), bcrypt.DefaultCost)

	db.Create((&model.Employee{
		ID:         uuid.New(),
		Username:   "admin",
		Password:   string(hashedAdmin),
		EmployeeId: "admin",
		Salary:     0,
		IsAdmin:    true,
	}))
}
