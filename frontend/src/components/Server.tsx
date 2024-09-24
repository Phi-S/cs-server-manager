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
        className={`d-flex flex-row rounded-2 w-100 h-100 px-2 align-items-center justify-items-center ${getBackground(
          defaultContext.status,
        )}`}
        title={`Current state: ${defaultContext.status.state}`}
      >
        <div
          className="text-start text-truncate flex-grow-1 fs-3 black"
          onClick={() => navigateTo(getConnectionUrl(defaultContext.status!))}
        >
          {defaultContext.status.hostname}
        </div>

        <div
          className="d-none d-sm-block text-nowrap text-end black ps-2"
          onClick={() => navigateTo(getConnectionUrl(defaultContext.status!))}
        >
          <span className="pe-2 fs-5">{defaultContext.status.map}</span>[
          <span className="fs-5">
            {defaultContext.status.player_count} /{" "}
            {defaultContext.status.max_player_count}]
          </span>
        </div>
        <div className="px-2 h-100 ">
          <div className="border-2 border-dark border-end black h-100"></div>
        </div>
        <div className="d-flex flex-nowrap black">
          <button
            className="btn bi-copy p-0 pe-1 fs-2 black"
            onClick={() =>
              copyToClipboard(getConnectionString(defaultContext.status!))
            }
          ></button>
          {getStartStopSpinner(defaultContext.status)}
        </div>
      </div>
    </>
  );
}
