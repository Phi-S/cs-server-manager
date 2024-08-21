package plugins

import (
	"cs-server-manager/download"
	"cs-server-manager/download/unzip"
	"cs-server-manager/event"
	"cs-server-manager/gvalidator"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type Plugin struct {
	Name        string    `json:"name" validate:"required,lt=32"`
	Description string    `json:"description" validate:"omitempty,printascii,lt=512"`
	URL         string    `json:"url" validate:"required,url,lt=256"`
	InstallDir  string    `json:"install_dir" validate:"required,dirpath,lt=256"`
	Versions    []Version `json:"versions" validate:"required,dive"`
}

type Version struct {
	Name         string       `json:"name" validate:"required,lt=32"`
	DownloadURL  string       `json:"download_url" validate:"required,url,lt=256"`
	Dependencies []Dependency `json:"dependencies" validate:"omitnil,dive"`
}

type Dependency struct {
	PluginName  string `json:"plugin_name" validate:"required,lt=32"`
	VersionName string `json:"version_name" validate:"required,lt=32"`
}

type InstalledPlugin struct {
	Name           string            `json:"name" validate:"required,lt=32"`
	Version        string            `json:"version" validate:"required,lt=32"`
	InstalledAtUtc time.Time         `json:"installed_at_utc" validate:"required,lt=32"`
	Files          []string          `json:"files" validate:"required"`
	Dependencies   []InstalledPlugin `json:"dependencies" validate:"dive"`
}

type PluginEventsPayload struct {
	Name    string
	Version string
}

type Instance struct {
	running                     atomic.Bool
	lock                        sync.Mutex
	installedPluginJsonFileLock sync.Mutex

	csgoDir                     string
	installedPluginJsonFilePath string

	plugins []Plugin

	onPluginInstallingEvent         event.InstanceWithData[PluginEventsPayload]
	onPluginInstalledEvent          event.InstanceWithData[PluginEventsPayload]
	onPluginInstallationFailedEvent event.InstanceWithData[PluginEventsPayload]
	onPluginUninstallingEvent       event.InstanceWithData[PluginEventsPayload]
	onPluginUninstalledEvent        event.InstanceWithData[PluginEventsPayload]
	onPluginUninstallFailedEvent    event.InstanceWithData[PluginEventsPayload]
}

var (
	//go:embed plugins.json
	defaultPluginsJsonData []byte
)

func New(csgoDir string, pluginsJsonFilePath string, installedPluginJsonPath string) (*Instance, error) {
	if err := gvalidator.Instance().Var(csgoDir, "required,dir"); err != nil {
		return nil, fmt.Errorf("csgoDir '%v' is not valid %w", csgoDir, err)
	}

	if err := gvalidator.Instance().Var(pluginsJsonFilePath, "required,filepath"); err != nil {
		return nil, fmt.Errorf("pluginsJsonFilePath '%v' is not valid %w", pluginsJsonFilePath, err)
	}

	if err := gvalidator.Instance().Var(installedPluginJsonPath, "required,filepath"); err != nil {
		return nil, fmt.Errorf("installedPluginJsonPath '%v' is not valid %w", installedPluginJsonPath, err)
	}

	unmarshalPluginsJson := func(pluginsJsonContent []byte) ([]Plugin, error) {
		var plugins []Plugin
		if err := json.Unmarshal(pluginsJsonContent, &plugins); err != nil {
			return nil, fmt.Errorf("failed to unmashal plugins.json '%v' %w", pluginsJsonFilePath, err)
		}

		if plugins == nil {
			return nil, fmt.Errorf("plugins is nil. failed to get plugins from '%v'", pluginsJsonFilePath)
		}

		if err := gvalidator.Instance().Var(plugins, "dive"); err != nil {
			return nil, fmt.Errorf("validation of plugins.json content failed '%v' %w", pluginsJsonFilePath, err)
		}

		return plugins, nil
	}

	var plugins []Plugin
	if _, err := os.Stat(pluginsJsonFilePath); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("failed to read plugins.json at '%v' %w", pluginsJsonFilePath, err)
		}

		plugins, err = unmarshalPluginsJson(defaultPluginsJsonData)
		if err != nil {
			return nil, fmt.Errorf("unmarshalPluginsJson(defaultPluginsJsonData): %w", err)
		}
	} else {
		pluginsJsonContent, err := os.ReadFile(pluginsJsonFilePath)
		if err != nil {
			return nil, fmt.Errorf("failed to read plugins.json %w", err)
		}

		plugins, err = unmarshalPluginsJson(pluginsJsonContent)
		if err != nil {
			return nil, fmt.Errorf(" unmarshalPluginsJson(pluginsJsonContent): %w", err)
		}
	}

	instance := &Instance{
		csgoDir:                     csgoDir,
		installedPluginJsonFilePath: installedPluginJsonPath,
		plugins:                     plugins,
	}

	if err := instance.writeInstalledPluginsJsonFile(nil); err != nil {
		return nil, fmt.Errorf("writeInstalledPluginsJsonFile(nil): %w", err)
	}

	return instance, nil
}

