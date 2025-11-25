# Use a lightweight Node.js image
FROM node:18-alpine

# Set working directory
WORKDIR /app

# Copy package files first (for better caching)
COPY package.json ./

# Install dependencies
RUN npm install --production

# Copy the rest of the application code
COPY . .

# Run as a non-root user for security
USER node

# Start the application
CMD ["node", "index.js"]
