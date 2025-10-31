package monitor

import (
	"check-monitor/config"
	"fmt"
	"health-check-monitor/config"
	"net/http"
	"time"
)

type CheckResult struct {
	Endpoint   config.Endpoint
	status     string
	StatusCode int
	ResponseTime time.Duration
	Error      error
	Timestamp time.Time
}

func CheckEndpoint(endpoint config.Endpoint) CheckResult {
	start := time.Now()

	client := &http.Client{
		Timeout: time.Duration(endpoint.Timeout) * time.Second,
	}

	req, err := http.NewRequest(endpoint.Method, endpoint.URL, nil)

	if err != nil {
		return CheckResult{
			Endpoint: endpoint,
			status:   "ERROR",
			Error:    err,
			Timestamp: time.Now(),
		}
	}

	resp, err := client.Do(req)
	responseTime := time.Since(start)

	result := CheckResult{
		Endpoint:   endpoint,
		ResponseTime: responseTime,
		Timestamp: time.Now(),
	}

	if err != nil {
		result.Status = "DOWN"
		result.Error = err
		return result
	}
	defer resp.Body.Close()

	result.StatusCode = resp.StatusCode

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		result.Status = "UP"
	} else {
		result.Status = "DEGRADED"
	}

	return result
}

func (r CheckResult) Print() {
    timestamp := r.Timestamp.Format("2006-01-02 15:04:05")
    
    if r.Status == "UP" {
        fmt.Printf("[%s] ✓ %s | %s | %d | %v\n", 
            timestamp, r.Endpoint.Name, r.Status, r.StatusCode, r.ResponseTime)
    } else if r.Status == "DEGRADED" {
        fmt.Printf("[%s] ⚠ %s | %s | %d | %v\n", 
            timestamp, r.Endpoint.Name, r.Status, r.StatusCode, r.ResponseTime)
    } else {
        fmt.Printf("[%s] ✗ %s | %s | Error: %v\n", 
            timestamp, r.Endpoint.Name, r.Status, r.Error)
    }
}