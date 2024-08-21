import "bootstrap"
import "bootstrap-icons/font/bootstrap-icons.css"
import "bootstrap/dist/css/bootstrap.min.css"
import './assets/main.css'

import {createApp} from 'vue'
import App from './App.vue'
import router from './router'
import {ConnectToWebsocket} from './webSocket'
import {SetupApi} from "@/api/api";
import {createPinia} from "pinia";

const backendAddress: string = import.meta.env.VITE_BACKEND_ADDRESS
const useTls: boolean = import.meta.env.VITE_TLS

const pinia = createPinia()
const app = createApp(App)

app.use(pinia)
app.use(router)

app.mount('#app')

SetupApi(backendAddress, useTls)
ConnectToWebsocket()
