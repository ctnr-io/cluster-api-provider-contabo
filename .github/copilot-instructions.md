# Copilot Instructions for cluster-api-provider-contabo

## Architecture Overview

This is a **Kubernetes Cluster API Infrastructure Provider** for Contabo VPS services. The provider implements the Cluster API v1beta2 interface to provision and manage Kubernetes clusters on Contabo infrastructure.

### Core Components

- **ContaboCluster** (`api/v1beta2/contabocluster_types.go`): Manages cluster-wide infrastructure (networking, control plane endpoint)
- **ContaboMachine** (`api/v1beta2/contabomachine_types.go`): Manages individual VPS instances with instance reuse pattern
- **ContaboMachineTemplate**: Template for creating machines with consistent configuration

### Key Architecture Patterns

**Instance Reuse Pattern**: The provider implements a sophisticated instance reuse system through display name state management:
```go
// States in contabo_helpers.go
StateAvailable = "capc-available"     // Instance ready for assignment
StateProvisioning = "capc-provisioning" // Instance being configured  
StateClusterBound = "capc-cluster-bound" // Instance successfully assigned
StateError = "capc-error"             // Instance failed provisioning
```

**Condition Handling**: Uses `meta.SetStatusCondition` with `[]metav1.Condition` (NOT the legacy `conditions.MarkFalse/MarkTrue`):
```go
meta.SetStatusCondition(&resource.Status.Conditions, metav1.Condition{
    Type:    infrastructurev1beta2.ReadyCondition,
    Status:  metav1.ConditionFalse,
    Reason:  infrastructurev1beta2.CreatingReason,
    Message: fmt.Sprintf("Error: %s", err.Error()), // Optional
})
```

## Critical Development Workflows

### Code Generation
```bash
make generate          # Generate deepcopy methods and CRDs
make generate-api-client  # Regenerate Contabo API client from OpenAPI spec
```

### Testing
```bash
make test             # Unit tests
make test-e2e        # E2E tests (requires credentials and Kind cluster)
make lint            # golangci-lint
```

### Authentication Setup
**Required environment variables** for development and testing:
```bash
export CONTABO_CLIENT_ID="oauth2-client-id"
export CONTABO_CLIENT_SECRET="oauth2-client-secret"  
export CONTABO_API_USER="contabo-username"
export CONTABO_API_PASSWORD="contabo-password"
```

### Local Development
```bash
make dev-deploy      # Build image, load to Kind, deploy
make dev-redeploy    # Rebuild and redeploy (includes undeploy)
```

## Project-Specific Conventions

### API Client Generation
- OpenAPI spec: `pkg/contabo/openapi.json`
- Generated client: `pkg/contabo/client/` 
- Generated models: `pkg/contabo/models/`
- **Run `make generate-api-client` after updating OpenAPI spec**

### Controller Patterns
- Always initialize conditions: `if resource.Status.Conditions == nil { resource.Status.Conditions = []metav1.Condition{} }`
- Use `meta.SetStatusCondition` for all condition updates
- Include error messages in condition messages for debugging
- Controllers use `patchHelper.Patch(ctx, resource)` pattern for status updates

### Instance Management Patterns
- Instances are reused between cluster lifecycles (not destroyed on cluster deletion)
- Instance state is tracked via display name formatting: `<id>-<state>-<cluster-id>`
- Use helpers in `contabo_helpers.go`: `BuildInstanceDisplayNameWithState()`, `GetInstanceState()`, `ParseProviderID()`

### Naming Conventions
- API group: `infrastructure.cluster.x-k8s.io/v1beta2`
- Provider ID format: `contabo://<instance-id>`
- Finalizers: `contabocluster.infrastructure.cluster.x-k8s.io`, `contabomachine.infrastructure.cluster.x-k8s.io`

### API Version Compatibility
- **ContaboCluster/ContaboMachine**: v1beta2 (this provider)
- **Core CAPI resources**: May use v1beta1 or v1beta2 depending on clusterctl version
- **E2E tests**: Use `Cluster` v1beta1 with `ContaboCluster` v1beta2 for compatibility

## Integration Points

### Contabo API
- OAuth2 authentication with client credentials flow
- Base URL: `https://api.contabo.com`
- Auth endpoint: `https://auth.contabo.com/auth/realms/contabo/protocol/openid-connect/token`
- Rate limiting and retry logic built into client

### Cluster API Integration
- Implements standard Cluster API infrastructure provider interface
- Watches for owner references: Machine -> ContaboMachine, Cluster -> ContaboCluster  
- Uses standard Cluster API conditions and status patterns

### Key Files for Understanding Data Flow
- `cmd/main.go`: Controller setup and credential loading
- `internal/controller/contabocluster_controller.go`: Cluster lifecycle management
- `internal/controller/contabomachine_controller.go`: Instance lifecycle and reuse logic
- `internal/controller/contabo_helpers.go`: Instance state management utilities

## Testing and Debugging

### E2E Test Setup
Run `test/setup-e2e.sh` to validate credentials before running `make test-e2e`

### Common Issues
- **Condition errors**: Ensure using `meta.SetStatusCondition`, not legacy `conditions.Mark*` 
- **API client issues**: Regenerate client with `make generate-api-client` after OpenAPI changes
- **Authentication failures**: Verify all 4 environment variables are set correctly
- **Instance not found**: Check instance display name format and state management in helpers
- **clusterctl version mismatch**: If e2e tests fail with "v1beta2 management clusters, v1beta1 detected", clean existing CAPI installation first with `clusterctl delete --all --include-crd --include-namespace` before running tests