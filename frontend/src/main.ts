import "bootstrap"
import "bootstrap-icons/font/bootstrap-icons.css"
import "bootstrap/dist/css/bootstrap.min.css"
import './assets/main.css'

import { Setup } from "@/state"
import { createApp } from 'vue'
import App from './components/App.vue'

Setup()

createApp(App).mount('#app')
