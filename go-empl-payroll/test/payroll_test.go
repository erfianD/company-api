package test

import (
	"bytes"
	"encoding/json"
	"go-empl-payroll/config"
	"go-empl-payroll/routes"

	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func TestPayrollProcess(t *testing.T) {
	db := config.ConnectDatabase()
	r := gin.Default()
	routes.SetupRouter(r, db)

	for i := 1; i <= 100; i++ {
		employeeID := "employee_" + strconv.Itoa(i)

		payload := map[string]interface{}{
			"employee_id": employeeID,
			"date":        time.Now().Format("2006-01-02"), // hari ini
			"hours":       0,
		}
		jsonPayload, _ := json.Marshal(payload)

		req := httptest.NewRequest("POST", "/employee/attendance", bytes.NewBuffer(jsonPayload))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Failed to submit attendance for %s. Status: %d, Body: %s", employeeID, w.Code, w.Body.String())
		}
	}

	for i := 1; i <= 100; i++ {
		employeeID := "employee_" + strconv.Itoa(i)

		payload := map[string]interface{}{
			"employee_id": employeeID,
			"date":        time.Now().Format("2006-01-02"), // hari ini
			"hours":       3,
		}
		jsonPayload, _ := json.Marshal(payload)

		req := httptest.NewRequest("POST", "/employee/overtime", bytes.NewBuffer(jsonPayload))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Failed to submit attendance for %s. Status: %d, Body: %s", employeeID, w.Code, w.Body.String())
		}
	}

	for i := 1; i <= 100; i++ {
		employeeID := "employee_" + strconv.Itoa(i)

		payload := map[string]interface{}{
			"employee_id": employeeID,
			"amount":      200000,
			"description": "Migration",
		}
		jsonPayload, _ := json.Marshal(payload)

		req := httptest.NewRequest("POST", "/employee/reimbursement", bytes.NewBuffer(jsonPayload))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Failed to submit attendance for %s. Status: %d, Body: %s", employeeID, w.Code, w.Body.String())
		}
	}

}
