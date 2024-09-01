import { useContext } from "react";
import { startServer, State, Status, stopServer } from "../api/server";
import { DefaultContext } from "../contexts/DefaultContext";
import { copyToClipboard, navigateTo } from "../util";
import Loading from "./Loading";

export default function Server() {
  const defaultContext = useContext(DefaultContext);
  if (defaultContext === undefined) {
    return <Loading />;
  }

  function getConnectionUrl(status: Status) {
    let connectUrl = `steam://connect/${status.ip}:${status.port}`;
    if (status.password !== "") {
      connectUrl = `${connectUrl}/${status.password}`;
    }
    return connectUrl;
  }

  function getConnectionString(status: Status) {
    let connectString = `connect ${status.ip}:${status.port}`;
    if (status.password !== "") {
      connectString = `${connectString}; password ${status.password}`;
    }
    return connectString;
  }

  function getStartStopSpinner(status: Status) {
    if (status.state === State.ServerStarted) {
      return (
        <>
          <button
            onClick={stopServer}
            className="btn bi-stop p-0 fs-1 black"
          ></button>
        </>
      );
    } else if (status.state !== State.Idle) {
      return (
        <>
          <button className="spinner-grow border-0 align-self-center black">
            <span className="visually-hidden">Loading...</span>
          </button>
        </>
      );
    } else {
      return (
        <>
          <button
            onClick={startServer}
            className="btn bi-play p-0 fs-1 black"
          ></button>
        </>
      );
    }
  }

  function getBackground(status: Status) {
    if (status.state === State.ServerStarted) {
      return "bg-success";
    } else {
      return "bg-warning";
    }
  }

  return (
    <>
      <div
        className={`d-flex flex-row flex-nowrap justify-content-between rounded-2 w-100 h-100 px-2 ${getBackground(
          defaultContext.status
        )}`}
        title={`Current state: ${defaultContext.status.state}`}
      >
        <div
          onClick={() => navigateTo(getConnectionUrl(defaultContext.status!))}
          className="d-flex w-100 btn align-items-center"
          style={{ color: "black" }}
        >
          <span className="col-8 text-start text-truncate fs-3">
            {defaultContext.status.hostname}
          </span>
          <div className="col-4 d-none d-sm-block text-nowrap text-end">
            <span className="pe-2 fs-5">{defaultContext.status.map}</span>[
            <span className="fs-5">
              {defaultContext.status.player_count} /{" "}
              {defaultContext.status.max_player_count}]
            </span>
          </div>
        </div>
        <div></div>

        <div className="d-flex flex-row align-content-center">
          <button
            onClick={() =>
              copyToClipboard(getConnectionString(defaultContext.status!))
            }
            className="btn bi-copy p-0 fs-2 px-2 black"
          ></button>
          {getStartStopSpinner(defaultContext.status)}
        </div>
      </div>
    </>
  );
}
