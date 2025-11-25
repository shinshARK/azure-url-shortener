# Project: Cloud-Native URL Shortener

This is the monorepo for our Cloud Technology project.

## Architecture

The system follows a cloud-native event-driven architecture using Azure Functions, Azure Service Bus, and Azure Kubernetes Service (AKS).

1.  **User Clicks Link** → Hits **Azure Function (Redirect Service)**.
2.  **Redirect Service** → Reads **Azure SQL** (Fast lookup) → Redirects User (HTTP 301).
3.  **Redirect Service** → Asynchronously pushes an event ("User Clicked") to **Azure Service Bus** (Queue).
4.  **Analytics Service (AKS Pod)** → Wakes up, pulls message from Queue → Writes to **Cosmos DB**.
5.  **Link Management (AKS Pod)** → CRUD API to create links → Writes to **Azure SQL**.

## Services

- **/services/link-management-service**: (AKS) Creates and manages short links.
- **/services/redirect-service**: (Azure Function) Handles the high-performance redirects.
- **/services/analytics-query-service**: (AKS) Serves aggregated analytics data.
- **/services/analytics-processing-service**: (AKS) Asynchronously processes click data from Service Bus.
