import { useEffect, useState } from "react";
import { Outlet } from "react-router-dom";
import Alert from "./components/Alert.tsx";
import NavBar from "./components/NavBar.tsx";
import Server from "./components/Server.tsx";
import AlertContextWrapper from "./contexts/AlertContext.tsx";
import WebSocketContextWrapper from "./contexts/WebSocketContext.tsx";
import { deleteCookie, getCookie, setCookie } from "./util.ts";

export default function App() {
  const navbarCollapseLimit = 992;
  const navbarCollapsedCookieKey = "is_navbar_collapsed";
  const [isNavbarCollapsedByUser, setIsNavbarCollapsedByUser] = useState<
    boolean | undefined
  >(navbarCollapsedCookie());
  const [currentWidth, setCurrentWidth] = useState(window.innerWidth);

  function navbarCollapsedCookie() {
    const cookie = getCookie(navbarCollapsedCookieKey);
    if (cookie !== undefined) {
      if (cookie.toLowerCase() === "true") {
        return true;
      } else if (cookie.toLowerCase() === "false") {
        return false;
      } else {
        deleteCookie(navbarCollapsedCookieKey);
      }
    }

    return undefined;
  }

  useEffect(() => {
    window.addEventListener("resize", () => {
      const width = window.innerWidth;
      setCurrentWidth(width);
    });
  });

  function isNavbarCollapsedFinal(): boolean {
    if (currentWidth < navbarCollapseLimit) {
      if (isNavbarCollapsedByUser === false) {
        return false;
      }

      if (isNavbarCollapsedByUser === true) {
        setIsNavbarCollapsedByUser(undefined);
      }

      return true;
    }

    if (isNavbarCollapsedByUser === true) {
      return true;
    }

    if (isNavbarCollapsedByUser === false) {
      setIsNavbarCollapsedByUser(undefined);
    }
    return false;
  }

  function expandNavbar() {
    setCookie(navbarCollapsedCookieKey, "false", 365);
    setIsNavbarCollapsedByUser(false);
  }

  function collapseNavbar() {
    setCookie(navbarCollapsedCookieKey, "true", 365);
    setIsNavbarCollapsedByUser(true);
  }

  function getNavbarWidth(): string {
    return isNavbarCollapsedFinal() ? "50px" : "150px";
  }

  return (
    <div
      className="m-auto p-3"
      style={{
        height: "100vh",
        width: "100vw",
        maxWidth: "1200px",
        maxHeight: "100vh",
      }}
    >
      <AlertContextWrapper>
        <Alert />
        <WebSocketContextWrapper>
          <div className="d-flex justify-content-center">
            <div
              className="w-100"
              style={{ maxWidth: "700px", height: "50px" }}
            >
              <Server />
            </div>
          </div>
          <hr />
          <div className="d-flex w-100 h-100">
            <div className="pe-2" style={{ width: getNavbarWidth() }}>
              <NavBar
                isCollapsed={isNavbarCollapsedFinal()}
                expandNavbar={expandNavbar}
                collapseNavbar={collapseNavbar}
              />
            </div>
            <div style={{ width: `calc(100% - ${getNavbarWidth()})` }}>
              <Outlet />
            </div>
          </div>
        </WebSocketContextWrapper>
      </AlertContextWrapper>
    </div>
  );
}
