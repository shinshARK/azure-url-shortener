<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { use } from 'echarts/core'
import { CanvasRenderer } from 'echarts/renderers'
import { LineChart, PieChart, BarChart } from 'echarts/charts'
import {
  GridComponent,
  TooltipComponent,
  LegendComponent,
  TitleComponent
} from 'echarts/components'
import VChart from 'vue-echarts'

use([
  CanvasRenderer,
  LineChart,
  PieChart,
  BarChart,
  GridComponent,
  TooltipComponent,
  LegendComponent,
  TitleComponent
])

const route = useRoute()
const router = useRouter()
const shortCode = route.params.shortCode
const stats = ref(null)
const loading = ref(true)
const error = ref('')
const countries = ref({}) 

// Chart Options Refs
const timelineOption = ref({})
const browserOption = ref({})
const osOption = ref({})

const fetchAnalytics = async () => {
  try {
    const response = await fetch(`/api/analytics/${shortCode}`)
    
    if (!response.ok) throw new Error('Failed to fetch data')
    
    const data = await response.json()
    stats.value = data
    
    // Process Data for Charts
    processCharts(data)

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

const processCharts = (data) => {
    // 1. Timeline (Clicks per Day)
    const clicksByDate = {}
    data.timeline.forEach(ts => {
        const date = new Date(ts).toLocaleDateString()
        clicksByDate[date] = (clicksByDate[date] || 0) + 1
    })
    
    const dates = Object.keys(clicksByDate).sort((a, b) => new Date(a) - new Date(b))
    const counts = dates.map(d => clicksByDate[d])

    timelineOption.value = {
        backgroundColor: 'transparent',
        tooltip: { 
            trigger: 'axis',
            backgroundColor: 'rgba(0,0,0,0.8)',
            borderColor: '#0044cc',
            textStyle: { color: '#fff', fontFamily: 'Departure Mono' }
        },
        grid: { left: '3%', right: '4%', bottom: '3%', containLabel: true },
        xAxis: {
            type: 'category',
            boundaryGap: false,
            data: dates,
            axisLine: { lineStyle: { color: '#888' } },
            axisLabel: { fontFamily: 'Departure Mono', color: '#666' }
        },
        yAxis: {
            type: 'value',
            splitLine: { lineStyle: { type: 'dashed', color: '#333' } },
            axisLabel: { fontFamily: 'Departure Mono', color: '#666' }
        },
        series: [{
            name: 'Clicks',
            type: 'line',
            data: counts,
            smooth: true,
            symbol: 'circle',
            symbolSize: 8,
            itemStyle: { color: '#0044cc' },
            areaStyle: {
                color: {
                    type: 'linear',
                    x: 0, y: 0, x2: 0, y2: 1,
                    colorStops: [
                        { offset: 0, color: 'rgba(0, 68, 204, 0.5)' },
                        { offset: 1, color: 'rgba(0, 68, 204, 0.0)' }
                    ]
                }
            }
        }]
    }

    // 2. Browsers (Pie)
    const browserData = Object.entries(data.browsers).map(([name, value]) => ({ name, value }))
    browserOption.value = {
        backgroundColor: 'transparent',
        tooltip: { 
            trigger: 'item',
            backgroundColor: 'rgba(0,0,0,0.8)',
            borderColor: '#0044cc',
            textStyle: { color: '#fff', fontFamily: 'Departure Mono' }
        },
        legend: { 
            bottom: '0%', 
            textStyle: { fontFamily: 'Departure Mono', color: '#666' } 
        },
        series: [{
            name: 'Browser',
            type: 'pie',
            radius: ['40%', '70%'],
            avoidLabelOverlap: false,
            itemStyle: {
                borderRadius: 5,
                borderColor: '#fff',
                borderWidth: 1
            },
            label: { show: false },
            data: browserData
        }]
    }

    // 3. OS (Bar)
    const osData = Object.entries(data.os).map(([name, value]) => ({ name, value }))
    osOption.value = {
        backgroundColor: 'transparent',
        tooltip: { 
            trigger: 'axis', 
            axisPointer: { type: 'shadow' },
            backgroundColor: 'rgba(0,0,0,0.8)',
            borderColor: '#0044cc',
            textStyle: { color: '#fff', fontFamily: 'Departure Mono' }
        },
        grid: { left: '3%', right: '4%', bottom: '3%', containLabel: true },
        xAxis: {
            type: 'category',
            data: osData.map(i => i.name),
            axisLabel: { fontFamily: 'Departure Mono', color: '#666' }
        },
        yAxis: {
            type: 'value',
            axisLabel: { fontFamily: 'Departure Mono', color: '#666' },
            splitLine: { lineStyle: { type: 'dashed', color: '#333' } }
        },
        series: [{
            name: 'OS',
            type: 'bar',
            data: osData.map(i => i.value),
            itemStyle: { color: '#4da6ff' }
        }]
    }
}

const processGeolocation = async (locations) => {
    const tempCountries = {}
    const ips = Object.keys(locations)
    
    for (const ip of ips) {
        try {
            if (ip === '::1' || ip === '127.0.0.1' || ip.startsWith('192.168.')) {
                addIpToCountry(tempCountries, 'Local Network', ip, locations[ip])
                continue
            }
            const res = await fetch(`https://ip-api.com/json/${ip}`)
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
      <button @click="goBack" class="back-btn"> BACK TO DASHBOARD</button>
      <h1>ANALYTICS: /{{ shortCode }}</h1>
    </header>

    <div v-if="loading" class="loading">DECODING DATA STREAM...</div>
    <div v-else-if="error" class="error">{{ error }}</div>
    
    <div v-else class="dashboard-grid">
      <!-- Top Row: Key Metrics -->
      <div class="card metric-card">
        <h3>TOTAL CLICKS</h3>
        <div class="big-number">{{ stats.totalClicks }}</div>
      </div>
      
      <div class="card metric-card">
        <h3>UNIQUE VISITORS</h3>
        <div class="big-number">{{ Object.keys(stats.locations).length }}</div>
      </div>

      <!-- Middle Row: Timeline (Full Width) -->
      <div class="card chart-card full-width">
        <h3>TRAFFIC TIMELINE</h3>
        <div class="chart-container">
            <v-chart class="chart" :option="timelineOption" autoresize />
        </div>
      </div>

      <!-- Bottom Row: Distribution -->
      <div class="card chart-card">
        <h3>BROWSERS</h3>
        <div class="chart-container">
            <v-chart class="chart" :option="browserOption" autoresize />
        </div>
      </div>

      <div class="card chart-card">
        <h3>OPERATING SYSTEMS</h3>
        <div class="chart-container">
            <v-chart class="chart" :option="osOption" autoresize />
        </div>
      </div>

      <!-- Locations List -->
      <div class="card full-width">
        <h3>GEOLOCATION DATA</h3>
        <div class="geo-grid">
            <div v-for="(data, country) in countries" :key="country" class="geo-item">
                <div class="country-name">{{ country }}</div>
                <div class="country-count">{{ data.total }} clicks</div>
            </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.analytics-page {
  max-width: 1200px;
  margin: 0 auto;
  padding: 2rem;
}

header {
  margin-bottom: 2rem;
}

.back-btn {
  background: none;
  border: 1px solid var(--text-color);
  color: var(--text-color);
  padding: 0.5rem 1rem;
  cursor: pointer;
  font-family: 'Departure Mono', monospace;
  margin-bottom: 1rem;
}

.back-btn:hover {
  background: var(--text-color);
  color: var(--bg-color);
}

.dashboard-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1.5rem;
}

.full-width {
  grid-column: 1 / -1;
}

.card {
  border: 1px solid var(--grid-color);
  padding: 1.5rem;
  background: rgba(255, 255, 255, 0.05);
}

.metric-card {
    text-align: center;
}

.big-number {
  font-size: 4rem;
  font-weight: bold;
  color: var(--accent-color);
  margin-top: 1rem;
}

.chart-container {
    height: 300px;
    width: 100%;
}

.chart {
    height: 100%;
    width: 100%;
}

.geo-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
    gap: 1rem;
    margin-top: 1rem;
}

.geo-item {
    border: 1px dashed var(--grid-color);
    padding: 1rem;
}

.country-name {
    font-weight: bold;
    margin-bottom: 0.5rem;
}

h3 {
    margin-top: 0;
    border-bottom: 1px solid var(--grid-color);
    padding-bottom: 0.5rem;
    font-size: 1rem;
    text-transform: uppercase;
    letter-spacing: 1px;
}

.loading, .error {
    text-align: center;
    font-size: 1.5rem;
    margin-top: 4rem;
    animation: blink 1s infinite;
}

@keyframes blink {
    50% { opacity: 0.5; }
}
</style>