func (i *Instance) readInstalledPluginsJsonFile() (*InstalledPlugin, error) {
	i.installedPluginJsonFileLock.Lock()
	defer i.installedPluginJsonFileLock.Unlock()

	content, err := os.ReadFile(i.installedPluginJsonFilePath)
	if err != nil {
		return nil, fmt.Errorf("os.ReadFile: %w", err)
	}

	contentJson := string(content)
	if contentJson == "{}" {
		return nil, nil
	}

	var installedPlugin InstalledPlugin
	if err := json.Unmarshal(content, &installedPlugin); err != nil {
		return nil, fmt.Errorf("json.Unmarshal: %w", err)
	}

	if err := gvalidator.Instance().Struct(installedPlugin); err != nil {
		return nil, fmt.Errorf("validate: %w", err)
	}

	return &installedPlugin, nil
}

func (i *Instance) writeInstalledPluginsJsonFile(plugin *InstalledPlugin) error {
	i.installedPluginJsonFileLock.Lock()
	defer i.installedPluginJsonFileLock.Unlock()

	if _, err := os.Stat(i.installedPluginJsonFilePath); err == nil {
		if err := os.Remove(i.installedPluginJsonFilePath); err != nil {
			return fmt.Errorf("remove old installedPluginJsonFile: %w", err)
		}
	} else if errors.Is(err, os.ErrNotExist) == false {
		return fmt.Errorf("os.Stats: %w", err)
	}

	if plugin == nil {
		if err := os.WriteFile(i.installedPluginJsonFilePath, []byte("{}"), os.ModePerm); err != nil {
			return fmt.Errorf("os.WriteFile empty: %w", err)
		}
	} else {
		if err := gvalidator.Instance().Struct(*plugin); err != nil {
			return fmt.Errorf("validation: %w", err)
		}

		jsonContent, err := json.MarshalIndent(*plugin, "", "    ")
		if err != nil {
			return fmt.Errorf("json.MarshalIndent: %w", err)
		}

		if err := os.WriteFile(i.installedPluginJsonFilePath, jsonContent, os.ModePerm); err != nil {
			return fmt.Errorf("os.WriteFile: %w", err)
		}
	}

	return nil
}

func (i *Instance) GetAllAvailablePlugins() []Plugin {
	return i.plugins
}

func (i *Instance) GetInstalledPlugin() (*InstalledPlugin, error) {
	installedPlugin, err := i.readInstalledPluginsJsonFile()
	if err != nil {
		return nil, fmt.Errorf("readInstalledPluginsJsonFile: %w", err)
	}
	return installedPlugin, nil
}

func (i *Instance) OnPluginInstalling(handler func(data event.PayloadWithData[PluginEventsPayload])) {
	i.onPluginInstallingEvent.Register(handler)
}

func (i *Instance) OnPluginInstalled(handler func(data event.PayloadWithData[PluginEventsPayload])) {
	i.onPluginInstalledEvent.Register(handler)
}

func (i *Instance) OnPluginInstallationFailedEvent(handler func(data event.PayloadWithData[PluginEventsPayload])) {
	i.onPluginInstallationFailedEvent.Register(handler)
}

