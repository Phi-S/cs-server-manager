import { useContext } from "react";
import {
  cancelUpdate,
  restartServer,
  sendCommandWithoutResponse,
  startServer,
  startUpdate,
  State,
  stopServer,
} from "../api/server";
import { WebSocketContext } from "../contexts/WebSocketContext";
import Loading from "./Loading";

export default function ServerControls() {
  const webSocketContext = useContext(WebSocketContext);
  if (webSocketContext === undefined || webSocketContext.status === undefined) {
    return <Loading message="Loading"></Loading>;
  }

  function UpdateOrCancelUpdateButton(state: State) {
    if (state === State.SteamcmdUpdating) {
      return (
        <>
          <button
            onClick={cancelUpdate}
            className="col-3 btn btn-outline-info"
            disabled={state !== State.SteamcmdUpdating}
          >
            Cancel Update
          </button>
        </>
      );
    } else {
      return (
        <>
          <button
            onClick={startUpdate}
            className="col-3 btn btn-outline-info"
            disabled={state !== State.Idle}
          >
            Update
          </button>
        </>
      );
    }
  }

  function StartStopButton(state: State) {
    if (state === State.Idle) {
      return (
        <button onClick={startServer} className="col-3 btn btn-outline-info">
          Start
        </button>
      );
    } else {
      return (
        <button
          onClick={stopServer}
          className="col-3 btn btn-outline-info"
          disabled={state !== State.ServerStarted}
        >
          Stop
        </button>
      );
    }
  }

  const maps = [
    "de_anubis",
    "de_inferno",
    "de_dust2",
    "de_ancient",
    "de_mirage",
    "de_vertigo",
    "de_nuke",
  ];

  async function changeMap(map: string) {
    return sendCommandWithoutResponse(`changelevel ${map}`);
  }

  return (
    <>
      <div className="input-group flex-nowrap w-100 h-100 m-0">
        {StartStopButton(webSocketContext.status.state)}
        <button onClick={restartServer} className="col-3 btn btn-outline-info">
          Restart
        </button>
        {UpdateOrCancelUpdateButton(webSocketContext.status.state)}

        <div className="dropdown col-3">
          <button
            className="btn btn-outline-info dropdown-toggle w-100 h-100"
            data-bs-toggle="dropdown"
            aria-expanded="false"
            disabled={webSocketContext.status.state !== State.ServerStarted}
          >
            Change Map
          </button>
          <ul className="dropdown-menu col-3 text-center w-100">
            {maps.map((m) => (
              <li key={m}>
                <button onClick={() => changeMap(m)} className="dropdown-item">
                  {m}
                </button>
              </li>
            ))}
          </ul>
        </div>
      </div>
    </>
  );
}
