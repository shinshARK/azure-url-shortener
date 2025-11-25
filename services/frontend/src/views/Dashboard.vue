<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const links = ref([])
const originalUrl = ref('')
const customAlias = ref('')
const error = ref('')
const createdLink = ref(null) // To show result to guest

const isLoggedIn = computed(() => !!localStorage.getItem('token'))

const fetchLinks = async () => {
  const token = localStorage.getItem('token')
  if (!token) return // Guests don't have history

  try {
    const response = await fetch('/api/links', {
      headers: { 'Authorization': `Bearer ${token}` }
    })
    if (response.ok) {
      const data = await response.json()
      links.value = data || [] // Handle null response
    }
  } catch (e) {
    console.error(e)
  }
}

const createLink = async () => {
  const token = localStorage.getItem('token')
  const headers = { 'Content-Type': 'application/json' }
  if (token) {
    headers['Authorization'] = `Bearer ${token}`
  }

  try {
    const response = await fetch('/api/links', {
      method: 'POST',
      headers: headers,
      body: JSON.stringify({
        originalUrl: originalUrl.value,
        customAlias: customAlias.value || undefined
      })
    })

    if (!response.ok) throw new Error('Failed to create link')

    const newLink = await response.json()
    
    if (isLoggedIn.value) {
      await fetchLinks()
    } else {
      createdLink.value = newLink // Show single result for guest
    }

    originalUrl.value = ''
    customAlias.value = ''
    error.value = ''
  } catch (e) {
    error.value = 'Failed to create link. Alias might be taken.'
  }
}

const logout = () => {
  localStorage.removeItem('token')
  router.push('/login')
}

const copyToClipboard = async (shortCode) => {
  // Use the custom domain (Root Level)
  const url = `${window.location.origin}/${shortCode}`
  
  try {
    if (navigator.clipboard && window.isSecureContext) {
      await navigator.clipboard.writeText(url)
    } else {
      // Fallback for HTTP/Insecure contexts
      const textArea = document.createElement("textarea")
      textArea.value = url
      textArea.style.position = "fixed"
      textArea.style.left = "-9999px"
      document.body.appendChild(textArea)
      textArea.focus()
      textArea.select()
      
      try {
        document.execCommand('copy')
      } catch (err) {
        console.error('Fallback: Oops, unable to copy', err)
      }
      
      document.body.removeChild(textArea)
    }
    alert('COPIED: ' + url)
  } catch (err) {
    console.error('Failed to copy: ', err)
    alert('Failed to copy to clipboard. Please copy manually: ' + url)
  }
}

onMounted(fetchLinks)
</script>

<template>
  <div class="dashboard">
    <header>
      <h1>LAZURUNE</h1>
      <div v-if="isLoggedIn">
        <button @click="logout" class="nav-btn">LOGOUT</button>
      </div>
      <div v-else>
        <router-link to="/login" class="nav-link">LOGIN</router-link>
        <router-link to="/register" class="nav-link">REGISTER</router-link>
      </div>
    </header>

    <div class="card create-section">
      <h2>NEW TRANSMISSION</h2>
      <form @submit.prevent="createLink" class="create-form">
        <input v-model="originalUrl" placeholder="TARGET URL (https://...)" required />
        <input v-if="isLoggedIn" v-model="customAlias" placeholder="ALIAS (OPTIONAL)" />
        <button type="submit">SHORTEN</button>
      </form>
      <p v-if="error" class="error">{{ error }}</p>
      
      <!-- Guest Result -->
      <div v-if="createdLink && !isLoggedIn" class="guest-result">
        <p>SUCCESS! YOUR LINK IS READY:</p>
        <div class="link-card">
          <span class="alias">/{{ createdLink.shortCode }}</span>
          <button @click="copyToClipboard(createdLink.shortCode)">COPY</button>
        </div>
        <p class="guest-note">Note: Guest links expire in 24 hours. <router-link to="/register">Register</router-link> to keep them forever.</p>
      </div>
    </div>

    <div v-if="isLoggedIn" class="links-section">
      <h2>ACTIVE LINKS</h2>
      <div v-if="links.length === 0" class="empty-state">NO ACTIVE TRANSMISSIONS FOUND</div>
      <div v-else class="grid">
        <div v-for="link in links" :key="link.id" class="link-card">
          <div class="link-info">
            <span class="alias">/{{ link.shortCode }}</span>
            <span class="target">{{ link.originalUrl }}</span>
          </div>
          <button @click="copyToClipboard(link.shortCode)">COPY</button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.dashboard {
  max-width: 800px;
  margin: 0 auto;
  text-align: left;
}

header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 2rem;
  border-bottom: 2px solid var(--accent-color);
  padding-bottom: 1rem;
}

h1 {
  font-size: 2rem;
  margin: 0;
}

.nav-btn, .nav-link {
  font-size: 0.9rem;
  margin-left: 1rem;
  text-decoration: none;
  color: var(--accent-color);
  font-weight: bold;
  cursor: pointer;
}

.nav-link:hover {
  text-decoration: underline;
}

.create-section {
  margin-bottom: 2rem;
}

.create-form {
  display: flex;
  gap: 1rem;
  flex-wrap: wrap;
}

.create-form input {
  flex: 1;
  min-width: 200px;
  margin-bottom: 0;
}

.links-section h2 {
  border-bottom: 1px solid var(--grid-color);
  padding-bottom: 0.5rem;
  margin-bottom: 1rem;
}

.link-card {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1rem;
  border: 1px solid var(--grid-color);
  margin-bottom: 0.5rem;
  background: rgba(255, 255, 255, 0.5);
}

.link-info {
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.alias {
  font-weight: bold;
  font-size: 1.2rem;
  color: var(--accent-color);
}

.target {
  font-size: 0.8rem;
  opacity: 0.7;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  max-width: 400px;
}

.error {
  color: #ff0000;
  margin-top: 0.5rem;
}

.empty-state {
  text-align: center;
  padding: 2rem;
  opacity: 0.5;
  font-style: italic;
}

.guest-result {
  margin-top: 2rem;
  padding-top: 1rem;
  border-top: 1px dashed var(--grid-color);
}

.guest-note {
  font-size: 0.8rem;
  margin-top: 1rem;
  opacity: 0.8;
}
</style>
