import "bootstrap"
import "bootstrap-icons/font/bootstrap-icons.css"
import "bootstrap/dist/css/bootstrap.min.css"
import './assets/main.css'

import { state } from "@/state"
import { createApp } from 'vue'
import App from './components/App.vue'

state.init()

createApp(App).mount('#app')
