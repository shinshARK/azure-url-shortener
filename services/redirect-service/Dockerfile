# Use a multi-stage build for a smaller, more secure production image
# Build stage
FROM node:18-alpine AS build
WORKDIR /usr/src/app
COPY package*.json ./
RUN npm install
COPY . .
# You can add a build step here if needed (e.g., for TypeScript)
# RUN npm run build

# Production stage
FROM node:18-alpine
WORKDIR /usr/src/app
COPY --from=build /usr/src/app/node_modules ./node_modules
COPY --from=build /usr/src/app/package.json ./
COPY --from=build /usr/src/app/src ./src
EXPOSE 8080
CMD [ "node", "src/index.js" ]
