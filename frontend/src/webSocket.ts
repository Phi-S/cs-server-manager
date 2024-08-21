import {ref} from 'vue';
import {LogEntry, Status} from "@/api/server";
import {ConnectToWebSocket, WebSocketMessage} from "@/api/api";
import {useStatusStore} from "@/stores/status";
import {useLogsStore} from "@/stores/logs";

export var connected = ref<boolean>(false)

export async function ConnectToWebsocket() {
    connected.value = false
    while (connected.value == false) {
        try {
            await initStatusAndLogs()
            SetupWebSocket()
            connected.value = true
            console.log("setup finished")
            break
        } catch (error) {
            console.log(error)
            await new Promise(r => setTimeout(r, 2000));
        }
    }
}

async function initStatusAndLogs() {
    const statusStore = useStatusStore()
    await statusStore.init()
    const logsStore = useLogsStore()
    await logsStore.init()
}

function SetupWebSocket() {
    const socket = ConnectToWebSocket();

    socket.onerror = async (event) => {
        socket.close()
        console.info(`websocket closed with error ${event}. Trying to reconnect`)
        ConnectToWebsocket()
    }

    socket.onclose = (event) => {
        socket.close()
        console.info(`websocket closed ${event}. Trying to reconnect`)
        ConnectToWebsocket()
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