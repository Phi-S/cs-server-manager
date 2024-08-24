import "bootstrap"
import "bootstrap-icons/font/bootstrap-icons.css"
import "bootstrap/dist/css/bootstrap.min.css"
import './assets/main.css'

import {createApp} from 'vue'
import App from './App.vue'
import router from './router'
import {ConnectToWebsocket} from './webSocket'
import {createPinia} from "pinia";
import {SetAPIUrl} from "@/api/api";

const backendAddress: string = import.meta.env.VITE_BACKEND_ADDRESS
const useTls: boolean = import.meta.env.VITE_TLS

const [API_URL, WEBSOCKET_URL] = GetApiAndWebsocketUrl(backendAddress, useTls)
console.log("API_URL: ", API_URL)
console.log("WEBSOCKET_URL: ", WEBSOCKET_URL)

const pinia = createPinia()
const app = createApp(App)

app.use(pinia)
app.use(router)

app.mount('#app')


SetAPIUrl(API_URL)
ConnectToWebsocket(WEBSOCKET_URL)


function GetApiAndWebsocketUrl(backendAddress: string | undefined, useTls: boolean | undefined): [API_URL: string, WEBSOCKET_URL: string] {
    const apiPath = "/api/v1"

    if (backendAddress == null || backendAddress == "" || !isValidHostname(backendAddress)) {
        if (window.location.protocol === "https") {
            return [`https://${window.location.host}${apiPath}`, `wss://${window.location.host}/ws`]
        } else {
            return [`http://${window.location.host}${apiPath}`, `ws://${window.location.host}/ws`]
        }
    } else {
        if (useTls == true) {
            return [`https://${backendAddress}${apiPath}`, `wss://${backendAddress}${apiPath}/ws`]
        } else {
            return [`http://${backendAddress}${apiPath}`, `ws://${backendAddress}${apiPath}/ws`]
        }
    }
}

function isValidHostname(hostname: string): boolean {
    const hostnameRegex = new RegExp("^[a-zA-Z0-9-.]+$")
    const hostnamePortRegex = new RegExp("^[a-zA-Z0-9-.]+:\\d{1,5}$")
    return hostnameRegex.test(hostname) || hostnamePortRegex.test(hostname);
}