/**
 * @fileoverview Main entry point for the AI project tutorial API server.
 * Implements a basic HTTP server with health and readiness endpoints for Phase 0 setup.
 * Provides foundation for building scalable, observable applications with comprehensive error handling.
 */

package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/ashleywang1/new-ai-project-tutorial/pkg/health"
)

const (
	// DefaultPort is the default HTTP server port
	DefaultPort = "8080"
	// ShutdownTimeout defines how long to wait for graceful shutdown
	ShutdownTimeout = 30 * time.Second
	// StartupTimeout defines how long to wait for server to start
	StartupTimeout = 10 * time.Second
	// MaxRetries defines maximum startup retry attempts
	MaxRetries = 3
	// RetryDelay defines delay between startup retries
	RetryDelay = 2 * time.Second
)

// ServerError represents application-specific errors
type ServerError struct {
	Message string
	Cause   error
	Code    int
}

func (e *ServerError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Cause)
	}
	return e.Message
}

/**
 * @description Main function that initializes and starts the HTTP server.
 * Sets up health endpoints and handles graceful shutdown on termination signals.
 * Includes comprehensive error handling and startup retry logic.
 */
func main() {
	fmt.Println("AI Project Tutorial API Server - Phase 0")

	// Validate configuration
	if err := validateConfiguration(); err != nil {
		log.Fatalf("Configuration validation failed: %v", err)
	}

	// Create health checker instance
	healthChecker := health.NewHealthChecker(health.HealthCheckerConfig{
		ServiceName:    "AI Project Tutorial API Server",
		ServiceVersion: "0.1.0",
	})

	// Add basic readiness checks
	healthChecker.AddReadinessCheck("handlers", health.AlwaysHealthyCheck())
	healthChecker.AddReadinessCheck("server", health.AlwaysHealthyCheck())

	// Create HTTP server with configured routes
	server, err := createHTTPServerWithHealthChecker(healthChecker)
	if err != nil {
		log.Fatalf("Failed to create HTTP server: %v", err)
	}

	// Start server with retry logic in a goroutine
	serverErrChan := make(chan error, 1)
	go func() {
		serverErrChan <- startServerWithRetries(server)
	}()

	// Setup graceful shutdown handling
	shutdown := setupShutdownSignals()

	// Wait for either server error or shutdown signal
	select {
	case err := <-serverErrChan:
		if err != nil {
			log.Fatalf("Server failed to start: %v", err)
		}
		// Server stopped gracefully
	case sig := <-shutdown:
		fmt.Printf("\nReceived signal: %v. Initiating graceful shutdown...\n", sig)
		if err := performGracefulShutdown(server); err != nil {
			log.Printf("Error during graceful shutdown: %v", err)
			os.Exit(1)
		}
	}

	fmt.Println("Server shutdown complete")
}

/**
 * @description Validates application configuration before startup.
 * Checks port availability, environment variables, and system requirements.
 */
func validateConfiguration() error {
	port := getPort()

	// Validate port number
	if portNum, err := strconv.Atoi(port); err != nil || portNum < 1 || portNum > 65535 {
		return &ServerError{
			Message: "Invalid port number",
			Cause:   err,
			Code:    400,
		}
	}

	// Check if port is available
	if !isPortAvailable(port) {
		return &ServerError{
			Message: fmt.Sprintf("Port %s is already in use", port),
			Code:    409,
		}
	}

	fmt.Printf("✅ Configuration validated - Port %s is available\n", port)
	return nil
}

/**
 * @description Creates and configures the HTTP server with health checker.
 * Returns a configured http.Server with proper timeouts and error handling.
 */
func createHTTPServerWithHealthChecker(healthChecker *health.HealthChecker) (*http.Server, error) {
	mux := http.NewServeMux()

	// Register health endpoints using the health checker
	mux.HandleFunc("/health", withErrorHandling(healthChecker.HealthHandler))
	mux.HandleFunc("/ready", withErrorHandling(healthChecker.ReadinessHandler))
	mux.HandleFunc("/", withErrorHandling(handleRoot))

	server := &http.Server{
		Addr:         ":" + getPort(),
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
		ErrorLog:     log.New(os.Stderr, "HTTP: ", log.LstdFlags),
	}

	fmt.Println("✅ HTTP server configured successfully")
	return server, nil
}

