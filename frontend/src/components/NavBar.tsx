import { NavLink, NavLinkRenderProps } from "react-router-dom";

export default function NavBar() {
  function getActiveClass(isActive: NavLinkRenderProps) {
    const defaultClasses = "btn nav-link px-2 me-2";
    return isActive.isActive
      ? defaultClasses + " active text-decoration-underline"
      : defaultClasses;
  }

  return (
    <>
      <div className="d-flex flex-row">
        <NavLink to={"/"} className={(isActive) => getActiveClass(isActive)}>
          Home
        </NavLink>
        <NavLink
          to={"/settings"}
          className={(isActive) => getActiveClass(isActive)}
        >
          Settings
        </NavLink>
        <NavLink
          to={"/plugins"}
          className={(isActive) => getActiveClass(isActive)}
        >
          Plugins
        </NavLink>
        <NavLink
          to={"/about"}
          className={(isActive) => getActiveClass(isActive)}
        >
          About
        </NavLink>
      </div>
    </>
  );
}
