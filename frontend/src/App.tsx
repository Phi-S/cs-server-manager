import { Outlet } from "react-router-dom";
import AlertHeader from "./components/AlertHeader.tsx";
import AlertModal from "./components/AlertModal.tsx";
import NavBar from "./components/NavBar.tsx";
import Server from "./components/Server.tsx";
import AlertContextWrapper from "./contexts/AlertContext.tsx";
import DefaultContextWrapper from "./contexts/DefaultContext.tsx";

export default function App() {
  return (
    <div
      className="m-auto p-3 vh-100 vw-100 overflow-hidden"
      style={{
        maxWidth: "1200px",
      }}
    >
      <AlertContextWrapper>
        <AlertModal />
        <DefaultContextWrapper>
          <div className="d-flex justify-content-center">
            <div
              className="w-100"
              style={{ maxWidth: "700px", height: "50px" }}
            >
              <Server />
            </div>
          </div>

          <div className="mt-3 mb-1" style={{ height: "30px" }}>
            <NavBar />
            <hr className="m-0" />
          </div>

          <AlertHeader />

          <div className="w-100 px-2" style={{ height: "calc(100vh - 120px)" }}>
            <Outlet />
          </div>
        </DefaultContextWrapper>
      </AlertContextWrapper>
    </div>
  );
}
