import { ChangeEvent, useContext, useEffect, useState } from "react";
import {
  getPlugins,
  installPlugin,
  PluginResp,
  uninstallPlugin,
} from "../api/plugins";
import { State } from "../api/server";
import ConfirmModal from "../components/ConfirmModal";
import Loading from "../components/Loading";
import { DefaultContext } from "../contexts/DefaultContext";
import { openInNewTab } from "../util";

export default function PluginsPage() {
  const defaultContext = useContext(DefaultContext);
  const [selectedVersions, setSelectedVersions] = useState<Map<string, string>>(
    new Map<string, string>()
  );
  const [plugins, setPlugins] = useState<PluginResp[]>();
  const [confirm, setConfirm] = useState<
    | { title?: string; message: string; handleConfirmation: () => void }
    | undefined
  >(undefined);

  useEffect(() => {
    updatePlugins();
  }, []);

  if (plugins === undefined || defaultContext === undefined) {
    return <Loading />;
  }

  function isAnyPluginInstalled(plugins: PluginResp[]) {
    for (const plugin of plugins) {
      for (const pluginVersion of plugin.versions) {
        if (pluginVersion.installed === true) {
          return true;
        }
      }
    }

    return false;
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

  function install(plugins: PluginResp[], pluginName: string, version: string) {
    if (isAnyPluginInstalled(plugins)) {
      setConfirm({
        message: `Are you sure you want to uninstall the current plugin and install ${pluginName} (${version})?`,
        handleConfirmation: () => {
          setConfirm(undefined);
          uninstallPlugin().then(() => {
            installPlugin(pluginName, version).then(() => {
              updatePlugins();
            });
          });
        },
      });
    } else {
      installPlugin(pluginName, version).then(() => {
        updatePlugins();
      });
    }
  }

  function updatePlugins() {
    getPlugins().then((value) => {
      value.forEach((plugin) => {
        plugin.versions = plugin.versions.reverse();
      });

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
      {confirm !== undefined && (
        <ConfirmModal
          title={confirm.title}
          message={confirm.message}
          handleConfirmation={confirm.handleConfirmation}
          handleClose={() => setConfirm(undefined)}
        />
      )}
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
                      install(
                        plugins,
                        plugin.name,
                        getSelectedVersion(plugin.name)
                      )
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
