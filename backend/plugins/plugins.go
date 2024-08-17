package plugins

import (
	"cs-server-manager/download"
	"cs-server-manager/download/unzip"
	"cs-server-manager/event"
	"cs-server-manager/gvalidator"
	"cs-server-manager/jfile"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type Plugin struct {
	Name        string    `json:"name" validate:"required,lt=32"`
	Description string    `json:"description" validate:"omitempty,printascii,lt=512"`
	URL         string    `json:"url" validate:"required,url,lt=256"`
	InstallDir  string    `json:"install_dir" validate:"required,dirpath,lt=256"`
	Versions    []Version `json:"versions" validate:"dive"`
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
	Name           string    `json:"name" validate:"required,lt=32"`
	Version        string    `json:"version" validate:"required,lt=32"`
	InstalledAtUtc time.Time `json:"installed_at_utc" validate:"required,lt=32"`
}

type OnPluginInstalledEventData struct {
	Name    string
	Version string
}

type Instance struct {
	lock       sync.Mutex
	csGamePath string

	installedPluginJfile *jfile.Instance[[]InstalledPlugin]
	plugins              []Plugin

	onPluginInstalledEvent event.InstanceWithData[OnPluginInstalledEventData]
}

var (
	//go:embed plugins.json
	defaultPluginsJsonData []byte
)

func New(csGamePath string, pluginsJsonFilePath string, installedPluginsJsonPath string) (*Instance, error) {
	if err := gvalidator.Instance().Var(csGamePath, "required,dir"); err != nil {
		return nil, fmt.Errorf("csGamePath %v is not valid %w", csGamePath, err)
	}

	if err := gvalidator.Instance().Var(pluginsJsonFilePath, "required,filepath"); err != nil {
		return nil, fmt.Errorf("pluginsJsonFilePath %v is not valid %w", pluginsJsonFilePath, err)
	}

	if err := gvalidator.Instance().Var(installedPluginsJsonPath, "required,filepath"); err != nil {
		return nil, fmt.Errorf("installedPluginsJsonPath %v is not valid %w", installedPluginsJsonPath, err)
	}

	var plugins []Plugin
	if _, err := os.Stat(pluginsJsonFilePath); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("failed to read plugins.json at '%v' %w", pluginsJsonFilePath, err)
		}

		// create default/embedded plugins.json is none exists
		if err := os.WriteFile(pluginsJsonFilePath, defaultPluginsJsonData, os.ModePerm); err != nil {
			return nil, fmt.Errorf("failed to create default plugins.json file at '%v' %w", pluginsJsonFilePath, err)
		}
	}

	pluginsJson, err := os.ReadFile(pluginsJsonFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read plugins.json %w", err)
	}

	if err := json.Unmarshal(pluginsJson, &plugins); err != nil {
		return nil, fmt.Errorf("failed to unmashal plugins.json '%v' %w", pluginsJsonFilePath, err)
	}

	if plugins == nil {
		return nil, fmt.Errorf("plugins is nil. failed to get plugins from '%v'", pluginsJsonFilePath)
	}

	if err := gvalidator.Instance().Var(plugins, "dive"); err != nil {
		return nil, fmt.Errorf("validation of plugins.json content failed '%v' %w", pluginsJsonFilePath, err)
	}

	installedPluginsJfile, err := jfile.New[[]InstalledPlugin](installedPluginsJsonPath, make([]InstalledPlugin, 0))
	if err != nil {
		return nil, fmt.Errorf("installedPluginsJson jfile.New: %w", err)
	}

	instance := &Instance{
		csGamePath:           csGamePath,
		installedPluginJfile: installedPluginsJfile,
		plugins:              plugins,
	}

	return instance, nil
}

func (i *Instance) GetAllAvailablePlugins() []Plugin {
	return i.plugins
}

func (i *Instance) GetInstalledPlugins() ([]InstalledPlugin, error) {
	installedPlugins, err := i.installedPluginJfile.Read()
	if err != nil {
		return nil, fmt.Errorf("installedPluginJfile: %w", err)
	}
	return *installedPlugins, nil
}

