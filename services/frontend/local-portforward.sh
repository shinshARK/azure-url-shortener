#!/bin/bash

# Define colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}Starting Port Forwarding to Azure Kubernetes...${NC}"

# 1. Auth Service (Local 8080 -> Remote 8080)
echo -e "Forwarding ${GREEN}Auth Service${NC} (8080 -> 8080)..."
kubectl port-forward deployment/auth-service-deployment 8080:8080 > /dev/null 2>&1 &
PID_AUTH=$!

# 2. Link Service (Local 8081 -> Remote 8080)
# NOTE: Mapping local 8081 to remote 8080 as discovered in logs
echo -e "Forwarding ${GREEN}Link Service${NC} (8081 -> 8080)..."
kubectl port-forward deployment/link-management-deployment 8081:8080 > /dev/null 2>&1 &
PID_LINK=$!

# 3. Analytics Service (Local 3001 -> Remote 3001)
echo -e "Forwarding ${GREEN}Analytics Service${NC} (3001 -> 3001)..."
kubectl port-forward deployment/analytics-query-deployment 3001:3001 > /dev/null 2>&1 &
PID_ANALYTICS=$!

# Function to kill processes when script stops
cleanup() {
    echo -e "\n${BLUE}Stopping all port forwards...${NC}"
    kill $PID_AUTH $PID_LINK $PID_ANALYTICS
    exit
}

# Trap Ctrl+C (SIGINT) to run cleanup
trap cleanup SIGINT

echo -e "${GREEN}All systems go!${NC} Localhost is now connected to Azure."
echo "Press Ctrl+C to stop."

# Keep script running to maintain tunnels
wait