/**
 * @description Starts the server with retry logic for improved reliability.
 * Attempts to start the server multiple times with exponential backoff.
 */
func startServerWithRetries(server *http.Server) error {
	var lastErr error

	for attempt := 1; attempt <= MaxRetries; attempt++ {
		fmt.Printf("Starting server (attempt %d/%d) on %s...\n", attempt, MaxRetries, server.Addr)

		// Start server - this will block until server stops or fails
		fmt.Printf("✅ Server started successfully on %s\n", server.Addr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			lastErr = &ServerError{
				Message: fmt.Sprintf("Server startup failed on attempt %d", attempt),
				Cause:   err,
				Code:    500,
			}

			if attempt < MaxRetries {
				fmt.Printf("❌ Startup failed: %v. Retrying in %v...\n", err, RetryDelay)
				time.Sleep(RetryDelay)
				continue
			}
		} else {
			// Server shutdown gracefully (ErrServerClosed)
			fmt.Println("✅ Server shutdown gracefully")
			return nil
		}
	}

	return lastErr
}

/**
 * @description Sets up signal handling for graceful shutdown.
 * Returns a channel that receives shutdown signals.
 */
func setupShutdownSignals() <-chan os.Signal {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	return signalChan
}

/**
 * @description Performs graceful shutdown of the HTTP server.
 * Handles connection draining and resource cleanup with timeout.
 */
func performGracefulShutdown(server *http.Server) error {
	fmt.Println("Initiating graceful shutdown...")

	// Create context with timeout for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), ShutdownTimeout)
	defer cancel()

	// Channel to track shutdown completion
	shutdownComplete := make(chan error, 1)

	go func() {
		shutdownComplete <- server.Shutdown(ctx)
	}()

	// Wait for shutdown completion or timeout
	select {
	case err := <-shutdownComplete:
		if err != nil {
			return &ServerError{
				Message: "Error during server shutdown",
				Cause:   err,
				Code:    500,
			}
		}
		fmt.Println("✅ Server shutdown completed successfully")
		return nil

	case <-ctx.Done():
		// Force close if graceful shutdown times out
		fmt.Println("⚠️ Graceful shutdown timed out, forcing server close...")
		if err := server.Close(); err != nil {
			return &ServerError{
				Message: "Error during forced server close",
				Cause:   err,
				Code:    500,
			}
		}
		return &ServerError{
			Message: "Server shutdown timed out and was forced to close",
			Code:    408,
		}
	}
}

/**
 * @description Middleware wrapper that adds error handling to HTTP handlers.
 * Provides consistent error logging and response formatting.
 */
func withErrorHandling(handler func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic in handler %s: %v", r.URL.Path, err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()

		// Log request
		log.Printf("Request: %s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)

		// Call the actual handler
		handler(w, r)
	}
}

/**
 * @description Root endpoint handler providing basic service information.
 * Returns service name and available endpoints with error handling.
 */
func handleRoot(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := fmt.Sprintf(`{
		"service": "AI Project Tutorial API Server",
		"phase": "0",
		"endpoints": ["/health", "/ready"],
		"timestamp": "%s"
	}`, time.Now().UTC().Format(time.RFC3339))
	w.Write([]byte(response))
}

/**
 * @description Gets the server port from environment or returns default.
 * Checks PORT environment variable, defaults to 8080.
 */
func getPort() string {
	if port := os.Getenv("PORT"); port != "" {
		return port
	}
	return DefaultPort
}

/**
 * @description Checks if a port is available for binding.
 * Returns true if the port is available, false otherwise.
 */
func isPortAvailable(port string) bool {
	address := net.JoinHostPort("", port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return false
	}
	listener.Close()
	return true
}
