/*
import { useState } from "react";

export var connected = useState<boolean>(false);

let WEBSOCKET_URL: string | undefined = undefined;
const websocketPath = "/api/v1/ws";

function GetWebSocketUrl() {
  if (WEBSOCKET_URL === undefined) {
    const BACKEND_HOST = import.meta.env.VITE_BACKEND_HOST;
    const BACKEND_USE_TLS = import.meta.env.VITE_BACKEND_USE_TLS;

    if (BACKEND_HOST !== undefined) {
      if (BACKEND_USE_TLS !== undefined && BACKEND_USE_TLS === "true") {
        WEBSOCKET_URL = `wss://${BACKEND_HOST}${websocketPath}`;
      } else {
        WEBSOCKET_URL = `ws://${BACKEND_HOST}${websocketPath}`;
      }
    } else {
      if (window.location.protocol === "https") {
        WEBSOCKET_URL = `wss://${window.location.host}${websocketPath}`;
      } else {
        WEBSOCKET_URL = `ws://${window.location.host}${websocketPath}`;
      }
    }
  }

  return WEBSOCKET_URL;
}

interface WebSocketMessage {
  type: string;
  message: any;
}

export function ConnectToWebsocket(webSocketUrl: string) {
  WEBSOCKET_URL = webSocketUrl;
  connected[1](false);
  setInterval(reconnectBackgroundTask, 2000);
}

async function reconnectBackgroundTask() {
  if (connected[0]) {
    return;
  }

  try {
    await initStatusAndLogs();
    SetupWebSocket(GetWebSocketUrl());
    console.info("websocket connected established");
    connected[1](true);
  } catch (error) {
    console.error(error);
  }
}

async function initStatusAndLogs() {
  const statusStore = useStatusStore();
  await statusStore.init();
  const logsStore = useLogsStore();
  await logsStore.init();
}

function SetupWebSocket(webSocketUrl: string) {
  const socket = new WebSocket(webSocketUrl);

  socket.onerror = async (event) => {
    socket.close();
    console.error(`websocket closed with error ${event}. Trying to reconnect`);
    connected.value = false;
  };

  socket.onclose = (event) => {
    socket.close();
    console.error(`websocket closed ${event}. Trying to reconnect`);
    connected.value = false;
  };

  socket.onmessage = (event) => {
    const msg = JSON.parse(event.data) as WebSocketMessage;

    if (msg.message === undefined) {
      console.log(`empty websocket message received ${event.data}`);
    }

    if (msg.type === "status") {
      useStatusStore().update(msg.message as Status);
    } else if (msg.type === "log") {
      useLogsStore().unshift(msg.message as LogEntry);
    } else {
      console.log(`unexpected websocket message received ${event.data}`);
    }
  };
}
*/