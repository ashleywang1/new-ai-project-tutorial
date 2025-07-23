/**
 * @fileoverview Health check package providing configurable health and readiness endpoints.
 * Supports custom health checks, uptime tracking, and structured JSON responses for monitoring.
 * Designed for scalable applications requiring comprehensive health monitoring.
 */

package health

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// HealthChecker provides health and readiness check functionality
type HealthChecker struct {
	serviceName     string
	serviceVersion  string
	startTime       time.Time
	readinessChecks map[string]CheckFunc
	healthChecks    map[string]CheckFunc
}

// CheckFunc represents a health check function that returns an error if unhealthy
type CheckFunc func() error

// CheckResult represents the result of a health check
type CheckResult struct {
	Status    string            `json:"status"`
	Checks    map[string]string `json:"checks,omitempty"`
	Timestamp string            `json:"timestamp"`
	Uptime    string            `json:"uptime,omitempty"`
	Service   string            `json:"service,omitempty"`
	Version   string            `json:"version,omitempty"`
}

// HealthCheckerConfig provides configuration options for the health checker
type HealthCheckerConfig struct {
	ServiceName    string
	ServiceVersion string
}

/**
 * @description Creates a new HealthChecker instance with the provided configuration.
 * Initializes check maps and sets the start time for uptime calculations.
 */
func NewHealthChecker(config HealthCheckerConfig) *HealthChecker {
	return &HealthChecker{
		serviceName:     config.ServiceName,
		serviceVersion:  config.ServiceVersion,
		startTime:       time.Now(),
		readinessChecks: make(map[string]CheckFunc),
		healthChecks:    make(map[string]CheckFunc),
	}
}

/**
 * @description Adds a readiness check with the given name and check function.
 * Readiness checks determine if the service is ready to accept traffic.
 */
func (hc *HealthChecker) AddReadinessCheck(name string, check CheckFunc) {
	hc.readinessChecks[name] = check
}

/**
 * @description Adds a health check with the given name and check function.
 * Health checks determine if the service is functioning properly.
 */
func (hc *HealthChecker) AddHealthCheck(name string, check CheckFunc) {
	hc.healthChecks[name] = check
}

/**
 * @description HTTP handler for the health endpoint.
 * Returns service health status and executes all registered health checks.
 */
func (hc *HealthChecker) HealthHandler(w http.ResponseWriter, r *http.Request) {
	result := hc.performChecks(hc.healthChecks)
	result.Service = hc.serviceName
	result.Version = hc.serviceVersion
	result.Uptime = time.Since(hc.startTime).String()

	hc.writeJSONResponse(w, result, http.StatusOK)
}

/**
 * @description HTTP handler for the readiness endpoint.
 * Returns service readiness status and executes all registered readiness checks.
 */
func (hc *HealthChecker) ReadinessHandler(w http.ResponseWriter, r *http.Request) {
	result := hc.performChecks(hc.readinessChecks)

	// Set appropriate status code based on check results
	statusCode := http.StatusOK
	if result.Status != "healthy" {
		statusCode = http.StatusServiceUnavailable
	}

	hc.writeJSONResponse(w, result, statusCode)
}

/**
 * @description Performs all checks in the provided map and returns aggregated results.
 * Returns "healthy" status only if all checks pass, "unhealthy" otherwise.
 */
func (hc *HealthChecker) performChecks(checks map[string]CheckFunc) CheckResult {
	result := CheckResult{
		Status:    "healthy",
		Checks:    make(map[string]string),
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}

	// If no checks are configured, default to healthy
	if len(checks) == 0 {
		result.Checks["default"] = "ok"
		return result
	}

	// Execute all checks
	hasFailures := false
	for name, checkFunc := range checks {
		if err := checkFunc(); err != nil {
			result.Checks[name] = fmt.Sprintf("failed: %v", err)
			hasFailures = true
		} else {
			result.Checks[name] = "ok"
		}
	}

	if hasFailures {
		result.Status = "unhealthy"
	}

	return result
}

/**
 * @description Writes a JSON response with proper headers and error handling.
 * Sets content type and handles JSON marshaling errors gracefully.
 */
func (hc *HealthChecker) writeJSONResponse(w http.ResponseWriter, result CheckResult, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(result); err != nil {
		// Fallback to simple error response if JSON encoding fails
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"status":"error","message":"failed to encode response"}`)
	}
}

/**
 * @description Returns the service uptime as a duration since start.
 * Useful for external monitoring and debugging.
 */
func (hc *HealthChecker) GetUptime() time.Duration {
	return time.Since(hc.startTime)
}

/**
 * @description Returns the service start time.
 * Useful for external monitoring and debugging.
 */
func (hc *HealthChecker) GetStartTime() time.Time {
	return hc.startTime
}
