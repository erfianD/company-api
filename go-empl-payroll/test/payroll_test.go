package test

import (
	"bytes"
	"encoding/json"
	"go-empl-payroll/config"
	"go-empl-payroll/routes"
	"strconv"
	"strings"
	"time"

	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestPayrollProcess(t *testing.T) {
	db := config.ConnectDatabase()
	r := gin.Default()
	routes.SetupRouter(r, db)

	req1 := httptest.NewRequest("POST", "/admin/attendance-period", strings.NewReader(`{
	    "period_id" : "payroll-2025-06",
		"start_date": "2025-06-01",
		"end_date": "2025-06-30"
	}`))
	res1 := httptest.NewRecorder()
	r.ServeHTTP(res1, req1)
	if res1.Code != http.StatusOK {
		t.Errorf("Failed to create period. Status: %d, Body: %s", res1.Code, res1.Body.String())
	}

	for i := 1; i <= 100; i++ {
		employeeID := "employee_" + strconv.Itoa(i)

		payload := map[string]interface{}{
			"employee_id": employeeID,
			"date":        time.Now().Format("2006-01-02"),
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
			t.Errorf("Failed to submit overtime for %s. Status: %d, Body: %s", employeeID, w.Code, w.Body.String())
		}
	}

	for i := 1; i <= 100; i++ {
		employeeID := "employee_" + strconv.Itoa(i)

		payload := map[string]interface{}{
			"employee_id": employeeID,
			"amount":      200000,
			"description": "Migration",
			"period_id":   "payroll-2025-06",
		}
		jsonPayload, _ := json.Marshal(payload)

		req := httptest.NewRequest("POST", "/employee/reimbursement", bytes.NewBuffer(jsonPayload))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Failed to submit reimbursement for %s. Status: %d, Body: %s", employeeID, w.Code, w.Body.String())
		}
	}

	req5 := httptest.NewRequest("POST", "/admin/payroll/payroll-2025-06", nil)
	res5 := httptest.NewRecorder()
	r.ServeHTTP(res5, req5)
	if res5.Code != http.StatusOK {
		t.Errorf("Failed to run payroll period. Status: %d, Body: %s", res5.Code, res5.Body.String())
	}

	reqSum := httptest.NewRequest("POST", "/admin/payroll/summary/payroll-2025-06", nil)
	resSum := httptest.NewRecorder()
	r.ServeHTTP(resSum, reqSum)
	if resSum.Code != http.StatusOK {
		t.Errorf("Failed to run payroll period. Status: %d, Body: %s", resSum.Code, resSum.Body.String())
	}

	req6 := httptest.NewRequest("GET", "/employee/payslip/payroll-2025-06/employee_40", nil)
	res6 := httptest.NewRecorder()
	r.ServeHTTP(res6, req6)
	if res6.Code != http.StatusOK {
		t.Errorf("Payslip not found. Status: %d, Body: %s", res6.Code, res6.Body.String())
	}
}
