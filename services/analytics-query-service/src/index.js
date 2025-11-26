require('dotenv').config();
const express = require('express');
const { CosmosClient } = require("@azure/cosmos");
const UAParser = require('ua-parser-js');
const cors = require('cors');
const crypto = require('crypto');

// Polyfill for Cosmos DB SDK on Node < 19
if (!global.crypto) {
    global.crypto = crypto;
}

const app = express();
const PORT = process.env.PORT || 3001;

// --- Configuration ---
const COSMOS_CONNECTION_STRING = process.env.COSMOS_CONNECTION_STRING;
const COSMOS_DATABASE_NAME = process.env.COSMOS_DATABASE_NAME || "analytics-db";
const COSMOS_CONTAINER_NAME = process.env.COSMOS_CONTAINER_NAME || "clicks";

if (!COSMOS_CONNECTION_STRING) {
    console.error("Error: Missing COSMOS_CONNECTION_STRING environment variable.");
    process.exit(1);
}

const cosmosClient = new CosmosClient(COSMOS_CONNECTION_STRING);
const container = cosmosClient.database(COSMOS_DATABASE_NAME).container(COSMOS_CONTAINER_NAME);

app.use(cors());
app.use(express.json());

// --- Helper: Process Analytics Data ---
function processAnalytics(items) {
    const stats = {
        totalClicks: items.length,
        browsers: {},
        os: {},
        locations: {}, // Group by IP
        timeline: []   // Clicks over time
    };

    items.forEach(item => {
        // 1. Parse User Agent
        const parser = new UAParser(item.user_agent);
        const browserName = parser.getBrowser().name || 'Unknown';
        const osName = parser.getOS().name || 'Unknown';

        // 2. Aggregate Browser
        stats.browsers[browserName] = (stats.browsers[browserName] || 0) + 1;

        // 3. Aggregate OS
        stats.os[osName] = (stats.os[osName] || 0) + 1;

        // 4. Aggregate Location (IP)
        const ip = item.ip_address || 'Unknown';
        stats.locations[ip] = (stats.locations[ip] || 0) + 1;

        // 5. Timeline (simplify to date string for now)
        const date = new Date(item.timestamp).toLocaleDateString();
        // We'll group by date in the frontend or here if needed. 
        // For now, let's just send the raw timestamp for the frontend to chart.
        stats.timeline.push(item.timestamp);
    });

    return stats;
}

// --- API Endpoint ---
app.get('/api/analytics/:shortCode', async (req, res) => {
    const { shortCode } = req.params;
    console.log(`Fetching analytics for: ${shortCode}`);

    try {
        // Query Cosmos DB
        // Note: We query by Partition Key (short_code) which is efficient.
        const querySpec = {
            query: "SELECT * FROM c WHERE c.short_code = @shortCode",
            parameters: [
                {
                    name: "@shortCode",
                    value: shortCode
                }
            ]
        };

        const { resources: items } = await container.items.query(querySpec).fetchAll();

        if (items.length === 0) {
            return res.json({
                totalClicks: 0,
                browsers: {},
                os: {},
                locations: {},
                timeline: []
            });
        }

        const analyticsData = processAnalytics(items);
        res.json(analyticsData);

    } catch (error) {
        console.error("Cosmos DB Error:", error);
        res.status(500).json({ error: "Failed to fetch analytics data" });
    }
});

app.get('/health', (req, res) => {
    res.send('Analytics Query Service is Healthy');
});

app.listen(PORT, () => {
    console.log(`Analytics Query Service running on port ${PORT}`);
});
