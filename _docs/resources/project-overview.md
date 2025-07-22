## Project Overview

This project delivers a Kubernetes-native API server that provides unified access to cluster telemetry (metrics, traces, and CRDs) stored in ClickHouse, serving as the data backbone for an operational UI dashboard.

**Key Components:**
- **Containerized API Server**: Go-based application deployed as a Docker container in Kubernetes
- **ClickHouse Integration**: High-performance connection to ClickHouse for time-series and operational data
- **Kubernetes-Native**: Built to run within cluster environments with proper RBAC integration
- **Observability Focus**: Unified access to metrics, tracing, and configuration data

---

## Phase 1: Core Implementation

### Objectives:
1. **Containerized API Foundation**
   - Develop Go-based apiserver with health checks, readiness probes, and metrics endpoints
   - Implement multi-architecture Docker build pipeline (amd64/arm64)
   - Configure Helm chart for Kubernetes deployment

2. **ClickHouse Data Gateway**
   - Establish optimized ClickHouse connection pool
   - Implement query builders for:
     - Cluster metrics aggregation
     - Distributed tracing data
     - Custom Resource Definition (CRD) state
   - Add query instrumentation for performance monitoring

3. **Core API Endpoints**
   - `/metrics` - Cluster health and performance indicators
   - `/traces` - Distributed tracing query interface
   - `/crds` - Custom resource state browser
   - `/query` - Raw ClickHouse query interface (with safeguards)

4. **Operational Foundations**
   - Structured logging with request correlation
   - OpenTelemetry instrumentation
   - RBAC integration using Kubernetes ServiceAccounts
   - Configuration via ConfigMaps/Secrets