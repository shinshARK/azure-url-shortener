<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()
const shortCode = route.params.shortCode
const stats = ref(null)
const loading = ref(true)
const error = ref('')

const fetchAnalytics = async () => {
  try {
    // In production, this should point to the Analytics Query Service
    // For local dev, we might need a proxy or direct URL
    const response = await fetch(`/api/analytics/${shortCode}`)
    
    if (!response.ok) throw new Error('Failed to fetch data')
    
    stats.value = await response.json()
  } catch (e) {
    error.value = 'Could not load analytics data.'
    console.error(e)
  } finally {
    loading.value = false
  }
}

const goBack = () => {
  router.push('/')
}

onMounted(fetchAnalytics)
</script>

<template>
  <div class="analytics-page">
    <header>
      <button @click="goBack" class="back-btn">← BACK TO DASHBOARD</button>
      <h1>ANALYTICS: /{{ shortCode }}</h1>
    </header>

    <div v-if="loading" class="loading">DECODING DATA STREAM...</div>
    <div v-else-if="error" class="error">{{ error }}</div>
    
    <div v-else class="stats-grid">
      <!-- Total Clicks -->
      <div class="card total-clicks">
        <h3>TOTAL CLICKS</h3>
        <div class="big-number">{{ stats.totalClicks }}</div>
      </div>

      <!-- Browsers -->
      <div class="card">
        <h3>BROWSERS</h3>
        <ul>
          <li v-for="(count, browser) in stats.browsers" :key="browser">
            <span class="label">{{ browser }}</span>
            <span class="value">{{ count }}</span>
          </li>
        </ul>
      </div>

      <!-- OS -->
      <div class="card">
        <h3>OPERATING SYSTEMS</h3>
        <ul>
          <li v-for="(count, os) in stats.os" :key="os">
            <span class="label">{{ os }}</span>
            <span class="value">{{ count }}</span>
          </li>
        </ul>
      </div>

      <!-- Locations (IPs) -->
      <div class="card">
        <h3>LOCATIONS (IP)</h3>
        <ul>
          <li v-for="(count, ip) in stats.locations" :key="ip">
            <span class="label">{{ ip }}</span>
            <span class="value">{{ count }}</span>
          </li>
        </ul>
      </div>
    </div>
  </div>
</template>

<style scoped>
.analytics-page {
  max-width: 1000px;
  margin: 0 auto;
  text-align: left;
}

header {
  margin-bottom: 2rem;
  border-bottom: 2px solid var(--accent-color);
  padding-bottom: 1rem;
}

.back-btn {
  background: none;
  border: none;
  color: var(--accent-color);
  cursor: pointer;
  font-family: inherit;
  font-size: 1rem;
  padding: 0;
  margin-bottom: 1rem;
}

.back-btn:hover {
  text-decoration: underline;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 1.5rem;
}

.card {
  border: 1px solid var(--grid-color);
  padding: 1.5rem;
  background: rgba(255, 255, 255, 0.5);
}

.total-clicks {
  text-align: center;
  grid-column: 1 / -1; /* Span full width */
}

.big-number {
  font-size: 4rem;
  font-weight: bold;
  color: var(--accent-color);
}

ul {
  list-style: none;
  padding: 0;
}

li {
  display: flex;
  justify-content: space-between;
  padding: 0.5rem 0;
  border-bottom: 1px dashed var(--grid-color);
}

li:last-child {
  border-bottom: none;
}

.label {
  font-weight: bold;
}
</style>
