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
        <form className="w-100" onSubmit={(e) => e.preventDefault()}>
          <div className="input-group mb-3">
            <span className="input-group-text">Hostname</span>
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

          <div className="input-group mb-3">
            <span className="input-group-text">Server password</span>
            <input
              v-model="settings.password"
              type="text"
              className="form-control"
              defaultValue={settings.password}
              onChange={(event: ChangeEvent<HTMLInputElement>) => {
                setSettings((s) => {
                  s!.password = event.target.value;
                  return s;
                });
              }}
            />
          </div>
          <div className="input-group mb-3">
            <span className="input-group-text">Start map</span>
            <input
              v-model="settings.start_map"
              type="text"
              className="form-control"
              aria-describedby="statMapHelp"
              defaultValue={settings.start_map}
              onChange={(event: ChangeEvent<HTMLInputElement>) => {
                setSettings((s) => {
                  s!.start_map = event.target.value;
                  return s;
                });
              }}
            />
            <small id="StartMapHelp" className="input-group text-muted">
              The map which is loaded when server starts
            </small>
          </div>
          <div className="input-group mb-3">
            <span className="input-group-text">May Players</span>
            <input
              type="number"
              className="form-control"
              defaultValue={settings.max_players}
              onChange={(event: ChangeEvent<HTMLInputElement>) => {
                setSettings((s) => {
                  s!.max_players = Number(event.target.value);
                  return s;
                });
              }}
            />
          </div>
          <div className="input-group mb-3">
            <span className="input-group-text">Steam login token</span>
            <input
              v-model="settings.steam_login_token"
              type="text"
              className="form-control"
              aria-describedby="steamLoginTokenHelp"
              defaultValue={settings.steam_login_token}
              onChange={(event: ChangeEvent<HTMLInputElement>) => {
                setSettings((s) => {
                  s!.steam_login_token = event.target.value;
                  return s;
                });
              }}
            />
            <small id="steamLoginTokenHelp" className="input-group text-muted">
              You can generate a token
              <a
                href="https://steamcommunity.com/dev/managegameservers"
                target="_blank"
              >
                here
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
