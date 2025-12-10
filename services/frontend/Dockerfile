# Stage 1: Build Vue App
FROM node:20-alpine AS builder

WORKDIR /app

# Install dependencies
COPY package*.json ./
RUN npm install

# Copy source and build
COPY . .
RUN npm run build

# Stage 2: Serve with Caddy
FROM caddy:alpine

# Copy built assets from builder
COPY --from=builder /app/dist /srv

# Expose port 80
EXPOSE 80

# Caddy runs automatically
