<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const username = ref('')
const password = ref('')
const error = ref('')
const errors = ref({})

const register = async () => {
  error.value = ''
  errors.value = {}
  
  try {
    const response = await fetch('/api/auth/register', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ username: username.value, password: password.value })
    })

    if (!response.ok) {
      const data = await response.json()
      if (data.errors) {
        errors.value = data.errors
        return
      }
      if (data.error) {
        throw new Error(data.error)
      }
      throw new Error('Registration failed')
    }

    // Auto login after register
    const loginResponse = await fetch('/api/auth/login', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ username: username.value, password: password.value })
    })
    
    const data = await loginResponse.json()
    localStorage.setItem('token', data.token)
    router.push('/dashboard')
  } catch (e) {
    error.value = e.message
  }
}
</script>

<template>
  <div class="card">
    <h1>NEW IDENTITY</h1>
    <form @submit.prevent="register">
      <div class="form-group">
        <label>IDENTITY</label>
        <input v-model="username" type="text" placeholder="USERNAME" required />
        <span v-if="errors.Username" class="field-error">{{ errors.Username }}</span>
      </div>
      <div class="form-group">
        <label>PASSPHRASE</label>
        <input v-model="password" type="password" placeholder="PASSWORD" required />
        <span v-if="errors.Password" class="field-error">{{ errors.Password }}</span>
      </div>
      <button type="submit">INITIALIZE</button>
    </form>
    <p v-if="error" class="error">{{ error }}</p>
    <p class="footer">
      ALREADY AUTHENTICATED? <router-link to="/login">ACCESS CONTROL</router-link>
    </p>
  </div>
</template>

<style scoped>
.form-group {
  margin-bottom: 1.5rem;
  text-align: left;
}

label {
  display: block;
  margin-bottom: 0.5rem;
  font-size: 0.8rem;
  letter-spacing: 1px;
}

.error {
  color: #ff0000;
  margin-top: 1rem;
}

.field-error {
  color: #ff0000;
  font-size: 0.8rem;
  margin-top: 0.25rem;
  display: block;
}

.footer {
  margin-top: 2rem;
  font-size: 0.9rem;
}
</style>
