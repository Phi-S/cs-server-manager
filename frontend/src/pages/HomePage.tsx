import Console from "../components/Console.tsx";
import SendCommand from "../components/SendCommand.tsx";
import ServerControls from "../components/ServerControls.tsx";

export default function HomePage() {
  return (
    <>
      <div className="pb-2" style={{ height: "60px" }}>
        <ServerControls />
      </div>
      <div className="pb-1" style={{ height: "40px" }}>
        <SendCommand />
      </div>
      <div style={{ height: "calc(100% - 100px)" }}>
        <Console />
      </div>
    </>
  );
}
