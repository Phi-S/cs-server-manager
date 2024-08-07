import "bootstrap"
import "bootstrap-icons/font/bootstrap-icons.css"
import "bootstrap/dist/css/bootstrap.min.css"
import './assets/main.css'

import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import { Setup } from './state'

export var API_URL: string
export var WEBSOCKET_URL: string

const envApiHost: string = import.meta.env.VITE_API_HOST
const envTLS: boolean = import.meta.env.VITE_TLS

const apiPath = "/api/v1"

if (envApiHost == undefined) {
    if (window.location.protocol === "https") {
        API_URL = `https://${window.location.host}${apiPath}`
        WEBSOCKET_URL = `wss://${window.location.host}/ws`
    } else {
        API_URL = `http://${window.location.host}${apiPath}`
        WEBSOCKET_URL = `ws://${window.location.host}/ws`
    }
} else {
    if (envTLS == true) {
        API_URL = `https://${envApiHost}${apiPath}`
        WEBSOCKET_URL = `wss://${envApiHost}${apiPath}/ws`
    } else {
        API_URL = `http://${envApiHost}${apiPath}`
        WEBSOCKET_URL = `ws://${envApiHost}${apiPath}/ws`
    }
}

console.log(`API_URL: ${API_URL}`)
console.log(`WEBSOCKET_URL: ${WEBSOCKET_URL}`)

Setup()

const app = createApp(App)

app.use(router)

app.mount('#app')
