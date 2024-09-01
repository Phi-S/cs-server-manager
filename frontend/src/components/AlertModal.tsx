import { useContext } from "react";
import { AlertContext } from "../contexts/AlertContext";

export default function AlertModal() {
  const alertCtx = useContext(AlertContext);

  if (
    alertCtx === undefined ||
    alertCtx.alert === undefined ||
    alertCtx.triggerAlert === undefined
  ) {
    return;
  }

  function closeAlert() {
    alertCtx?.triggerAlert(undefined);
  }

  return (
    <>
      <div className="modal d-flex">
        <div className="modal-dialog">
          <div className="modal-content">
            <div className="modal-header">
              <h1 className="modal-title fs-5">{alertCtx?.alert?.heading}</h1>
              <button
                type="button"
                className="btn-close"
                aria-label="Close"
                onClick={closeAlert}
              ></button>
            </div>
            <div className="modal-body">{alertCtx?.alert?.message}</div>
            <div className="modal-footer">
              <button
                className="btn btn-secondary"
                data-bs-dismiss="modal"
                onClick={closeAlert}
              >
                Close
              </button>
            </div>
          </div>
        </div>
      </div>
    </>
  );
}
