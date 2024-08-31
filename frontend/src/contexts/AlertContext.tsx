import { createContext, useState } from "react";
import { ErrorResponseError } from "../api/api";

export const AlertContext = createContext<IAlertContext | undefined>(undefined);

export interface IAlertContext {
  alert: AlertValue | undefined;
  triggerAlert: (alert: AlertValue | undefined) => void;
}

export interface AlertValue {
  heading: string;
  message: string;
  requestId: string;
}

interface Props {
  children: string | JSX.Element | JSX.Element[];
}

export default function AlertContextWrapper({ children }: Props) {
  const [alert, setAlert] = useState<AlertValue>();

  function triggerAlert(alert: AlertValue | undefined) {
    setAlert(alert);
  }

  window.onunhandledrejection = function (event: PromiseRejectionEvent) {
    const error = event.reason;
    if (error instanceof ErrorResponseError) {
      triggerAlert({
        heading: error.message,
        message: error.errorResponse.message,
        requestId: error.errorResponse.request_id,
      });
    } else {
      triggerAlert({
        heading: "Unexpected error occurred",
        message: error?.message !== undefined ? error.message : "????",
        requestId: "",
      });
    }
  };

  window.onerror = (message, source, lineno, colno, error) => {
    console.log(
      "in global error msg: ",
      message,
      " src: ",
      source,
      " lineno: ",
      lineno,
      " colno: ",
      colno
    );

    if (error instanceof ErrorResponseError) {
      triggerAlert({
        heading: error.message,
        message: error.errorResponse.message,
        requestId: error.errorResponse.request_id,
      });
    } else {
      triggerAlert({
        heading: "Unexpected error occurred",
        message: error?.message !== undefined ? error.message : "????",
        requestId: "",
      });
    }
  };

  return (
    <AlertContext.Provider
      value={{
        alert,
        triggerAlert,
      }}
    >
      {children}
    </AlertContext.Provider>
  );
}
