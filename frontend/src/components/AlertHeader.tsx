import { useContext } from "react";
import { startUpdate } from "../api/server";
import { DefaultContext } from "../contexts/DefaultContext";

export default function AlertHeader() {
  const defaultContext = useContext(DefaultContext);
  if (defaultContext === undefined) {
    return undefined;
  }

  const alertClasses = "alert alert-danger text-center fs-3";

  if (defaultContext.status.is_game_server_installed === false) {
    return (
      <div className={alertClasses}>
        <div>Game server is not yet installed</div>
        <button className="btn btn-outline-info" onClick={startUpdate}>
          Install now
        </button>
      </div>
    );
  }
  return;
}
