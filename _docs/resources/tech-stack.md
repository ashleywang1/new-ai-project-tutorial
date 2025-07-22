# Technology Stack Decisions

## Core Components

### Backend
- **Language**: Go (1.21+)
- **Runtime**: Docker container running in Kubernetes
- **Primary Database**: ClickHouse (v23.8+)

### API Layer
- **Protocol**: gRPC (with gRPC-Web for browser compatibility)
- **Libraries**:
  - `google.golang.org/grpc` (core gRPC)
  - `github.com/ClickHouse/clickhouse-go/v2` (ClickHouse client)
  - `github.com/grpc-ecosystem/grpc-gateway` (optional REST JSON transcoding)

### Frontend
- **Framework**: React (v18+)
- **Language**: TypeScript (v5+)
- **gRPC Integration**:
  - `@grpc/grpc-js` (gRPC client)
  - `grpc-web` (for browser compatibility)
- **UI Libraries**:
  - `react-query` (data fetching)
  - `material-ui` (component library)
  - `victory` or `recharts` (visualizations)

## Deployment Architecture
- **Packaging**: Helm (v3+)
- **Service Definition**: 
  - Kubernetes Service (ClusterIP)
  - Ingress (path-based routing to frontend/backend)
- **Configuration**:
  - ConfigMaps for environment variables
  - Secrets for credentials (managed externally)

## Development Tooling
- **Protobuf**: `buf.build` (for schema management)
- **Code Generation**:
  - `protoc-gen-go` (Go gRPC stubs)
  - `protoc-gen-js` (TypeScript types)
- **Build Pipeline**:
  - Multi-stage Docker builds
  - Helm chart linting (`ct` tool)

## Key Implementation Notes
1. All data access flows through ClickHouse (no secondary stores)
2. gRPC service definitions will enforce strict API contracts
3. Frontend will use gRPC-Web protocol via Envoy proxy
4. Helm chart will include both frontend and backend deployments 