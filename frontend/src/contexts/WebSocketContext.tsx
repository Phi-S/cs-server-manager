import {createContext, useEffect, useState} from "react";
import {getLogs, getStatus, LogEntry, State, Status} from "../api/server";
import Loading from "../components/Loading";
import {ChildrenProps} from "../misc";

export const WebSocketContext = createContext<
  WebSocketContextValue | undefined
>(undefined);

export interface WebSocketContextValue {
  status: Status;
  logs: LogEntry[];
}

export default function WebSocketContextWrapper({children}: ChildrenProps) {
  const [connected, setConnected] = useState(false);
  const [status, setStatus] = useState<Status>({
    state: State.Idle,
    hostname: "cs2 server",
    player_count: 0,
    max_player_count: 0,
    map: "",
    ip: "",
    port: "",
    password: "",
  });

  const [logs, setLogs] = useState<LogEntry[]>([]);

  useEffect(() => {
    const abortController = new AbortController();

    getStatus()
      .then((value) => {
        setStatus(value);
      })
      .catch(() => abortController.abort());

    getLogs(500)
      .then((value) => {
        setLogs(value);
      })
      .catch(() => abortController.abort());

    const socket = new WebSocket(GetWebSocketUrl());
    socket.onopen = () => {
      console.info("websocket connection established");
      setConnected(true);
    };

    socket.onclose = () => {
      socket.close();
      console.info(`websocket closed`);
      setConnected(false);
    };

    socket.onmessage = (event) => {
      const msg = JSON.parse(event.data) as {
        type: string;
        message: unknown;
      };

      if (msg.message === undefined) {
        console.warn(`empty websocket message received ${event.data}`);
        return;
      }

      //console.log(msg.message);

      if (msg.type === "status") {
        setStatus(msg.message as Status);
      } else if (msg.type === "log") {
        setLogs((prevState) => [msg.message as LogEntry, ...prevState]);
      } else {
        console.warn(`unexpected websocket message received ${event.data}`);
      }
    };

    return () => {
      abortController.abort();
      socket.close();
    };
  }, []);

  if (!connected || status === undefined || logs === undefined) {
    return <Loading/>;
  }

  return (
    <>
      <WebSocketContext.Provider value={{status, logs}}>
        {children}
      </WebSocketContext.Provider>
    </>
  );
}

let WEBSOCKET_URL: string | undefined = undefined;

function GetWebSocketUrl() {
  const websocketPath = "/api/v1/ws";

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
