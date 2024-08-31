import { ChangeEvent, useEffect, useState } from "react";
import { restartServer } from "../api/server";
import { getSettings, Settings, updateSettings } from "../api/settings";
import Loading from "../components/Loading";

export default function SettingsPage() {
  const [settings, setSettings] = useState<Settings | undefined>(undefined);

  useEffect(() => {
    getSettings().then((s) => setSettings(s));
  }, []);

  if (settings === undefined) {
    return <Loading />;
  }

  async function save(restart: boolean) {
    if (settings === undefined) {
      return;
    }

    await updateSettings(settings);
    if (restart) {
      await restartServer();
    }
    getSettings().then((s) => setSettings(s));
  }

  return (
    <>
      <div className="w-100 d-flex justify-content-center">
        <form className="w-100 mx-5" onSubmit={(e) => e.preventDefault()}>
          <div className="form-group mb-3">
            <label htmlFor="hostname">Hostname</label>
            <input
              v-model="settings.hostname"
              type="text"
              className="form-control"
              id="hostname"
              defaultValue={settings.hostname}
              onChange={(event: ChangeEvent<HTMLInputElement>) => {
                setSettings((s) => {
                  s!.hostname = event.target.value;
                  return s;
                });
              }}
            />
          </div>
          <div className="form-group mb-3">
            <label htmlFor="password">Server password</label>
            <input
              v-model="settings.password"
              type="text"
              className="form-control"
              id="password"
              defaultValue={settings.password}
              onChange={(event: ChangeEvent<HTMLInputElement>) => {
                setSettings((s) => {
                  s!.password = event.target.value;
                  return s;
                });
              }}
            />
          </div>
          <div className="form-group mb-3">
            <label htmlFor="startMap">Start map</label>
            <input
              v-model="settings.start_map"
              type="text"
              className="form-control"
              id="startMap"
              aria-describedby="statMapHelp"
              defaultValue={settings.start_map}
              onChange={(event: ChangeEvent<HTMLInputElement>) => {
                setSettings((s) => {
                  s!.start_map = event.target.value;
                  return s;
                });
              }}
            />
            <small id="StartMapHelp" className="form-text text-muted">
              The map which is loaded when server starts
            </small>
          </div>
          <div className="form-group mb-3">
            <label htmlFor="maxPlayers">May Players</label>
            <input
              type="number"
              className="form-control"
              id="maxPlayers"
              defaultValue={settings.max_players}
              onChange={(event: ChangeEvent<HTMLInputElement>) => {
                setSettings((s) => {
                  s!.max_players = Number(event.target.value);
                  return s;
                });
              }}
            />
          </div>
          <div className="form-group mb-3">
            <label htmlFor="steamLoginToken">Steam login token</label>
            <input
              v-model="settings.steam_login_token"
              type="text"
              className="form-control"
              id="steamLoginToken"
              aria-describedby="steamLoginTokenHelp"
              defaultValue={settings.steam_login_token}
              onChange={(event: ChangeEvent<HTMLInputElement>) => {
                setSettings((s) => {
                  s!.steam_login_token = event.target.value;
                  return s;
                });
              }}
            />
            <small id="steamLoginTokenHelp" className="form-text text-muted">
              You can generate a token{" "}
              <a href="https://steamcommunity.com/dev/managegameservers">
                Here
              </a>
            </small>
          </div>
          <button
            onClick={() => save(true)}
            type="submit"
            className="btn btn-outline-info me-2"
          >
            Save and Restart
          </button>
          <button
            onClick={() => save(false)}
            type="submit"
            className="btn btn-outline-info"
          >
            Save
          </button>
          <br />
          <small className="form-text text-muted">
            A server restart is required for those settings to take effect
          </small>
        </form>
      </div>
    </>
  );
}
