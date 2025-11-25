const express = require('express');
const app = express();
const PORT = process.env.PORT || 5000;

app.get('/', (req, res) => {
  res.json({ 
    service: 'link-management-service',
    status: 'running',
    message: 'Hello from Azure URL Shortener!'
  });
});

app.get('/health', (req, res) => {
  res.json({ status: 'healthy' });
});

app.listen(PORT, () => {
  console.log(`Link Management Service running on port ${PORT}`);
});