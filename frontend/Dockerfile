# Build stage
FROM node:20-alpine AS builder

WORKDIR /app

# Copy package files
COPY package*.json ./

# Install dependencies
RUN npm ci

# Copy source code
COPY . .

# Build the app
RUN npm run build

# Runtime stage
FROM nginxinc/nginx-unprivileged:alpine

# Copy built files from builder
COPY --from=builder /app/dist /usr/share/nginx/html

# Copy nginx configuration
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Copy entrypoint script
COPY docker-entrypoint.sh /docker-entrypoint.sh

# Switch to root temporarily to make script executable and set permissions
USER root
RUN chmod +x /docker-entrypoint.sh && \
    chown 1000:1000 /usr/share/nginx/html/config.js

# Use non-root user
USER 1000:1000

# Expose non-privileged port
EXPOSE 8080

# Use entrypoint to generate config at startup
ENTRYPOINT ["/docker-entrypoint.sh"]
CMD ["nginx", "-g", "daemon off;"]




