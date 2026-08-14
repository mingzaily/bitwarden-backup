import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import './style.css'

const savedTheme = window.localStorage.getItem('vaultsync-theme')
document.documentElement.dataset.theme = savedTheme === 'light' ? 'light' : 'dark'

const app = createApp(App)
app.use(router)
app.mount('#app')
