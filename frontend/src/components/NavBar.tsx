import { NavLink, NavLinkRenderProps } from "react-router-dom";

export default function NavBar({
  isCollapsed,
  expandNavbar,
  collapseNavbar,
}: {
  isCollapsed: boolean;
  expandNavbar: () => void;
  collapseNavbar: () => void;
}) {
  function getActiveClass(isActive: NavLinkRenderProps) {
    const defaultClasses =
      "btn btn-outline-info nav-link d-flex flex-nowrap w-100 fs-4 mb-2";
    return isActive.isActive ? defaultClasses + " active" : defaultClasses;
  }

  function navLinkText(icon: string, text: string) {
    if (isCollapsed) {
      return <div className={`bi ${icon} col-12`}></div>;
    }
    return (
      <>
        <div className={`bi ${icon} col-2 ps-3`}></div>
        <div className="col-10 text-start ps-4">{text}</div>
      </>
    );
  }

  return (
    <>
      <div className="w-100 d-flex justify-content-end">
        <button
          className={`btn fs-3 bi ${isCollapsed ? "bi-arrow-bar-right" : "bi-arrow-bar-left"}`}
          onClick={isCollapsed ? expandNavbar : collapseNavbar}
        />
      </div>
      <div className="navbar flex-column align-content-start me-2">
        <NavLink to={"/"} className={(isActive) => getActiveClass(isActive)}>
          {navLinkText("bi-house", "Home")}
        </NavLink>
        <NavLink
          to={"/settings"}
          className={(isActive) => getActiveClass(isActive)}
        >
          {navLinkText("bi-gear", "Settings")}
        </NavLink>
        <NavLink
          to={"/plugins"}
          className={(isActive) => getActiveClass(isActive)}
        >
          {navLinkText("bi-plugin", "Plugins")}
        </NavLink>
        <NavLink
          to={"/about"}
          className={(isActive) => getActiveClass(isActive)}
        >
          {navLinkText("bi-info-square", "About")}
        </NavLink>
      </div>
    </>
  );
}
