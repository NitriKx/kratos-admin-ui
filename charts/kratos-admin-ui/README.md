# Kratos Admin UI Helm Chart

A Helm chart for deploying [Kratos Admin UI](https://github.com/Nitrikx/kratos-admin-ui) - a modern administration console for [Ory Kratos](https://www.ory.sh/kratos/) identity management system.

## Prerequisites

- Kubernetes 1.19+
- Helm 3.2+
- A running Ory Kratos instance accessible from within the cluster

## Installation

### Add the Helm Repository

```bash
helm repo add kratos-admin-ui https://nitrikx.github.io/kratos-admin-ui
helm repo update
```

### Install the Chart

```bash
helm install my-release kratos-admin-ui/kratos-admin-ui
```

### Install with Custom Values

```bash
helm install my-release kratos-admin-ui/kratos-admin-ui \
  --set backend.config.kratosAdminUrl=http://kratos:4434 \
  --set backend.config.kratosPublicUrl=http://kratos:4433
```

Or using a values file:

```bash
helm install my-release kratos-admin-ui/kratos-admin-ui -f my-values.yaml
```

## Configuration

### Global Parameters

| Parameter | Description | Default |
|-----------|-------------|---------|
| `nameOverride` | Override the chart name | `""` |
| `fullnameOverride` | Override the full release name | `""` |
| `imagePullSecrets` | Image pull secrets | `[]` |

### Backend Parameters

| Parameter | Description | Default |
|-----------|-------------|---------|
| `backend.image.repository` | Backend image repository | `ghcr.io/nitrikx/kratos-admin-backend` |
| `backend.image.pullPolicy` | Image pull policy | `IfNotPresent` |
| `backend.image.tag` | Image tag (defaults to chart appVersion) | `""` |
| `backend.replicaCount` | Number of backend replicas | `1` |
| `backend.service.type` | Backend service type | `ClusterIP` |
| `backend.service.port` | Backend service port | `8080` |
| `backend.resources` | Backend resource requests/limits | `{}` |
| `backend.nodeSelector` | Node selector for backend pods | `{}` |
| `backend.tolerations` | Tolerations for backend pods | `[]` |
| `backend.affinity` | Affinity rules for backend pods | `{}` |

#### Backend Authentication

| Parameter | Description | Default |
|-----------|-------------|---------|
| `backend.auth.existingSecret` | Name of existing secret containing admin password | `""` |
| `backend.auth.existingSecretKey` | Key in the existing secret | `"password"` |
| `backend.auth.jwtSecret.existingSecret` | Name of existing secret containing JWT secret | `""` |
| `backend.auth.jwtSecret.existingSecretKey` | Key in the existing secret for JWT | `"jwt-secret"` |
| `backend.auth.jwtSecret.value` | JWT secret value (ignored if existingSecret is set) | `"change-me-in-production"` |

> **Note**: If `backend.auth.existingSecret` is empty, a random admin password will be auto-generated and stored in a Kubernetes secret.

#### Backend Configuration

| Parameter | Description | Default |
|-----------|-------------|---------|
| `backend.config.kratosAdminUrl` | Kratos Admin API URL | `"http://kratos:4434"` |
| `backend.config.kratosPublicUrl` | Kratos Public API URL | `"http://kratos:4433"` |
| `backend.config.corsOrigins` | CORS allowed origins (comma-separated, empty allows all) | `""` |

### Frontend Parameters

| Parameter | Description | Default |
|-----------|-------------|---------|
| `frontend.image.repository` | Frontend image repository | `ghcr.io/nitrikx/kratos-admin-frontend` |
| `frontend.image.pullPolicy` | Image pull policy | `IfNotPresent` |
| `frontend.image.tag` | Image tag (defaults to chart appVersion) | `""` |
| `frontend.replicaCount` | Number of frontend replicas | `1` |
| `frontend.service.type` | Frontend service type | `ClusterIP` |
| `frontend.service.port` | Frontend service port | `80` |
| `frontend.resources` | Frontend resource requests/limits | `{}` |
| `frontend.nodeSelector` | Node selector for frontend pods | `{}` |
| `frontend.tolerations` | Tolerations for frontend pods | `[]` |
| `frontend.affinity` | Affinity rules for frontend pods | `{}` |

### Ingress Parameters

| Parameter | Description | Default |
|-----------|-------------|---------|
| `ingress.enabled` | Enable ingress | `false` |
| `ingress.className` | Ingress class name | `""` |
| `ingress.annotations` | Ingress annotations | `{}` |
| `ingress.hosts` | Ingress hosts configuration | See below |
| `ingress.tls` | TLS configuration | `[]` |

Default ingress hosts:

```yaml
hosts:
  - host: kratos-admin.local
    paths:
      - path: /api
        pathType: Prefix
        backend: backend
      - path: /
        pathType: Prefix
        backend: frontend
```

### Service Account Parameters

| Parameter | Description | Default |
|-----------|-------------|---------|
| `serviceAccount.create` | Create a service account | `false` |
| `serviceAccount.annotations` | Service account annotations | `{}` |
| `serviceAccount.name` | Service account name | `""` |

### Pod Parameters

| Parameter | Description | Default |
|-----------|-------------|---------|
| `podAnnotations` | Additional pod annotations | `{}` |

## Security

Both backend and frontend deployments are configured with security best practices:

- Run as non-root user (UID 1000)
- Privilege escalation disabled
- Capabilities dropped
- Backend uses read-only root filesystem

## Creating Secrets

If you prefer to manage secrets yourself instead of using auto-generated values, create them before installing the chart.

### Admin Password Secret

Create a secret containing the admin password:

```bash
kubectl create secret generic kratos-admin-password \
  --from-literal=password='your-secure-admin-password'
```

Then reference it in your values:

```yaml
backend:
  auth:
    existingSecret: "kratos-admin-password"
    existingSecretKey: "password"
```

### JWT Secret

Create a secret containing the JWT signing key:

```bash
kubectl create secret generic kratos-admin-jwt \
  --from-literal=jwt-secret='your-secure-jwt-secret-at-least-32-chars'
```

Then reference it in your values:

```yaml
backend:
  auth:
    jwtSecret:
      existingSecret: "kratos-admin-jwt"
      existingSecretKey: "jwt-secret"
```

### Combined Example

Create both secrets at once:

```bash
# Generate secure random values
ADMIN_PASSWORD=$(openssl rand -base64 24)
JWT_SECRET=$(openssl rand -base64 32)

# Create the secrets
kubectl create secret generic kratos-admin-credentials \
  --from-literal=admin-password="$ADMIN_PASSWORD" \
  --from-literal=jwt-secret="$JWT_SECRET"

# Print the admin password for initial login
echo "Admin password: $ADMIN_PASSWORD"
```

Reference the combined secret:

```yaml
backend:
  auth:
    existingSecret: "kratos-admin-credentials"
    existingSecretKey: "admin-password"
    jwtSecret:
      existingSecret: "kratos-admin-credentials"
      existingSecretKey: "jwt-secret"
```

### Using SealedSecrets or External Secrets

For production environments, consider using:

- [Bitnami Sealed Secrets](https://github.com/bitnami-labs/sealed-secrets)
- [External Secrets Operator](https://external-secrets.io/)

Example with External Secrets Operator:

```yaml
apiVersion: external-secrets.io/v1beta1
kind: ExternalSecret
metadata:
  name: kratos-admin-credentials
spec:
  refreshInterval: 1h
  secretStoreRef:
    name: vault-backend
    kind: SecretStore
  target:
    name: kratos-admin-credentials
  data:
    - secretKey: admin-password
      remoteRef:
        key: secret/kratos-admin
        property: admin-password
    - secretKey: jwt-secret
      remoteRef:
        key: secret/kratos-admin
        property: jwt-secret
```

## Examples

### Basic Installation with Ingress

```yaml
# values.yaml
backend:
  config:
    kratosAdminUrl: "http://kratos:4434"
    kratosPublicUrl: "http://kratos:4433"

ingress:
  enabled: true
  className: nginx
  hosts:
    - host: kratos-admin.example.com
      paths:
        - path: /api
          pathType: Prefix
          backend: backend
        - path: /
          pathType: Prefix
          backend: frontend
  tls:
    - secretName: kratos-admin-tls
      hosts:
        - kratos-admin.example.com
```

### Using Existing Secrets

```yaml
# values.yaml
backend:
  auth:
    existingSecret: "my-admin-secret"
    existingSecretKey: "admin-password"
    jwtSecret:
      existingSecret: "my-jwt-secret"
      existingSecretKey: "secret"
```

### With Resource Limits

```yaml
# values.yaml
backend:
  resources:
    limits:
      cpu: 100m
      memory: 128Mi
    requests:
      cpu: 50m
      memory: 64Mi

frontend:
  resources:
    limits:
      cpu: 100m
      memory: 128Mi
    requests:
      cpu: 50m
      memory: 64Mi
```

## Accessing the UI

After installation, you can access the UI using one of these methods:

### With Ingress Enabled

Navigate to your configured host (e.g., `https://kratos-admin.example.com`)

### Without Ingress (Port Forward)

```bash
# Forward the frontend
kubectl port-forward svc/<release-name>-kratos-admin-ui-frontend 8080:80

# Access at http://localhost:8080
```

### Get Admin Password

If no existing secret was provided, retrieve the auto-generated password:

```bash
kubectl get secret <release-name>-kratos-admin-ui-admin -o jsonpath="{.data.password}" | base64 -d; echo
```

Default username: `admin`

## Uninstalling

```bash
helm uninstall my-release
```

## Source Code

- Application: <https://github.com/Nitrikx/kratos-admin-ui>
- Helm Chart: <https://github.com/Nitrikx/kratos-admin-ui/tree/main/charts/kratos-admin-ui>

## License

MIT