func (i *Instance) InstallPluginByName(pluginName string, versionName string) error {
	i.lock.Lock()
	defer i.lock.Unlock()
	if err := i.installPluginByNameInternal(pluginName, versionName); err != nil {
		return err
	}

	i.onPluginInstalledEvent.Trigger(OnPluginInstalledEventData{Name: pluginName, Version: versionName})
	return nil
}

func (i *Instance) OnPluginInstalled(handler func(data event.PayloadWithData[OnPluginInstalledEventData])) {
	i.onPluginInstalledEvent.Register(handler)
}

func (i *Instance) installPluginByNameInternal(pluginName string, versionName string) error {
	if err := gvalidator.Instance().Var(pluginName, "required,lt=64"); err != nil {
		return fmt.Errorf("pluginName %v is not valid %w", pluginName, err)
	}

	if err := gvalidator.Instance().Var(versionName, "required,lt=64"); err != nil {
		return fmt.Errorf("versionName %v is not valid %w", versionName, err)
	}

	plugin, version, err := i.getPluginAndVersionByName(pluginName, versionName)
	if err != nil {
		return fmt.Errorf("getPluginAndVersionByName: %w", err)
	}

	for _, dependency := range version.Dependencies {
		err := i.installPluginByNameInternal(dependency.PluginName, dependency.VersionName)
		if err != nil {
			return fmt.Errorf("InstallPluginDependency: plugin %v | version %v | dependecy plugin %v | dependency version %v %w",
				plugin.Name, version.Name, dependency.PluginName, dependency.VersionName, err)
		}
	}

	err = i.installedPluginJfile.Update(func(currentData *[]InstalledPlugin) {
		installedPlugin := InstalledPlugin{
			Name:           plugin.Name,
			Version:        version.Name,
			InstalledAtUtc: time.Now().UTC(),
		}
		if currentData == nil {
			*currentData = []InstalledPlugin{installedPlugin}
		} else {
			*currentData = append(*currentData, installedPlugin)
		}
	})
	if err != nil {
		return fmt.Errorf("installedPluginJfile.Update: %w", err)
	}

	if downloadAndUnzipErr := i.downloadAndUnzip(plugin, version.DownloadURL); downloadAndUnzipErr != nil {
		err := i.installedPluginJfile.Update(func(currentData *[]InstalledPlugin) {
			for i, installedPlugin := range *currentData {
				if installedPlugin.Name == plugin.Name && installedPlugin.Version == version.Name {
					*currentData = append((*currentData)[:i], (*currentData)[i+1:]...)
					return
				}
			}
		})
		if err != nil {
			return fmt.Errorf("revert installedPluginJfile update: %w", err)
		}
		return fmt.Errorf("downloadAndUnzip: %w", downloadAndUnzipErr)
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

func (i *Instance) downloadAndUnzip(plugin Plugin, downloadUrl string) error {
	tempDirPath := filepath.Join(i.csGamePath, fmt.Sprintf("temp_%v_%v", plugin.Name, time.Now().UTC()))
	if err := os.Mkdir(tempDirPath, os.ModePerm); err != nil {
		return fmt.Errorf("os.Mkdir temp dir: %w", err)
	}

	downloadedFilePath, err := download.Download(downloadUrl, tempDirPath)
	if err != nil {
		return fmt.Errorf("download.Download: %w", err)
	}

	destPath := filepath.Join(i.csGamePath, plugin.InstallDir)

	if err := unzipDownload(downloadedFilePath, destPath); err != nil {
		return fmt.Errorf("unzipDownload: %w", err)
	}

	if err := os.RemoveAll(tempDirPath); err != nil {
		return fmt.Errorf("os.RemoveAll temp dir: %w", err)
	}

	return nil
}

func unzipDownload(sourcePath, destDir string) error {
	if strings.HasSuffix(sourcePath, ".zip") {
		if err := unzip.Zip(sourcePath, destDir); err != nil {
			return fmt.Errorf("unzip zip: %w", err)

		}
	} else if strings.HasSuffix(sourcePath, ".tar.gz") {
		if err := unzip.TarGz(sourcePath, destDir); err != nil {
			return fmt.Errorf("unzip tar.gz: %w", err)
		}
	} else {
		return fmt.Errorf("downloaded filetype not supported '%v'", sourcePath)
	}

	return nil
}