func (i *Instance) OnPluginUninstallingEvent(handler func(data event.PayloadWithData[PluginEventsPayload])) {
	i.onPluginUninstallingEvent.Register(handler)
}

func (i *Instance) OnPluginUninstalledEvent(handler func(data event.PayloadWithData[PluginEventsPayload])) {
	i.onPluginUninstalledEvent.Register(handler)
}

func (i *Instance) OnPluginUninstallFailedEvent(handler func(data event.PayloadWithData[PluginEventsPayload])) {
	i.onPluginUninstallFailedEvent.Register(handler)
}

func (i *Instance) InstallPluginByName(pluginName string, versionName string) error {
	if i.running.Load() {
		return fmt.Errorf("another plugin is currently getting installed/uninstalled")
	}

	i.running.Store(true)
	defer i.running.Store(false)

	i.lock.Lock()
	defer i.lock.Unlock()
	eventPayload := PluginEventsPayload{Name: pluginName, Version: versionName}
	i.onPluginInstallingEvent.Trigger(eventPayload)
	installedPlugin, err := i.installPluginByNameInternal(pluginName, versionName, false)
	if err != nil {
		i.onPluginInstallationFailedEvent.Trigger(eventPayload)
		return err
	}

	if err := i.writeInstalledPluginsJsonFile(&installedPlugin); err != nil {
		if uninstallErr := i.uninstallInternal(installedPlugin); uninstallErr != nil {
			slog.Error("failed to uninstall plugin ", "plugin", pluginName, "error", err)
		}
		return fmt.Errorf("installedPluginJfile.Update: %w", err)
	}

	i.onPluginInstalledEvent.Trigger(eventPayload)
	return nil
}

func (i *Instance) installPluginByNameInternal(pluginName string, versionName string, isDependency bool) (InstalledPlugin, error) {
	plugin, version, err := i.getPluginAndVersionByName(pluginName, versionName)
	if err != nil {
		return InstalledPlugin{}, fmt.Errorf("getPluginAndVersionByName: %w", err)
	}

	installedPlugins, err := i.GetInstalledPlugin()
	if err != nil {
		return InstalledPlugin{}, fmt.Errorf("GetInstalledPlugin: %w", err)
	}

	if isDependency == false && installedPlugins != nil {
		return InstalledPlugin{}, fmt.Errorf("another plugin is already installed '%v'", installedPlugins)
	}

	dependencies := make([]InstalledPlugin, 0)
	for _, dependency := range version.Dependencies {
		installedPlugin, err := i.installPluginByNameInternal(dependency.PluginName, dependency.VersionName, true)
		if err != nil {
			return InstalledPlugin{}, fmt.Errorf("InstallPluginDependency: for: plugin %v | version %v | dependecy: plugin %v | version %v %w",
				plugin.Name, version.Name, dependency.PluginName, dependency.VersionName, err)
		}
		dependencies = append(dependencies, installedPlugin)
	}

	files, err := i.downloadAndInstall(plugin.Name, plugin.InstallDir, version.DownloadURL)
	if err != nil {
		return InstalledPlugin{}, fmt.Errorf("downloadAndInstall: %w", err)
	}

	return InstalledPlugin{
		Name:           plugin.Name,
		Version:        version.Name,
		InstalledAtUtc: time.Now().UTC(),
		Dependencies:   dependencies,
		Files:          files,
	}, nil
}

func (i *Instance) Uninstall(pluginName string) error {
	if i.running.Load() {
		return fmt.Errorf("another plugin is currently getting installed/uninstalled")
	}

	i.running.Store(true)
	defer i.running.Store(false)

	i.lock.Lock()
	defer i.lock.Unlock()

	installedPlugin, err := i.GetInstalledPlugin()
	if err != nil {
		return fmt.Errorf("GetInstalledPlugin: %w", err)
	}
	eventPayload := PluginEventsPayload{Name: pluginName}

	i.onPluginUninstallingEvent.Trigger(eventPayload)
	if err := i.uninstallInternal(*installedPlugin); err != nil {
		i.onPluginUninstallFailedEvent.Trigger(eventPayload)
		return fmt.Errorf("uninstallInternal: %w", err)
	}

	if err := i.writeInstalledPluginsJsonFile(nil); err != nil {
		return fmt.Errorf("installedPluginJfile.Update: %w", err)
	}

	i.onPluginUninstalledEvent.Trigger(eventPayload)
	return nil
}

