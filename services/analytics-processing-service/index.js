require('dotenv').config();
const { ServiceBusClient } = require("@azure/service-bus");
const { CosmosClient } = require("@azure/cosmos");

// --- Configuration ---
const SERVICE_BUS_CONNECTION_STRING = process.env.SERVICE_BUS_CONNECTION_STRING;
const QUEUE_NAME = "analytics-queue";

const COSMOS_CONNECTION_STRING = process.env.COSMOS_CONNECTION_STRING;
const COSMOS_DATABASE_NAME = process.env.COSMOS_DATABASE_NAME || "analytics-db";
const COSMOS_CONTAINER_NAME = process.env.COSMOS_CONTAINER_NAME || "clicks";

if (!SERVICE_BUS_CONNECTION_STRING || !COSMOS_CONNECTION_STRING) {
    console.error("Error: Missing required environment variables (SERVICE_BUS_CONNECTION_STRING or COSMOS_CONNECTION_STRING).");
    process.exit(1);
}

// --- Clients ---
const sbClient = new ServiceBusClient(SERVICE_BUS_CONNECTION_STRING);
const cosmosClient = new CosmosClient(COSMOS_CONNECTION_STRING);

async function main() {
    console.log("Starting Analytics Processing Service...");

    // 1. Initialize Cosmos DB
    console.log(`Connecting to Cosmos DB: ${COSMOS_DATABASE_NAME} / ${COSMOS_CONTAINER_NAME}`);
    const { database } = await cosmosClient.databases.createIfNotExists({ id: COSMOS_DATABASE_NAME });
    const { container } = await database.containers.createIfNotExists({ 
        id: COSMOS_CONTAINER_NAME,
        partitionKey: "/short_code" // Partition by short_code for efficient queries
    });

    // 2. Initialize Service Bus Receiver
    const receiver = sbClient.createReceiver(QUEUE_NAME);
    console.log(`Listening on Service Bus Queue: ${QUEUE_NAME}`);

    // 3. Message Handler
    const myMessageHandler = async (messageReceived) => {
        try {
            console.log(`Received message: ${messageReceived.body.short_code}`);
            
            const analyticsEvent = messageReceived.body;

            // Add a unique ID for Cosmos DB if not present
            if (!analyticsEvent.id) {
                analyticsEvent.id = `${analyticsEvent.short_code}-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`;
            }

            // Write to Cosmos DB
            await container.items.create(analyticsEvent);
            console.log(`Saved to Cosmos DB: ${analyticsEvent.id}`);

            // Complete the message (remove from queue)
            await receiver.completeMessage(messageReceived);

        } catch (err) {
            console.error("Error processing message:", err);
            // If it's a transient error, we might want to abandon (retry).
            // If it's a data error, we might want to dead-letter.
            // For now, we'll abandon to retry.
            await receiver.abandonMessage(messageReceived);
        }
    };

    // 4. Error Handler
    const myErrorHandler = async (error) => {
        console.error("Service Bus Error:", error);
    };

    // 5. Start Subscription
    receiver.subscribe({
        processMessage: myMessageHandler,
        processError: myErrorHandler
    });

    // Keep the process alive
    console.log("Service is running. Press Ctrl+C to exit.");
}

// Handle graceful shutdown
process.on('SIGINT', async () => {
    console.log("Shutting down...");
    await sbClient.close();
    process.exit(0);
});

main().catch((err) => {
    console.error("Fatal Error:", err);
    process.exit(1);
});
