import {ref} from 'vue';
import type {LogEntry, Status} from "@/api/server";
import {useStatusStore} from "@/stores/status";
import {useLogsStore} from "@/stores/logs";

export var connected = ref<boolean>(false)

let WEBSOCKET_URL: string

interface WebSocketMessage {
    type: string
    message: any
}

export function ConnectToWebsocket(webSocketUrl: string) {
    WEBSOCKET_URL = webSocketUrl
    connected.value = false
    setInterval(reconnectBackgroundTask, 2000)
}

async function reconnectBackgroundTask() {
    if (connected.value) {
        return
    }

    try {
        await initStatusAndLogs()
        SetupWebSocket(WEBSOCKET_URL)
        console.info("websocket connected established")
        connected.value = true
    } catch (error) {
        console.error(error)
    }
}

async function initStatusAndLogs() {
    const statusStore = useStatusStore()
    await statusStore.init()
    const logsStore = useLogsStore()
    await logsStore.init()
}

function SetupWebSocket(webSocketUrl: string) {
    const socket = new WebSocket(webSocketUrl)

    socket.onerror = async (event) => {
        socket.close()
        console.error(`websocket closed with error ${event}. Trying to reconnect`)
        connected.value = false
    }

    socket.onclose = (event) => {
        socket.close()
        console.error(`websocket closed ${event}. Trying to reconnect`)
        connected.value = false
    }

    socket.onmessage = (event) => {
        const msg = JSON.parse(event.data) as WebSocketMessage

        if (msg.message === undefined) {
            console.log(`empty websocket message received ${event.data}`)
        }

        if (msg.type === "status") {
            useStatusStore().update(msg.message as Status)
        } else if (msg.type === "log") {
            useLogsStore().unshift(msg.message as LogEntry)
        } else {
            console.log(`unexpected websocket message received ${event.data}`)
        }
    }
}