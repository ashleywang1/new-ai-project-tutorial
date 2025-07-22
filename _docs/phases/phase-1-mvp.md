# Phase 1: MVP - Core Functionality

**@fileoverview** Minimal Viable Product phase implementing core telemetry API endpoints, basic React dashboard, Kubernetes deployment, and essential ClickHouse query capabilities.

---

## Objectives

Transform the foundation into a working product that delivers core value: unified access to cluster telemetry data through both API and web interface.

**Deliverables:**
- Functional API endpoints for metrics, traces, and CRDs
- Basic React dashboard with data visualization
- ClickHouse query builders for telemetry data
- Kubernetes deployment via Helm chart
- Basic authentication and RBAC integration

---

## Features & Tasks

### 1. Core API Endpoints Implementation
**Goal:** Implement functional gRPC services with real ClickHouse queries
- [ ] Build `/metrics` endpoint with cluster health and performance indicators
- [ ] Implement `/traces` endpoint for distributed tracing query interface
- [ ] Create `/crds` endpoint for custom resource state browsing
- [ ] Add `/query` endpoint for raw ClickHouse queries with basic safeguards
- [ ] Integrate request correlation and structured error responses

### 2. ClickHouse Query Engine
**Goal:** Robust query builders for telemetry data retrieval and aggregation
- [ ] Implement metrics aggregation queries (CPU, memory, network, error rates)
- [ ] Build distributed tracing query builders with span relationships
- [ ] Create CRD state queries with YAML parsing and status extraction
- [ ] Add query parameter validation and SQL injection protection
- [ ] Implement query performance monitoring and optimization

### 3. React Frontend Foundation
**Goal:** Basic dashboard UI for telemetry data visualization
- [ ] Setup React project with TypeScript and gRPC-Web integration
- [ ] Create layout components (header, navigation, main content areas)
- [ ] Implement metrics dashboard with basic charts (CPU, memory, network)
- [ ] Build trace explorer with list view and basic filtering
- [ ] Add CRD browser with searchable table and YAML viewer

### 4. Kubernetes Deployment
**Goal:** Production-ready Helm chart for cluster deployment
- [ ] Create Helm chart with configurable values for different environments
- [ ] Implement Kubernetes Service (ClusterIP) and Ingress configurations
- [ ] Add ConfigMaps for application configuration and environment variables
- [ ] Configure Secrets management for ClickHouse credentials
- [ ] Setup proper RBAC with ServiceAccounts and role bindings

### 5. Authentication & Security Integration
**Goal:** Basic authentication flow with Kubernetes RBAC
- [ ] Implement JWT token validation for API requests
- [ ] Integrate Kubernetes ServiceAccount-based authorization
- [ ] Add CORS configuration for frontend-backend communication
- [ ] Create basic RBAC rules for cluster metrics and resource access
- [ ] Implement request logging for security auditing

---

## Success Criteria

- [x] **API Functionality**: All four core endpoints return real data from ClickHouse
- [x] **Dashboard Access**: React frontend displays metrics, traces, and CRDs
- [x] **Kubernetes Deployment**: Helm chart deploys successfully to test cluster
- [x] **Data Integration**: End-to-end data flow from ClickHouse to UI works
- [x] **Basic Security**: Authentication and basic RBAC controls access
- [x] **Error Handling**: Graceful error responses and user feedback

---

## Technical Implementation

### API Service Architecture
```
gRPC Services:
├── MetricsService          # GET /metrics - cluster health indicators
├── TracingService          # GET /traces - distributed tracing queries  
├── CRDService             # GET /crds - custom resource browsing
└── QueryService           # POST /query - raw ClickHouse interface
```

### Frontend Component Structure
```
React Components:
├── Dashboard/             # Landing page with key metrics widgets
├── MetricsExplorer/       # Time-series charts and filtering
├── TraceViewer/           # Trace list and waterfall visualization
├── CRDBrowser/            # Resource table and YAML viewer
└── QueryLab/             # SQL editor for advanced users
```

### ClickHouse Query Patterns
```sql
-- Metrics aggregation example
SELECT 
  timestamp,
  cluster_id,
  avg(cpu_usage) as avg_cpu,
  max(memory_usage) as max_memory
FROM cluster_metrics 
WHERE timestamp >= ? AND timestamp <= ?
GROUP BY timestamp, cluster_id
ORDER BY timestamp DESC

-- Trace queries example  
SELECT trace_id, span_id, operation_name, duration
FROM distributed_traces
WHERE service_name = ? AND timestamp >= ?
ORDER BY timestamp DESC
LIMIT 100
```

### Helm Chart Configuration
```yaml
# Basic values.yaml structure
apiServer:
  image: registry/apiserver:latest
  replicas: 2
  resources:
    requests:
      cpu: 100m
      memory: 128Mi

clickhouse:
  host: clickhouse.monitoring.svc.cluster.local
  port: 9000
  database: telemetry

frontend:
  image: registry/frontend:latest
  ingress:
    enabled: true
    host: telemetry.company.com
```

---

## Data Models & Interfaces

### Core gRPC Messages
```protobuf
// Metrics request/response
message GetMetricsRequest {
  string cluster_id = 1;
  google.protobuf.Timestamp start_time = 2;
  google.protobuf.Timestamp end_time = 3;
  repeated string metric_types = 4;
}

message MetricsResponse {
  repeated MetricSeries series = 1;
  map<string, string> metadata = 2;
}

// Tracing request/response  
message GetTracesRequest {
  string service_name = 1;
  google.protobuf.Duration lookback = 2;
  int32 limit = 3;
}
```

### React TypeScript Interfaces
```typescript
interface MetricDataPoint {
  timestamp: Date;
  value: number;
  labels: Record<string, string>;
}

interface TraceSpan {
  traceId: string;
  spanId: string;
  operationName: string;
  startTime: Date;
  duration: number;
  tags: Record<string, any>;
}

interface CustomResource {
  apiVersion: string;
  kind: string;
  metadata: K8sObjectMeta;
  spec: any;
  status?: any;
}
```

---

## Quality & Testing

### Testing Strategy
- **Unit Tests**: 80%+ coverage for all Go packages and React components
- **Integration Tests**: End-to-end API testing with test ClickHouse instance
- **UI Tests**: Component testing with React Testing Library
- **Helm Tests**: Chart validation and deployment verification

### Performance Targets
- **API Response Time**: < 500ms for metrics queries, < 1s for complex traces
- **Frontend Load Time**: < 3s initial load, < 1s subsequent navigation
- **ClickHouse Queries**: Optimized with proper indexes and materialized views
- **Resource Usage**: < 256MB memory, < 200m CPU per replica

---

## Deployment & Operations

### Local Development
```bash
# Start ClickHouse for testing
docker-compose up -d clickhouse

# Run backend
make run-server

# Run frontend  
cd frontend && npm start

# Access dashboard
open http://localhost:3000
```

### Production Deployment
```bash
# Deploy to Kubernetes
helm install telemetry-api ./deployments/helm \
  --set clickhouse.host=prod-clickhouse \
  --set frontend.ingress.host=telemetry.company.com

# Verify deployment
kubectl get pods -l app=telemetry-api
kubectl get ingress telemetry-api
```

---

## Next Phase Preview

Phase 2 (Enhanced MVP) will add:
- Google OIDC authentication integration
- Advanced UI features (time-range picker, auto-refresh, export)
- Query performance optimization and caching
- Enhanced error handling and user feedback
- Monitoring and alerting integration 