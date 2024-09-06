interface Props {
  title?: string;
  message: string;
  handleConfirmation: () => void;
  handleClose: () => void;
}

export default function ConfirmModal({
  title,
  message,
  handleConfirmation,
  handleClose,
}: Props) {
  if (title === undefined) {
    title = "Are you sure?";
  }

  return (
    <>
      <div className="position-absolute modal d-flex">
        <div className="modal-dialog">
          <div
            className="modal-content m-4"
            style={{ maxWidth: "90vw", width: "500px" }}
          >
            <div className="modal-header">
              <h1 className="modal-title fs-5">{title}</h1>
              <button
                type="button"
                className="btn-close"
                aria-label="Close"
                onClick={handleClose}
              ></button>
            </div>
            <div className="modal-body">{message}</div>
            <div className="modal-footer justify-content-start">
              <button
                className="btn btn-outline-success"
                data-bs-dismiss="modal"
                onClick={() => {
                  handleConfirmation();
                }}
              >
                Yes
              </button>
              <button
                className="btn btn-outline-warning"
                data-bs-dismiss="modal"
                onClick={handleClose}
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
