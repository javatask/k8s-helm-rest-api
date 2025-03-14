# Helm API Server

A RESTful API server that exposes Helm operations as HTTP endpoints, providing a service-oriented interface to Helm functionality.

## Features

- **Chart Management Operations**: Install, upgrade, and uninstall Helm charts
- **Release Information**: List releases, get details, history, and status of deployments
- **Repository Management**: Add and manage chart repositories
- **Kubernetes Integration**: Seamless integration with Kubernetes clusters

## API Endpoints

### Health Check
- `GET /api/v1/health` - Check server health

### Chart Operations
- `POST /api/v1/charts/install` - Install a chart
- `PUT /api/v1/charts/upgrade` - Upgrade a chart
- `DELETE /api/v1/charts/uninstall` - Uninstall a chart

### Release Information
- `GET /api/v1/releases` - List all releases
- `GET /api/v1/releases/{name}` - Get details of a specific release
- `GET /api/v1/releases/{name}/history` - Get release revision history
- `GET /api/v1/releases/{name}/status` - Get status of a release

## Request and Response Examples

### Installing a Chart

**Request:**
```json
POST /api/v1/charts/install
{
  "releaseName": "test-nginx",
  "chartName": "nginx-ingress",
  "repoURL": "https://helm.nginx.com/stable",
  "namespace": "test"
  "values": {
    "replicaCount": 2,
    "service": {
      "type": "ClusterIP"
    }
  },
  "wait": true,
  "timeout": 300,
  "createNamespace": true
}
```

**Response:**
```json
{
  "success": true,
  "message": "Chart installed successfully",
  "data": {
    "name": "test-nginx",
    "namespace": "test",
    "version": 1,
    "status": "deployed",
    "lastDeployed": "2025-03-14T12:00:00Z",
    "chart": "nginx",
    "appVersion": "1.25.3"
  }
}
```

### Listing Releases

**Request:**
```
GET /api/v1/releases?namespace=default
```

**Response:**
```json
{
  "success": true,
  "data": [
    {
      "name": "test-nginx",
      "namespace": "test",
      "version": 1,
      "status": "deployed",
      "lastDeployed": "2025-03-14T12:00:00Z",
      "chart": "nginx",
      "appVersion": "1.25.3"
    }
  ]
}
```

## Environment Variables

- `PORT` - Server port (default: 8080)
- `KUBECONFIG` - Path to kubeconfig file (if not using in-cluster config)
- `IN_CLUSTER` - Set to "true" to use in-cluster configuration
- `HELM_DRIVER` - Helm driver (default: secrets)
- `HELM_REGISTRY_URL` - Optional default registry URL

## Building and Running

### Using Go

```bash
# Build the server
go build -o helm-api-server .

# Run the server
./helm-api-server
```

## License

MIT
