import "bootstrap"
import "bootstrap-icons/font/bootstrap-icons.css"
import "bootstrap/dist/css/bootstrap.min.css"
import './assets/main.css'

import {createApp} from 'vue'
import App from './App.vue'
import router from './router'
import {Setup as SetupState} from './state'
import {Setup as SetupApi} from "@/api/api";

const backendAddress: string = import.meta.env.VITE_BACKEND_ADDRESS
const useTls: boolean = import.meta.env.VITE_TLS

SetupApi(backendAddress, useTls)
SetupState()

const app = createApp(App)

app.use(router)

app.mount('#app')
