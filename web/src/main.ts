import './assets/main.css'

import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import { createApi, installApi } from './api'
const pinia = createPinia()
const app = createApp(App)

app.use(pinia)
const api = createApi(import.meta.env.VITE_API_URL || "http://localhost:8080")
installApi(app, api)
app.mount('#app')