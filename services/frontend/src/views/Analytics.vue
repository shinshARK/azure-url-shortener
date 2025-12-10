<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()
const shortCode = route.params.shortCode
const stats = ref(null)
const loading = ref(true)
const error = ref('')
const countries = ref({}) // { "Country Name": { count: 10, ips: { "1.2.3.4": 5 } } }

const fetchAnalytics = async () => {
  try {
    const response = await fetch(`/api/analytics/${shortCode}`)
    
    if (!response.ok) throw new Error('Failed to fetch data')
    
    const data = await response.json()
    stats.value = data
    
    // Process Geolocation
    if (data.locations) {
        await processGeolocation(data.locations)
    }

  } catch (e) {
    error.value = 'Could not load analytics data.'
    console.error(e)
  } finally {
    loading.value = false
  }
}

const processGeolocation = async (locations) => {
    const tempCountries = {}
    
    // Get unique IPs
    const ips = Object.keys(locations)
    
    // Fetch country for each IP
    // Note: In a real app, we'd batch this or cache it on the backend.
    // Using ip-api.com (free, rate limited)
    for (const ip of ips) {
        try {
            // Skip local IPs for lookup (optional, but good for dev)
            if (ip === '::1' || ip === '127.0.0.1' || ip.startsWith('192.168.')) {
                addIpToCountry(tempCountries, 'Local Network', ip, locations[ip])
                continue
            }

            const res = await fetch(`http://ip-api.com/json/${ip}`)
            const geo = await res.json()
            
            const country = geo.status === 'success' ? geo.country : 'Unknown Location'
            addIpToCountry(tempCountries, country, ip, locations[ip])
            
        } catch (e) {
            addIpToCountry(tempCountries, 'Unknown Location', ip, locations[ip])
        }
    }
    countries.value = tempCountries
}

const addIpToCountry = (store, country, ip, count) => {
    if (!store[country]) {
        store[country] = { total: 0, ips: {} }
    }
    store[country].total += count
    store[country].ips[ip] = count
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

      <!-- Locations (Geo) -->
      <div class="card locations-card">
        <h3>LOCATIONS</h3>
        <div v-if="Object.keys(countries).length === 0">No location data</div>
        <div v-else>
            <div v-for="(data, country) in countries" :key="country" class="country-group">
                <div class="country-header">
                    <span class="country-name">{{ country }}</span>
                    <span class="country-total">{{ data.total }}</span>
                </div>
                <ul class="ip-list">
                    <li v-for="(count, ip) in data.ips" :key="ip" class="ip-item">
                        <span class="ip-addr">{{ ip }}</span>
                        <span class="ip-count">{{ count }}</span>
                    </li>
                </ul>
            </div>
        </div>
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
  background: transparent; /* Transparent background */
}

.total-clicks {
  text-align: center;
  grid-column: 1 / -1; /* Span full width */
  border-color: var(--accent-color);
  border-width: 2px;
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

/* Location Styles */
.country-group {
    margin-bottom: 1rem;
}

.country-header {
    display: flex;
    justify-content: space-between;
    font-weight: bold;
    color: var(--accent-color);
    padding: 0.5rem 0;
    border-bottom: 1px solid var(--accent-color);
}

.ip-list {
    padding-left: 1rem;
}

.ip-item {
    font-size: 0.9rem;
    opacity: 0.8;
    border-bottom: 1px dotted var(--grid-color);
}
</style>
