#!/bin/sh
set -e

# Generate runtime configuration from environment variables
# API_URL: The URL of the backend API (empty string means relative URLs/same origin)

cat > /usr/share/nginx/html/config.js << EOF
// Runtime configuration - generated at container startup
// Do not edit manually - values come from environment variables
window.__RUNTIME_CONFIG__ = {
  apiUrl: '${API_URL:-}'
};
EOF

echo "Generated /usr/share/nginx/html/config.js with API_URL=${API_URL:-'(relative URLs)'}"

# Execute the main command (nginx)
exec "$@"
