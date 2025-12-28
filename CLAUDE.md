# CLAUDE.md

This file provides guidance to Claude Code when working with the Kratos Admin UI repository.

## Project Overview

Kratos Admin UI is a web-based administration interface for Ory Kratos. It consists of:
- **Backend**: Go-based API server that proxies Kratos Admin API
- **Frontend**: Vue 3 application with TypeScript
- **Helm Chart**: For Kubernetes deployment

## Version Management

**IMPORTANT: NEVER edit `version.txt` or chart versions directly.**

- Versions are managed automatically by **release-please** CI workflow
- The workflow analyzes conventional commit messages to determine version bumps
- When commits are pushed to `main`, release-please creates a PR with version updates
- Merging that PR triggers:
  - Docker image builds (tagged with the new version)
  - Helm chart publishing
  - GitHub release creation

### Conventional Commits

Use conventional commit format for all commits:
- `fix:` - Patch version bump (0.3.0 → 0.3.1)
- `feat:` - Minor version bump (0.3.0 → 0.4.0)
- `feat!:` or `BREAKING CHANGE:` - Major version bump (0.3.0 → 1.0.0)

## Development Workflow

### Building Locally

```bash
# Install dependencies
task install

# Build both backend and frontend
task build

# Run in development mode
task dev

# Build Docker images
task docker:build
```

### Testing Changes

Use docker-compose for local testing with a real Kratos instance:

```bash
docker compose up -d
```

## Schema-Driven UI

The frontend is **schema-driven** and dynamically generates forms based on the Kratos identity schema:

1. Frontend fetches schemas from Kratos Admin API via backend
2. `IdentityModal.vue` parses the schema and generates form fields
3. Field types, validation, and constraints come from the schema definition

**Key files:**
- `frontend/src/components/IdentityModal.vue` - Dynamic form generation
- `frontend/src/views/IdentitiesView.vue` - Identity list and search
- `frontend/src/views/IdentityDetailView.vue` - Identity details

### OIDC-Standard Schema

The UI now supports OIDC-standard identity schemas with:
- `preferred_username` (instead of legacy `username`)
- `groups` array (instead of legacy `role` string)
- Standard OIDC claims: `email`, `name`, `given_name`, `family_name`, `picture`, etc.

The UI maintains backward compatibility with legacy field names for display purposes.

## Architecture

### Backend (Go)

- Exposes REST API at `:8080/api`
- Proxies requests to Kratos Admin API
- Handles authentication via JWT
- Admin password authentication

**Key endpoints:**
- `POST /api/auth/login` - Admin login
- `GET /api/identities` - List identities
- `GET /api/identities/:id` - Get identity details
- `PUT /api/identities/:id` - Update identity
- `DELETE /api/identities/:id` - Delete identity
- `POST /api/identities/:id/password` - Reset password
- `GET /api/schemas` - Get identity schemas

### Frontend (Vue 3 + TypeScript)

- Vite-based build system
- Vue Router for navigation
- Pinia for state management
- TailwindCSS for styling
- Lucide icons

**Key components:**
- `IdentityModal` - Create/edit identity form
- `DataTable` - Reusable data table
- `ConfirmDialog` - Confirmation dialogs
- `StatusBadge` - Status indicators

### Helm Chart

Located in `charts/kratos-admin-ui/`:

- Deploys both backend and frontend as separate pods
- ConfigMaps for configuration
- Secrets for sensitive data
- Ingress support
- Resource limits and requests

**Key values:**
- `backend.config.kratosAdminUrl` - Kratos Admin API URL
- `backend.config.kratosPublicUrl` - Kratos Public API URL (for sessions)
- `auth.password.value` - Admin UI password
- `ingress.enabled` - Enable ingress

## CI/CD Workflows

### release-please.yml
- Triggers on push to `main`
- Creates PR with version bumps based on conventional commits
- Updates `version.txt`, `CHANGELOG.md`, and `charts/kratos-admin-ui/Chart.yaml`

### docker-publish.yml
- Triggers on tag push (`v*.*.*`)
- Builds multi-arch images (amd64, arm64)
- Pushes to GitHub Container Registry (`ghcr.io`)
- Tags: `<version>`, `<version>-<sha>`, `latest`

### chart-release.yml
- Publishes Helm chart to GitHub Pages
- Chart available at: `https://nitrikx.github.io/kratos-admin-ui/`

### update-chart-appversion.yml
- Updates chart `appVersion` when new release is created

## Deployment

### Standalone Deployment

```bash
helm repo add kratos-admin-ui https://nitrikx.github.io/kratos-admin-ui/
helm install kratos-admin-ui kratos-admin-ui/kratos-admin-ui \
  --set backend.config.kratosAdminUrl=http://kratos-admin:4434 \
  --set backend.config.kratosPublicUrl=http://kratos-public:4433 \
  --set auth.password.value=your-secure-password
```

### As Dependency in Umbrella Chart

Used in the SkynetK8S `auth` umbrella chart:

```yaml
dependencies:
  - name: kratos-admin-ui
    version: "0.4.0"
    repository: https://nitrikx.github.io/kratos-admin-ui/
```

## Common Tasks

### Adding a New Feature

1. Create a feature branch
2. Implement the feature
3. Commit with `feat:` prefix
4. Push and create PR
5. After merge, release-please creates version bump PR
6. Merge version bump PR to trigger release

### Fixing a Bug

1. Create a bugfix branch
2. Fix the bug
3. Commit with `fix:` prefix
4. Follow same PR process as features

### Updating Dependencies

Frontend dependencies are in `frontend/package.json`.
Backend dependencies are in `backend/go.mod`.

Run `task install` after updating dependencies.

## Troubleshooting

### Kratos Connection Issues

Check backend logs:
```bash
kubectl logs -l app.kubernetes.io/name=kratos-admin-ui-backend -n auth
```

Verify Kratos URLs are correct in the deployment.

### Schema Not Updating

The UI caches schemas. To force refresh:
1. Restart the backend pods
2. Clear browser cache
3. Check Kratos has the updated schema

### Docker Build Failures

Ensure all dependencies are installed before building:
```bash
cd frontend && npm install
cd ../backend && go mod download
```
