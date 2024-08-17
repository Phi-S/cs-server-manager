package handlers

import (
	"cs-server-manager/constants"
	"cs-server-manager/gvalidator"
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

	installedPlugins, err := pluginsInstance.GetInstalledPlugins()
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

			if isPluginVersionInstalled(plugin, version.Name, installedPlugins) {
				versionResponse.Installed = true
			}

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

func isPluginVersionInstalled(plugin plugins.Plugin, versionName string, installedPlugins []plugins.InstalledPlugin) bool {
	for _, installedPlugin := range installedPlugins {
		if plugin.Name == installedPlugin.Name && versionName == installedPlugin.Version {
			return true
		}
	}

	return false
}

func InstallPluginHandler(c fiber.Ctx) error {
	pluginsInstance, err := GetFromLocals[*plugins.Instance](c, constants.PluginsKey)
	if err != nil {
		return NewInternalServerErrorWithInternal(c, err)
	}

	pluginName := c.Query("name")
	versionName := c.Query("version")

	if err := gvalidator.Instance().Var(pluginName, "required,lte=32"); err != nil {
		return NewErrorWithInternal(c, fiber.StatusBadRequest, "plugin name is not valid", err)
	}

	if err := gvalidator.Instance().Var(versionName, "required,lte=32"); err != nil {
		return NewErrorWithInternal(c, fiber.StatusBadRequest, "version name is not valid", err)
	}

	installedPlugins, err := pluginsInstance.GetInstalledPlugins()
	if err != nil {
		return NewInternalServerErrorWithInternal(c, fmt.Errorf("GetInstalledPlugins: %w", err))
	}

	for _, plugin := range installedPlugins {
		if plugin.Name == pluginName && plugin.Version == versionName {
			return c.SendStatus(fiber.StatusAlreadyReported)
		}
	}

	if err := pluginsInstance.InstallPluginByName(pluginName, versionName); err != nil {
		return NewInternalServerErrorWithInternal(c, err)
	}

	return c.SendStatus(fiber.StatusOK)
}
