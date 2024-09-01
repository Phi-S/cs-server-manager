import { createContext, useEffect, useState } from "react";
import { getLogs, getStatus, LogEntry, State, Status } from "../api/server";
import Loading from "../components/Loading";
import { ChildrenProps } from "../misc";

export const WebSocketContext = createContext<
  WebSocketContextValue | undefined
>(undefined);

export interface WebSocketContextValue {
  status: Status;
  logs: LogEntry[];
}

export default function WebSocketContextWrapper({ children }: ChildrenProps) {
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
  const [error, setError] = useState<string | undefined>(undefined);

  useEffect(() => {
    getStatus()
      .then((value) => {
        setStatus(value);
      })
      .catch(() => {
        setError("Connection to server lost");
      });

    getLogs(500)
      .then((value) => {
        setLogs(value);
      })
      .catch(() => {
        setError("Connection to server lost");
      });

    const socket = new WebSocket(GetWebSocketUrl());
    socket.onopen = () => {
      console.info("websocket connection established");
      setError(undefined);
      setConnected(true);
    };

    socket.onclose = () => {
      socket.close();
      console.info(`websocket closed`);
      setError("Connection to server lost");
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
      socket.onopen = () => {};
      socket.onmessage = () => {};
      socket.onclose = () => {};
      socket.close();
      setError(undefined);
      setConnected(false);
    };
  }, []);

  if (error !== undefined && error !== "") {
    (async () => {
      await new Promise((f) => setTimeout(f, 5000));
      if (error !== undefined && error !== "") {
        window.location.reload();
      }
    })();
    return (
      <>
        <div className="flex flex-column text-center align-content-center w-100">
          <h1>{error}</h1>
          <h2>Trying to reconnect...</h2>
          <button
            className="btn btn-primary align-self-center fs-2"
            onClick={() => {
              window.location.reload();
            }}
          >
            Reload
          </button>
        </div>
      </>
    );
  }

  if (connected === false || status === undefined || logs === undefined) {
    return <Loading />;
  }

  return (
    <>
      <WebSocketContext.Provider value={{ status, logs }}>
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
