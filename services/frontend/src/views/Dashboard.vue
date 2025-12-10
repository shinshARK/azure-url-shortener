<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import QRCode from 'qrcode'

const router = useRouter()
const links = ref([])
const originalUrl = ref('')
const customAlias = ref('')
const error = ref('')
const createdLink = ref(null)
const username = ref('')

// Edit State
const editingLink = ref(null)
const editUrl = ref('')
const editError = ref('')

// QR State
const qrCodeUrl = ref('')
const showQrModal = ref(false)
const qrLinkShortCode = ref('')

const isLoggedIn = computed(() => !!localStorage.getItem('token'))

const fetchLinks = async () => {
  const token = localStorage.getItem('token')
  if (!token) return

  try {
    const response = await fetch('/api/links', {
      headers: { 'Authorization': `Bearer ${token}` }
    })
    if (response.ok) {
      const data = await response.json()
      links.value = data || []
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

    if (!response.ok) {
        const data = await response.json()
        throw new Error(data.error || 'Failed to create link')
    }

    const newLink = await response.json()
    
    if (isLoggedIn.value) {
      await fetchLinks()
    } else {
      createdLink.value = newLink
    }

    originalUrl.value = ''
    customAlias.value = ''
    error.value = ''
  } catch (e) {
    error.value = e.message
  }
}

const deleteLink = async (shortCode) => {
    if (!confirm('Are you sure you want to delete this link?')) return

    const token = localStorage.getItem('token')
    try {
        const response = await fetch(`/api/links/${shortCode}`, {
            method: 'DELETE',
            headers: { 'Authorization': `Bearer ${token}` }
        })
        if (response.ok) {
            await fetchLinks()
        } else {
            alert('Failed to delete link')
        }
    } catch (e) {
        console.error(e)
    }
}

const startEdit = (link) => {
    editingLink.value = link
    editUrl.value = link.originalUrl
    editError.value = ''
}

const cancelEdit = () => {
    editingLink.value = null
    editUrl.value = ''
    editError.value = ''
}

const saveEdit = async () => {
    const token = localStorage.getItem('token')
    try {
        const response = await fetch(`/api/links/${editingLink.value.shortCode}`, {
            method: 'PUT',
            headers: { 
                'Authorization': `Bearer ${token}`,
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ originalUrl: editUrl.value })
        })

        if (!response.ok) {
            const data = await response.json()
            throw new Error(data.error || 'Failed to update link')
        }

        await fetchLinks()
        cancelEdit()
    } catch (e) {
        editError.value = e.message
    }
}

const logout = () => {
  localStorage.removeItem('token')
  localStorage.removeItem('username')
  router.push('/login')
}

