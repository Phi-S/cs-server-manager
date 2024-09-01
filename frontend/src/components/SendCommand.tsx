import { useContext, useState } from "react";
import { sendCommandWithoutResponse, State } from "../api/server";
import Loading from "./Loading";
import { DefaultContext } from "../contexts/DefaultContext";

export default function SendCommand() {
  const defaultContext = useContext(DefaultContext);
  if (defaultContext === undefined) {
    return <Loading />;
  }

  const [command, setCommand] = useState("");

  function sendCommandIfServerIsRunning(state: State) {
    if (state !== State.ServerStarted) {
      return;
    }

    if (command.trim().length == 0) {
      return;
    }

    sendCommandWithoutResponse(command);
  }

  function onEnter(key: string, state: State) {
    if (key !== "Enter") {
      return;
    }

    sendCommandIfServerIsRunning(state);
  }

  return (
    <div className="text-nowrap d-flex justify-content-center h-100">
      <div className="input-group justify-content-center d-flex w-100">
        <input
          id="command-input"
          value={command}
          className="input-group-text"
          style={{ width: "70%" }}
          placeholder="Server command"
          onKeyUp={(e) => onEnter(e.key, defaultContext.status.state)}
          onChange={(e) => setCommand(e.target.value)}
          autoFocus
        />
        <button
          onClick={() =>
            sendCommandIfServerIsRunning(defaultContext.status.state)
          }
          disabled={defaultContext.status.state !== State.ServerStarted}
          className="btn btn-outline-info"
          style={{ width: "30%" }}
        >
          Send
        </button>
      </div>
    </div>
  );
}
