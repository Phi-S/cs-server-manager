export default function Loading({ message }: { message?: string }) {
  if (message === undefined) {
    message = "Loading...";
  }
  return (
    <>
      <div className="position-fixed z-max bg-dark text-center fs-1 bg-opacity-75 pt-3 justify-content-center h-100 w-100">
        <div className="bi spinner-border"></div>
        <div className="ps-2">{message}</div>
      </div>
    </>
  );
}
