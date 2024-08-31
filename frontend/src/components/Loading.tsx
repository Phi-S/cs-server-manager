export default function Loading({ message }: { message?: string }) {
  if (message === undefined) {
    message = "Loading...";
  }
  return (
    <>
      <div className="position-absolute bg-dark w-100 h-100 text-center fs-1 bg-opacity-75 z-3 pt-3 justify-content-center">
        <div className="bi spinner-border"></div>
        <div className="ps-2">{message}</div>
      </div>
    </>
  );
}
