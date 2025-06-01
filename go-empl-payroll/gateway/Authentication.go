package gateway

import (
	"go-empl-payroll/model"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func AuthGateway(requiredRole string, db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		token = strings.TrimPrefix(token, "Bearer")

		var employees model.Employee
		if err := db.Where("token = ?", token).First(&employees).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
			c.Abort()
			return
		}

		if requiredRole == "admin" && !employees.IsAdmin {
			c.JSON(http.StatusForbidden, gin.H{"message": "Forbidden. Access Denied"})
			c.Abort()
			return
		}

		c.Set("employee_id", employees.ID)
		c.Next()
	}
}

func GetEmployeeId(c *gin.Context) uuid.UUID {
	if id, exists := c.Get("employee_id"); exists {
		return id.(uuid.UUID)
	}
	return uuid.Nil
}
