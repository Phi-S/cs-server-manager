import { useContext } from "react";
import {
  cancelUpdate,
  restartServer,
  sendCommandWithoutResponse,
  startServer,
  startUpdate,
  State,
  Status,
  stopServer,
} from "../api/server";
import { DefaultContext } from "../contexts/DefaultContext";
import Loading from "./Loading";

export default function ServerControls() {
  const defaultContext = useContext(DefaultContext);
  if (defaultContext === undefined) {
    return <Loading message="Loading"></Loading>;
  }

  function UpdateOrCancelUpdateButton(status: Status) {
    if (status.state === State.SteamcmdUpdating) {
      return (
        <>
          <button
            onClick={cancelUpdate}
            className="col-3 btn btn-outline-info"
            disabled={status.state !== State.SteamcmdUpdating}
          >
            {status.is_game_server_installed
              ? "Cancel Update"
              : "Cancel install"}
          </button>
        </>
      );
    } else {
      return (
        <>
          <button
            onClick={startUpdate}
            className="col-3 btn btn-outline-info"
            disabled={status.state !== State.Idle}
          >
            {status.is_game_server_installed ? "Update" : "Install"}
          </button>
        </>
      );
    }
  }

  function StartStopButton(status: Status) {
    if (status.state === State.Idle) {
      return (
        <button
          onClick={startServer}
          className="col-3 btn btn-outline-info"
          disabled={status.is_game_server_installed === false}
        >
          Start
        </button>
      );
    } else {
      return (
        <button
          onClick={stopServer}
          className="col-3 btn btn-outline-info"
          disabled={status.state !== State.ServerStarted}
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
        {StartStopButton(defaultContext.status)}
        <button
          onClick={restartServer}
          className="col-3 btn btn-outline-info"
          disabled={
            defaultContext.status.is_game_server_installed === false ||
            (defaultContext.status.state !== State.ServerStarted &&
              defaultContext.status.state !== State.Idle)
          }
        >
          Restart
        </button>
        {UpdateOrCancelUpdateButton(defaultContext.status)}

        <div className="dropdown col-3">
          <button
            className="btn btn-outline-info dropdown-toggle w-100 h-100"
            data-bs-toggle="dropdown"
            aria-expanded="false"
            disabled={defaultContext.status.state !== State.ServerStarted}
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
