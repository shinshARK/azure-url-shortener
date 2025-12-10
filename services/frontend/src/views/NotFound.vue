<script setup>
import { useRouter } from 'vue-router'

const router = useRouter()

const goHome = () => {
  router.push('/')
}
</script>

<template>
  <div class="not-found">
    <div class="content">
      <h1>404</h1>
      <h2>TRANSMISSION LOST</h2>
      <p>The signal you are looking for has faded into the void.</p>
      <button @click="goHome" class="home-btn">RETURN TO BASE</button>
    </div>
  </div>
</template>

<style scoped>
.not-found {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 80vh;
  text-align: center;
  font-family: 'DepartureMono', monospace;
}

.content {
  border: 2px solid var(--accent-color);
  padding: 3rem;
  box-shadow: 0 0 20px rgba(0, 255, 255, 0.1);
}

h1 {
  font-size: 6rem;
  margin: 0;
  color: var(--accent-color);
  text-shadow: 0 0 10px var(--accent-color);
}

h2 {
  font-size: 2rem;
  margin: 1rem 0;
  letter-spacing: 2px;
}

p {
  margin-bottom: 2rem;
  opacity: 0.8;
}

.home-btn {
  background: transparent;
  border: 1px solid var(--accent-color);
  color: var(--accent-color);
  padding: 1rem 2rem;
  font-size: 1.2rem;
  cursor: pointer;
  transition: all 0.3s ease;
  font-family: inherit;
}

.home-btn:hover {
  background: var(--accent-color);
  color: black;
  box-shadow: 0 0 15px var(--accent-color);
}
</style>
