<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const username = ref('')
const password = ref('')
const error = ref('')
const errors = ref({})

const login = async () => {
  error.value = ''
  errors.value = {}

  try {
    const response = await fetch('/api/auth/login', {
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
      throw new Error('Login failed')
    }

    const data = await response.json()
    localStorage.setItem('token', data.token)
    localStorage.setItem('username', username.value)
    router.push('/dashboard')
  } catch (e) {
    error.value = e.message
  }
}
</script>

<template>
  <div class="card">
    <h1>ACCESS CONTROL</h1>
    <form @submit.prevent="login">
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
      <button type="submit">AUTHENTICATE</button>
    </form>
    <p v-if="error" class="error">{{ error }}</p>
    <p class="footer">
      NO CREDENTIALS? <router-link to="/register">INITIALIZE IDENTITY</router-link>
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
