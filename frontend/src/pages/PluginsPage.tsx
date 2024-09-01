import { ChangeEvent, useContext, useEffect, useState } from "react";
import {
  getPlugins,
  installPlugin,
  PluginResp,
  uninstallPlugin,
} from "../api/plugins";
import { State } from "../api/server";
import Loading from "../components/Loading";
import { DefaultContext } from "../contexts/DefaultContext";
import { openInNewTab } from "../util";

export default function PluginsPage() {
  const defaultContext = useContext(DefaultContext);
  const [selectedVersions, setSelectedVersions] = useState<Map<string, string>>(
    new Map<string, string>()
  );
  const [plugins, setPlugins] = useState<PluginResp[]>();

  useEffect(() => {
    updatePlugins();
  }, []);

  if (plugins === undefined || defaultContext === undefined) {
    return <Loading />;
  }

  function isPluginInstalled(
    plugins: PluginResp[],
    name: string,
    version: string
  ): boolean {
    for (const plugin of plugins) {
      if (plugin.name === name) {
        for (const pluginVersion of plugin.versions) {
          if (pluginVersion.name === version) {
            return pluginVersion.installed;
          }
        }
      }
    }

    return false;
  }

  function updatePlugins() {
    getPlugins().then((value) => {
      setPlugins(value);
      for (const plugin of value) {
        if (plugin.versions.length <= 0) {
          continue;
        }

        setSelectedVersions((v) => v.set(plugin.name, plugin.versions[0].name));
      }
    });
  }

  function getSelectedVersion(pluginName: string): string {
    const version = selectedVersions.get(pluginName);
    if (version === undefined) {
      throw new Error(`selected version for ${pluginName} is undefined`);
    }

    return version;
  }

  function onSelectedVersionChanged(
    pluginName: string,
    event: ChangeEvent<HTMLSelectElement>
  ) {
    setSelectedVersions((s) => s.set(pluginName, event.target.value));
  }

  return (
    <>
      {defaultContext.status.state === State.PluginInstalling && (
        <Loading message="Installing plugin" />
      )}
      {defaultContext.status.state === State.PluginUninstalling && (
        <Loading message="Uninstalling plugin" />
      )}
      <table className="table">
        <tbody>
          {plugins.map((plugin) => (
            <tr key={plugin.url} className="row pb-4">
              <td className="col-2 text-center">
                <a
                  className="link btn"
                  onClick={(e) => {
                    e.preventDefault();
                    openInNewTab(plugin.url);
                  }}
                >
                  {plugin.name}
                </a>
              </td>
              <td className="col-6">{plugin.description}</td>
              <td className="col-2">
                <select
                  className="w-100 form-select"
                  onChange={(e) => onSelectedVersionChanged(plugin.name, e)}
                >
                  {plugin.versions.map((version) => (
                    <option key={version.name} value={version.name}>
                      {version.name}
                    </option>
                  ))}
                </select>
              </td>
              <td className="col-2">
                {isPluginInstalled(
                  plugins,
                  plugin.name,
                  getSelectedVersion(plugin.name)
                ) ? (
                  <button
                    className="btn btn-outline-info w-100"
                    onClick={() =>
                      uninstallPlugin().then((_) => updatePlugins())
                    }
                  >
                    Uninstall
                  </button>
                ) : (
                  <button
                    className="btn btn-outline-info w-100"
                    onClick={() =>
                      installPlugin(
                        plugin.name,
                        getSelectedVersion(plugin.name)
                      ).then((_) => updatePlugins())
                    }
                  >
                    Install
                  </button>
                )}
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </>
  );
}
