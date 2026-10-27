# Build stage
FROM registry.redhat.io/ubi9/nodejs-24:9.8-1793067767 AS builder

# Set working directory
WORKDIR /opt/app-root/src

# Copy package files
COPY package*.json ./

# Install dependencies
RUN npm ci

# Copy source code
COPY . .

# Build the application (standalone version)
RUN npm run build:standalone

# Runtime stage - use UBI9 nginx
FROM registry.redhat.io/ubi9/nginx-124:1-113

# Copy built files to nginx
COPY --from=builder /opt/app-root/src/dist /opt/app-root/src

# Copy nginx configuration
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Run as non-root user (nginx user)
USER 1001

# Expose port
EXPOSE 8080

# Run nginx
CMD ["nginx", "-g", "daemon off;"]