func (i *Instance) uninstallInternal(plugin InstalledPlugin) error {
	for _, dependency := range plugin.Dependencies {
		if err := i.uninstallInternal(dependency); err != nil {
			return fmt.Errorf("failed to uninstall dependencies '%v' of plugin '%v' %w", dependency.Name, plugin.Name, err)
		}
	}

	for _, file := range plugin.Files {
		if err := os.Remove(file); err != nil {
			return fmt.Errorf("failed to remove file '%v' %w", file, err)
		}
	}

	// additional plugin actions
	if plugin.Name == "metamod_source" {
		gaminfoPaht := filepath.Join(i.csgoDir, "gameinfo.gi")
		if err := metamodUninstall(gaminfoPaht); err != nil {
			// TODO: should return error????
			slog.Error(fmt.Errorf("metamod uninstalled successfully but failed to perfomace additional action after plugin removal '%v' %w", plugin.Name, err).Error())
		}
	}

	return nil
}

func (i *Instance) getPluginAndVersionByName(pluginName string, versionName string) (Plugin, Version, error) {
	pluginFound := false
	for _, plugin := range i.plugins {
		if plugin.Name != pluginName {
			continue
		}
		pluginFound = true

		for _, version := range plugin.Versions {
			if version.Name == versionName {
				return plugin, version, nil
			}
		}
	}

	if pluginFound {
		return Plugin{}, Version{}, fmt.Errorf("plugin found %v but version %v dose not exists in plugin", pluginName, versionName)
	}

	return Plugin{}, Version{}, errors.New("No plugin with name " + pluginName + " found")
}

func (i *Instance) downloadAndInstall(pluginName string, pluginInstallDir, downloadUrl string) ([]string, error) {
	tempDirPath := filepath.Join(i.csgoDir, fmt.Sprintf("temp_%v_%v", pluginName, time.Now().UTC()))
	if err := os.Mkdir(tempDirPath, os.ModePerm); err != nil {
		return nil, fmt.Errorf("os.Mkdir temp dir: %w", err)
	}

	downloadedFilePath, err := download.Download(downloadUrl, tempDirPath)
	if err != nil {
		return nil, fmt.Errorf("download.Download: %w", err)
	}

	destPath := filepath.Join(i.csgoDir, pluginInstallDir)

	unzippedFiles, err := unzipDownload(downloadedFilePath, destPath)
	if err != nil {
		return nil, fmt.Errorf("unzipDownload: %w", err)
	}

	if err := os.RemoveAll(tempDirPath); err != nil {
		return nil, fmt.Errorf("os.RemoveAll temp dir: %w", err)
	}

	// additional plugin actions
	if pluginName == "metamod_source" {
		gaminfoPaht := filepath.Join(i.csgoDir, "gameinfo.gi")
		if err := metamodInstall(gaminfoPaht); err != nil {
			return nil, fmt.Errorf("metamodInstall: %w", err)
		}
	}

	return unzippedFiles, nil
}

func unzipDownload(sourcePath, destDir string) ([]string, error) {
	var unzippedFiles []string
	var err error
	if strings.HasSuffix(sourcePath, ".zip") {
		unzippedFiles, err = unzip.Zip(sourcePath, destDir)
		if err != nil {
			return nil, fmt.Errorf("unzip zip: %w", err)
		}
	} else if strings.HasSuffix(sourcePath, ".tar.gz") {
		if unzippedFiles, err = unzip.TarGz(sourcePath, destDir); err != nil {
			return nil, fmt.Errorf("unzip tar.gz: %w", err)
		}
	} else {
		return nil, fmt.Errorf("downloaded filetype not supported '%v'", sourcePath)
	}

	return unzippedFiles, err
}
