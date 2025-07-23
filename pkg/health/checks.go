/**
 * @fileoverview Common health check implementations for typical application dependencies.
 * Provides ready-to-use check functions for databases, external services, and system resources.
 * Designed to be composable and easily integrated with the HealthChecker.
 */

package health

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"time"
)

/**
 * @description Creates a check that verifies if a TCP port is available for binding.
 * Useful for checking if the application's port is ready to accept connections.
 */
func PortAvailableCheck(port string) CheckFunc {
	return func() error {
		address := net.JoinHostPort("", port)
		listener, err := net.Listen("tcp", address)
		if err != nil {
			return fmt.Errorf("port %s is not available: %w", port, err)
		}
		listener.Close()
		return nil
	}
}

/**
 * @description Creates a check that verifies if a TCP connection can be established to a host:port.
 * Useful for checking database connections, external service dependencies, etc.
 */
func TCPConnectionCheck(host, port string, timeout time.Duration) CheckFunc {
	return func() error {
		address := net.JoinHostPort(host, port)
		conn, err := net.DialTimeout("tcp", address, timeout)
		if err != nil {
			return fmt.Errorf("failed to connect to %s: %w", address, err)
		}
		conn.Close()
		return nil
	}
}

/**
 * @description Creates a check that performs an HTTP GET request to verify service availability.
 * Useful for checking external HTTP dependencies and health endpoints.
 */
func HTTPCheck(url string, timeout time.Duration, expectedStatusCode int) CheckFunc {
	return func() error {
		client := &http.Client{
			Timeout: timeout,
		}

		resp, err := client.Get(url)
		if err != nil {
			return fmt.Errorf("HTTP request failed to %s: %w", url, err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != expectedStatusCode {
			return fmt.Errorf("unexpected status code from %s: got %d, expected %d",
				url, resp.StatusCode, expectedStatusCode)
		}

		return nil
	}
}

/**
 * @description Creates a simple check that always returns healthy.
 * Useful for basic health endpoints when no specific checks are needed.
 */
func AlwaysHealthyCheck() CheckFunc {
	return func() error {
		return nil
	}
}

/**
 * @description Creates a check that validates required environment variables are set.
 * Useful for ensuring configuration is properly loaded.
 */
func EnvironmentVariableCheck(envVars []string) CheckFunc {
	return func() error {
		for _, envVar := range envVars {
			if value := getEnvVar(envVar); value == "" {
				return fmt.Errorf("required environment variable %s is not set", envVar)
			}
		}
		return nil
	}
}

/**
 * @description Creates a composite check that runs multiple checks and fails if any fail.
 * Useful for grouping related checks together.
 */
func CompositeCheck(name string, checks ...CheckFunc) CheckFunc {
	return func() error {
		for i, check := range checks {
			if err := check(); err != nil {
				return fmt.Errorf("%s check %d failed: %w", name, i+1, err)
			}
		}
		return nil
	}
}

// Helper function to get environment variable
func getEnvVar(key string) string {
	return os.Getenv(key)
}
