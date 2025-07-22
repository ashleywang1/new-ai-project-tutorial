# Project Rules

**@fileoverview** Comprehensive development guidelines for our AI-first Kubernetes-native telemetry project. This document establishes directory structure, naming conventions, coding standards, and AI compatibility requirements to ensure a modular, scalable, and maintainable codebase.

---

## Table of Contents

1. [AI-First Development Principles](#ai-first-development-principles)
2. [Directory Structure](#directory-structure)
3. [File Naming Conventions](#file-naming-conventions)
4. [Code Organization Standards](#code-organization-standards)
5. [Technology-Specific Guidelines](#technology-specific-guidelines)
6. [Documentation Requirements](#documentation-requirements)
7. [API Design Standards](#api-design-standards)

---

## AI-First Development Principles

### Core Requirements
- **Maximum File Size**: 500 lines per file (strict limit)
- **File Documentation**: Every file must include `@fileoverview` at the top
- **Function Documentation**: All functions require JSDoc/TSDoc/GoDoc block comments
- **Descriptive Naming**: Use clear, descriptive names for files, functions, and variables
- **Modular Design**: Prefer composition over inheritance, functional over imperative

### AI Compatibility Standards
- **Semantic Search Optimization**: Use descriptive variable names with auxiliary verbs (`isLoading`, `hasError`, `shouldRetry`)
- **Grep-Friendly Patterns**: Consistent naming patterns for easy regex searches
- **Context-Rich Comments**: Block comments explaining business logic and architectural decisions
- **Explicit Dependencies**: Clear import statements and dependency declarations

---

## Directory Structure

### Backend (Go)
```
/
├── cmd/
│   └── apiserver/           # Main application entry points
├── pkg/
│   ├── api/                 # gRPC service implementations
│   │   ├── metrics/         # Metrics endpoint handlers
│   │   ├── traces/          # Tracing endpoint handlers
│   │   ├── crds/           # CRD browser handlers
│   │   └── query/          # Raw query interface
│   ├── clickhouse/         # ClickHouse client and query builders
│   ├── auth/               # RBAC and authentication
│   ├── health/             # Health checks and readiness probes
│   └── telemetry/          # OpenTelemetry instrumentation
├── proto/                  # Protocol buffer definitions
├── deployments/
│   ├── helm/               # Helm charts
│   └── docker/             # Dockerfiles and build scripts
└── tests/
    ├── integration/        # Integration tests
    └── unit/              # Unit tests
```

### Frontend (React/TypeScript)
```
/frontend/
├── src/
│   ├── components/         # Reusable UI components
│   │   ├── common/         # Shared components
│   │   ├── metrics/        # Metrics visualization components
│   │   ├── traces/         # Tracing components
│   │   └── crds/          # CRD browser components
│   ├── pages/             # Route-level page components
│   ├── hooks/             # Custom React hooks
│   ├── services/          # gRPC client services
│   ├── types/             # TypeScript type definitions
│   ├── utils/             # Utility functions
│   └── store/             # State management
├── public/                # Static assets
└── tests/                # Frontend tests
```

### Documentation
```
/_docs/
├── resources/             # Images, diagrams, examples
├── api/                  # API documentation
├── deployment/           # Deployment guides
└── development/          # Development setup guides
```

---

## File Naming Conventions

### General Principles
- **Kebab-case** for directories: `custom-resources/`, `query-builder/`
- **PascalCase** for React components: `MetricsExplorer.tsx`, `TraceWaterfall.tsx`
- **camelCase** for TypeScript/JavaScript files: `authService.ts`, `queryUtils.ts`
- **snake_case** for Go files: `clickhouse_client.go`, `metrics_handler.go`
- **lowercase** for configuration: `dockerfile`, `makefile`, `helm.yaml`

### Descriptive Naming Patterns
- **Handlers**: `*_handler.go` (e.g., `metrics_handler.go`)
- **Services**: `*Service.ts` (e.g., `authService.ts`, `grpcService.ts`)
- **Components**: Noun-based names (e.g., `Dashboard.tsx`, `QueryEditor.tsx`)
- **Hooks**: `use*` prefix (e.g., `useMetrics.ts`, `useAuthentication.ts`)
- **Types**: Domain-specific suffixes (e.g., `MetricsTypes.ts`, `TraceModels.ts`)
- **Tests**: `*_test.go` or `*.test.ts` matching the file under test

---

## Code Organization Standards

### File Structure Template
```typescript
/**
 * @fileoverview [Brief description of file purpose and main exports]
 * 
 * @description [Detailed explanation of the module's role in the system]
 * @author [Team/Individual responsible]
 * @since [Version/Date when file was created]
 */

// External imports (libraries, frameworks)
import { ... } from 'external-library';

// Internal imports (other project modules)
import { ... } from '../internal/module';

// Type definitions (if not in separate file)
interface LocalInterface { ... }

// Constants and configuration
const CONFIG = { ... };

// Main implementation
export function primaryFunction() { ... }

// Helper functions (private to module)
function helperFunction() { ... }

// Default export (if applicable)
export default MainComponent;
```

### Function Documentation Standards
```typescript
/**
 * @description Retrieves cluster metrics from ClickHouse with time-based filtering
 * @param {string} clusterId - Unique identifier for the target cluster
 * @param {TimeRange} timeRange - Start and end timestamps for data retrieval
 * @param {MetricType[]} metricTypes - Array of specific metrics to fetch
 * @returns {Promise<MetricsResponse>} Aggregated metrics data with timestamps
 * @throws {ClickHouseError} When database connection fails or query times out
 * @example
 * ```typescript
 * const metrics = await fetchClusterMetrics('prod-cluster', {
 *   start: new Date('2024-01-01'),
 *   end: new Date('2024-01-02')
 * }, ['cpu', 'memory']);
 * ```
 */
```

### Go Documentation Standards
```go
// Package metrics provides handlers for cluster telemetry data retrieval.
//
// This package implements the gRPC service definitions for metrics endpoints,
// handling authentication, query validation, and ClickHouse integration.
package metrics

// FetchClusterMetrics retrieves aggregated cluster metrics from ClickHouse.
//
// The function validates the time range, applies RBAC filters based on the
// service account, and executes optimized ClickHouse queries for the requested
// metric types.
//
// Parameters:
//   - ctx: Request context with authentication and tracing information
//   - req: MetricsRequest containing cluster ID, time range, and metric types
//
// Returns:
//   - MetricsResponse: Aggregated metrics with timestamps and metadata
//   - error: ClickHouseError for database issues, ValidationError for invalid input
func FetchClusterMetrics(ctx context.Context, req *MetricsRequest) (*MetricsResponse, error) {
    // Implementation
}
```

---

## Technology-Specific Guidelines

### Go Backend Standards
- **Package Organization**: One package per API domain (metrics, traces, crds)
- **Error Handling**: Use wrapped errors with context (`fmt.Errorf("operation failed: %w", err)`)
- **Context Propagation**: Pass `context.Context` as first parameter in all functions
- **Struct Naming**: Use noun-based names (`MetricsHandler`, `ClickHouseClient`)
- **Interface Definitions**: Small, focused interfaces (`QueryExecutor`, `AuthValidator`)

### React/TypeScript Frontend Standards
- **Component Types**: Functional components with TypeScript interfaces for props
- **State Management**: React Query for server state, local state for UI-only data
- **Event Handlers**: Prefix with `handle` or `on` (`handleSubmit`, `onMetricSelect`)
- **Custom Hooks**: Encapsulate business logic and API calls
- **Type Safety**: Strict TypeScript configuration, no `any` types in production code

### gRPC Protocol Standards
- **Service Definitions**: Domain-based service grouping (`MetricsService`, `TracingService`)
- **Message Naming**: Descriptive request/response pairs (`GetMetricsRequest`, `GetMetricsResponse`)
- **Field Validation**: Use proto validation rules for input constraints
- **Versioning**: Include version in package names (`v1.MetricsService`)

### ClickHouse Integration Standards
- **Query Builders**: Composable query construction with parameter binding
- **Connection Pooling**: Shared connection pool with health monitoring
- **Query Optimization**: Use materialized views and proper indexing strategies
- **Error Mapping**: Convert ClickHouse errors to domain-specific error types

---

## Documentation Requirements

### README Standards
- **Project Overview**: Clear description of purpose and key features
- **Quick Start**: Step-by-step setup instructions for local development
- **API Documentation**: Links to generated gRPC documentation
- **Deployment Guide**: Helm chart usage and configuration options

### API Documentation
- **Endpoint Documentation**: Auto-generated from proto files
- **Usage Examples**: cURL commands and code snippets for each endpoint
- **Error Codes**: Comprehensive list of possible error responses
- **Rate Limiting**: Documentation of any API limits or throttling

### Architecture Documentation
- **System Diagrams**: Mermaid diagrams showing component relationships
- **Data Flow**: Visual representation of request/response cycles
- **Security Model**: RBAC configuration and authentication flows
- **Deployment Architecture**: Kubernetes manifests and networking setup

---

## API Design Standards

### gRPC Service Design
- **Stateless Services**: No server-side session state
- **Idempotent Operations**: Safe retry behavior for all read operations
- **Resource-Based URLs**: RESTful-style resource identification
- **Consistent Error Handling**: Standardized error codes and messages

### Authentication & Authorization
- **JWT Tokens**: Short-lived tokens with proper expiration handling
- **RBAC Integration**: Kubernetes ServiceAccount-based permissions
- **Request Tracing**: Correlation IDs for debugging and monitoring
- **Audit Logging**: Security-relevant events with proper context

### Performance Guidelines
- **Response Pagination**: Limit large result sets with cursor-based pagination
- **Caching Strategy**: Cache static data with appropriate TTL values
- **Query Optimization**: Use database query optimization and explain plans
- **Resource Limits**: Set appropriate timeouts and connection limits

---

## Quality Assurance

### Testing Requirements
- **Unit Test Coverage**: Minimum 80% coverage for all packages
- **Integration Tests**: End-to-end API testing with real ClickHouse instance
- **Load Testing**: Performance validation under expected traffic patterns
- **Security Testing**: Authentication and authorization verification

### Code Review Standards
- **AI-First Review**: Ensure code is optimized for AI tool compatibility
- **Performance Impact**: Review ClickHouse query efficiency and caching
- **Security Review**: Validate authentication, authorization, and input validation
- **Documentation Review**: Verify all functions have proper documentation

---

This document serves as the foundation for maintaining code quality and AI compatibility throughout the project lifecycle. All team members should reference these guidelines when contributing to the codebase. 