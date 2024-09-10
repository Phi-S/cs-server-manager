import moment from "moment";
import { useContext } from "react";
import { LogEntry } from "../api/server";
import { DefaultContext } from "../contexts/DefaultContext";
import Loading from "./Loading";

export default function Console() {
  const defaultContext = useContext(DefaultContext);
  if (defaultContext === undefined) {
    return <Loading />;
  }

  function timestampString(timestampUtc: string): string {
    const offset = new Date().getTimezoneOffset();
    return moment(timestampUtc).add(offset).format("yyyy.MM.DD HH:mm:ss");
  }

  function getLogBackgroundColor(log: LogEntry): string {
    if (log.log_type === "system_info") {
      return "table-success bg-opacity-75";
    } else if (log.log_type === "system_error") {
      return "table-danger bg-opacity-75";
    } else {
      return "";
    }
  }

  return (
    <>
      <div className="overflow-y-scroll rounded-3 border border-2 h-100">
        <table className="table table-sm table-striped">
          <tbody>
            {defaultContext.logs.map((log) => (
              <tr
                key={log.message + log.timestamp}
                className={`border-bottom ${getLogBackgroundColor(log)}`}
              >
                <td className="ps-2 pe-2 pt-1 text-nowrap border-end">
                  {timestampString(log.timestamp as string)}
                </td>
                <td className="ps-2">{log.message}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </>
  );
}
