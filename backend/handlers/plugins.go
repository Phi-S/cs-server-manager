package handlers

import (
	"cs-server-manager/constants"
	"cs-server-manager/plugins"
	"fmt"
	"github.com/gofiber/fiber/v3"
)

type PluginResponse struct {
	Name        string                  `json:"name"`
	Description string                  `json:"description"`
	URL         string                  `json:"url"`
	Versions    []PluginVersionResponse `json:"versions"`
}

type PluginVersionResponse struct {
	Name         string                              `json:"name"`
	Installed    bool                                `json:"installed"`
	Dependencies []PluginVersionDependenciesResponse `json:"dependencies"`
}

type PluginVersionDependenciesResponse struct {
	PluginName  string `json:"plugin_name"`
	VersionName string `json:"version_name"`
}

func GetAvailablePluginsHandler(c fiber.Ctx) error {
	pluginsInstance, err := GetFromLocals[*plugins.Instance](c, constants.PluginsKey)
	if err != nil {
		return NewInternalServerErrorWithInternal(c, err)
	}

	installedPlugin, err := pluginsInstance.GetInstalledPlugin()
	if err != nil {
		return NewInternalServerErrorWithInternal(c, err)
	}

	availablePlugins := pluginsInstance.GetAllAvailablePlugins()

	result := make([]PluginResponse, 0, len(availablePlugins))

	for _, plugin := range availablePlugins {
		versions := make([]PluginVersionResponse, 0)
		for _, version := range plugin.Versions {
			versionDependencies := make([]PluginVersionDependenciesResponse, 0)
			for _, dependency := range version.Dependencies {
				versionDependencies = append(versionDependencies, PluginVersionDependenciesResponse{
					PluginName:  dependency.PluginName,
					VersionName: dependency.VersionName,
				})
			}

			versionResponse := PluginVersionResponse{
				Name:         version.Name,
				Installed:    false,
				Dependencies: versionDependencies,
			}

			versionResponse.Installed = installedPlugin != nil && installedPlugin.Name == plugin.Name && installedPlugin.Version == version.Name
			versions = append(versions, versionResponse)
		}

		result = append(result, PluginResponse{
			Name:        plugin.Name,
			Description: plugin.Description,
			URL:         plugin.URL,
			Versions:    versions,
		})
	}

	return c.Status(fiber.StatusOK).JSON(result)
}

func InstallPluginHandler(c fiber.Ctx) error {
	pluginsInstance, err := GetFromLocals[*plugins.Instance](c, constants.PluginsKey)
	if err != nil {
		return NewInternalServerErrorWithInternal(c, err)
	}

	lock, serverInstance, steamcmdInstance, err := GetServerSteamcmdInstances(c)
	if err != nil {
		return NewInternalServerErrorWithInternal(c, fmt.Errorf("GetServerSteamcmdInstances: %w", err))
	}

	lock.Lock()
	defer lock.Unlock()
	if serverInstance.IsRunning() {
		return NewErrorWithMessage(c, fiber.StatusInternalServerError, "can not install plugins while server is running")
	}

	if steamcmdInstance.IsRunning() {
		return NewErrorWithMessage(c, fiber.StatusInternalServerError, "can not install plugins while steamcmd is running")
	}

	pluginName := c.Query("name")
	versionName := c.Query("version")

	installedPlugin, err := pluginsInstance.GetInstalledPlugin()
	if err != nil {
		return NewInternalServerErrorWithInternal(c, fmt.Errorf("GetInstalledPlugin: %w", err))
	}

	if installedPlugin != nil && installedPlugin.Name == pluginName && installedPlugin.Version == versionName {
		return c.SendStatus(fiber.StatusAlreadyReported)
	}

	if err := pluginsInstance.InstallPluginByName(pluginName, versionName); err != nil {
		return NewInternalServerErrorWithInternal(c, err)
	}

	return c.SendStatus(fiber.StatusOK)
}

func UninstallPluginHandler(c fiber.Ctx) error {
	pluginsInstance, err := GetFromLocals[*plugins.Instance](c, constants.PluginsKey)
	if err != nil {
		return NewInternalServerErrorWithInternal(c, err)
	}

	lock, serverInstance, steamcmdInstance, err := GetServerSteamcmdInstances(c)
	if err != nil {
		return NewInternalServerErrorWithInternal(c, fmt.Errorf("GetServerSteamcmdInstances: %w", err))
	}

	lock.Lock()
	defer lock.Unlock()
	if serverInstance.IsRunning() {
		return NewErrorWithMessage(c, fiber.StatusInternalServerError, "can not uninstall plugins while server is running")
	}

	if steamcmdInstance.IsRunning() {
		return NewErrorWithMessage(c, fiber.StatusInternalServerError, "can not uninstall plugins while steamcmd is running")
	}

	pluginName := c.Query("name")

	if err := pluginsInstance.Uninstall(pluginName); err != nil {
		return NewInternalServerErrorWithInternal(c, fmt.Errorf("pluginsInstance.Uninstall: %w", err))
	}

	return c.SendStatus(fiber.StatusOK)
}