const copyToClipboard = async (shortCode) => {
  const url = `${window.location.origin}/${shortCode}`
  try {
    if (navigator.clipboard && window.isSecureContext) {
      await navigator.clipboard.writeText(url)
    } else {
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

const generateQR = async (shortCode) => {
    const url = `${window.location.origin}/${shortCode}`
    try {
        qrCodeUrl.value = await QRCode.toDataURL(url, { width: 300, margin: 2 })
        qrLinkShortCode.value = shortCode
        showQrModal.value = true
    } catch (err) {
        console.error(err)
        alert('Failed to generate QR code')
    }
}

const closeQrModal = () => {
    showQrModal.value = false
    qrCodeUrl.value = ''
}

const downloadQR = () => {
    const link = document.createElement('a')
    link.download = `qr-${qrLinkShortCode.value}.png`
    link.href = qrCodeUrl.value
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
}

const copyQRImage = async () => {
    try {
        const response = await fetch(qrCodeUrl.value)
        const blob = await response.blob()
        await navigator.clipboard.write([
            new ClipboardItem({ 'image/png': blob })
        ])
        alert('QR Code copied to clipboard!')
    } catch (err) {
        console.error(err)
        alert('Failed to copy image')
    }
}

onMounted(() => {
    username.value = localStorage.getItem('username') || 'USER'
    fetchLinks()
})
</script>

<template>
  <div class="dashboard">
    <header>
      <h1>LAZURUNE</h1>
      <div v-if="isLoggedIn" class="user-info">
        <span>WELCOME, {{ username }}</span>
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
          <div class="actions">
            <button @click="copyToClipboard(createdLink.shortCode)" title="Copy Link">
                <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="9" y="9" width="13" height="13" rx="2" ry="2"></rect><path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"></path></svg>
            </button>
            <button @click="generateQR(createdLink.shortCode)" title="QR Code">
                <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="3" width="7" height="7"></rect><rect x="14" y="3" width="7" height="7"></rect><rect x="14" y="14" width="7" height="7"></rect><path d="M3 14h7v7H3z"></path></svg>
            </button>
          </div>
        </div>
        <p class="guest-note">Note: Guest links expire in 24 hours. <router-link to="/register">Register</router-link> to keep them forever.</p>
      </div>
    </div>

    <div v-if="isLoggedIn" class="links-section">
      <h2>ACTIVE LINKS</h2>
      <div v-if="links.length === 0" class="empty-state">NO ACTIVE TRANSMISSIONS FOUND</div>
      <div v-else class="grid">
        <div v-for="link in links" :key="link.shortCode" class="link-card">
          <div class="link-info">
            <span class="alias">/{{ link.shortCode }}</span>
            <span class="target">{{ link.originalUrl }}</span>
          </div>
          <div class="actions">
            <button @click="router.push(`/analytics/${link.shortCode}`)" class="icon-btn" title="Analytics">
                <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="18" y1="20" x2="18" y2="10"></line><line x1="12" y1="20" x2="12" y2="4"></line><line x1="6" y1="20" x2="6" y2="14"></line></svg>
            </button>
            <button @click="copyToClipboard(link.shortCode)" class="icon-btn" title="Copy Link">
                <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="9" y="9" width="13" height="13" rx="2" ry="2"></rect><path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"></path></svg>
            </button>
            <button @click="generateQR(link.shortCode)" class="icon-btn" title="QR Code">
                <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="3" width="7" height="7"></rect><rect x="14" y="3" width="7" height="7"></rect><rect x="14" y="14" width="7" height="7"></rect><path d="M3 14h7v7H3z"></path></svg>
            </button>
            <button v-if="link.customAlias" @click="startEdit(link)" class="icon-btn" title="Edit Target">
                <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"></path><path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"></path></svg>
            </button>
            <button @click="deleteLink(link.shortCode)" class="icon-btn delete-btn" title="Delete">
                <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="3 6 5 6 21 6"></polyline><path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"></path><line x1="10" y1="11" x2="10" y2="17"></line><line x1="14" y1="11" x2="14" y2="17"></line></svg>
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Edit Modal -->
    <div v-if="editingLink" class="modal-overlay">
        <div class="modal">
            <h3>EDIT TARGET URL</h3>
            <p class="alias-display">Alias: /{{ editingLink.shortCode }}</p>
            <input v-model="editUrl" placeholder="NEW TARGET URL" class="edit-input" />
            <p v-if="editError" class="error">{{ editError }}</p>
            <div class="modal-actions">
                <button @click="saveEdit">SAVE</button>
                <button @click="cancelEdit" class="cancel-btn">CANCEL</button>
            </div>
        </div>
    </div>

    <!-- QR Modal -->
    <div v-if="showQrModal" class="modal-overlay" @click.self="closeQrModal">
        <div class="modal qr-modal">
            <h3>QR CODE</h3>
            <img :src="qrCodeUrl" alt="QR Code" class="qr-image"/>
            <div class="modal-actions centered">
                <button @click="downloadQR">DOWNLOAD</button>
                <button @click="copyQRImage">COPY IMG</button>
                <button @click="closeQrModal" class="cancel-btn">CLOSE</button>
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

.link-card {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1rem;
  border: 2px solid var(--accent-color); /* Thicker and Accent Color */
  margin-bottom: 0.5rem;
  background: transparent;
  transition: box-shadow 0.3s ease;
}

.link-card:hover {
    box-shadow: 0 0 15px rgba(0, 255, 255, 0.1);
}

.actions {
  display: flex;
  gap: 0.5rem;
}

.icon-btn {
  background: transparent;
  border: 1px solid var(--grid-color);
  color: var(--text-color);
  padding: 0.5rem;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 4px;
  transition: all 0.2s ease;
}

.icon-btn:hover {
  border-color: var(--accent-color);
  color: var(--accent-color);
  background: rgba(0, 255, 255, 0.05);
}

.delete-btn {
    border-color: #ff4444;
    color: #ff4444;
}

.delete-btn:hover {
    background: #ff4444;
    color: white;
    border-color: #ff4444;
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

.user-info {
    display: flex;
    align-items: center;
    gap: 1rem;
    font-family: 'DepartureMono', monospace;
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

.link-info {
  display: flex;
  flex-direction: column;
  overflow: hidden;
  margin-right: 1rem;
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
  max-width: 300px;
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

/* Modal Styles */
.modal-overlay {
    position: fixed;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    background: rgba(0, 0, 0, 0.8);
    display: flex;
    justify-content: center;
    align-items: center;
    z-index: 1000;
}

.modal {
    background: #000;
    border: 1px solid var(--accent-color);
    padding: 2rem;
    width: 100%;
    max-width: 500px;
    box-shadow: 0 0 20px rgba(0, 255, 255, 0.2);
    text-align: center;
}

.edit-input {
    width: 100%;
    margin: 1rem 0;
}

.modal-actions {
    display: flex;
    gap: 1rem;
    justify-content: flex-end;
    margin-top: 1rem;
}

.modal-actions.centered {
    justify-content: center;
}

.cancel-btn {
    background: transparent;
    border: 1px solid #666;
    color: #666;
}

.cancel-btn:hover {
    border-color: #fff;
    color: #fff;
}

.qr-image {
    width: 200px;
    height: 200px;
    margin: 1rem auto;
    display: block;
    border: 5px solid white;
}
</style>
