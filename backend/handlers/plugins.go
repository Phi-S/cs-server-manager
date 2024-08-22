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

// GetAvailablePluginsHandler	GetPlugins
// @Summary				Gets all available plugins
// @Tags         		plugins
// @Success     		200  {object}  PluginResponse
// @Failure				400  {object}  handlers.ErrorResponse
// @Failure				500  {object}  handlers.ErrorResponse
// @Router       		/plugins [get]
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

type InstallPluginRequest struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

// InstallPluginHandler
// @Summary				Installs the given plugin or updates to given version
// @Tags         		plugins
// @Param		 		plugin body InstallPluginRequest true "The plugin and the version that should be installed"
// @Success     		200
// @Failure				400  {object}  handlers.ErrorResponse
// @Failure				500  {object}  handlers.ErrorResponse
// @Router       		/plugins [post]
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

	var installPluginRequest InstallPluginRequest
	if err := c.Bind().JSON(installPluginRequest); err != nil {
		return NewErrorWithMessage(c, fiber.StatusBadRequest, "request is not valid")
	}

	installedPlugin, err := pluginsInstance.GetInstalledPlugin()
	if err != nil {
		return NewInternalServerErrorWithInternal(c, fmt.Errorf("GetInstalledPlugin: %w", err))
	}

	if installedPlugin != nil && installedPlugin.Name == installPluginRequest.Name && installedPlugin.Version == installPluginRequest.Version {
		return c.SendStatus(fiber.StatusAlreadyReported)
	}

	if err := pluginsInstance.InstallPluginByName(installPluginRequest.Name, installPluginRequest.Version); err != nil {
		return NewInternalServerErrorWithInternal(c, err)
	}

	return c.SendStatus(fiber.StatusOK)
}

// UninstallPluginHandler
// @Summary				Uninstalls the currently installed plugin
// @Tags         		plugins
// @Success     		200
// @Failure				400  {object}  handlers.ErrorResponse
// @Failure				500  {object}  handlers.ErrorResponse
// @Router       		/plugins [delete]
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

	if err := pluginsInstance.Uninstall(); err != nil {
		return NewInternalServerErrorWithInternal(c, fmt.Errorf("pluginsInstance.Uninstall: %w", err))
	}

	return c.SendStatus(fiber.StatusOK)
}
