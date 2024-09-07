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
	Name         string                     `json:"name"`
	Installed    bool                       `json:"installed"`
	Dependencies []PluginDependencyResponse `json:"dependencies"`
}

type PluginDependencyResponse struct {
	Name         string                     `json:"name"`
	InstallDir   string                     `json:"install_dir"`
	Version      string                     `json:"version"`
	DownloadURL  string                     `json:"download_url"`
	Dependencies []PluginDependencyResponse `json:"dependencies"`
}

func RegisterPlugins(r fiber.Router) {
	r.Get("/plugins", getPluginsHandler)
	r.Post("/plugins", installPluginHandler)
	r.Delete("/plugins", uninstallPluginHandler)
}

// getPluginsHandler
// @Summary				Get all available plugins
// @Tags         		plugins
// @Produce      		json
// @Success     		200  {object}  PluginResponse
// @Failure				400  {object}  handlers.ErrorResponse
// @Failure				500  {object}  handlers.ErrorResponse
// @Router       		/plugins [get]
func getPluginsHandler(c fiber.Ctx) error {
	pluginsInstance, err := GetFromLocals[*plugins.Instance](c, constants.PluginsKey)
	if err != nil {
		return NewInternalServerErrorWithInternal(c, fmt.Errorf("failed to get plugins instance from context: %w", err))
	}

	installedPlugin, err := pluginsInstance.GetInstalledPlugin()
	if err != nil {
		return NewInternalServerErrorWithInternal(c, fmt.Errorf("failed to get installed plugins: %w", err))
	}

	availablePlugins := pluginsInstance.GetAllAvailablePlugins()
	if availablePlugins == nil {
		return c.Status(fiber.StatusOK).JSON(make([]PluginResponse, 0))
	}

	result := make([]PluginResponse, 0, len(availablePlugins))

	for _, plugin := range availablePlugins {
		versions := make([]PluginVersionResponse, 0, len(plugin.Versions))
		for _, version := range plugin.Versions {
			versions = append(versions, PluginVersionResponse{
				Name:         version.Name,
				Installed:    installedPlugin != nil && installedPlugin.Name == plugin.Name && installedPlugin.Version == version.Name,
				Dependencies: mapPluginDependencyToPluginDependencyResponses(version.Dependencies),
			})
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

func mapPluginDependencyToPluginDependencyResponses(dependencies []plugins.PluginDependency) []PluginDependencyResponse {
	if len(dependencies) == 0 {
		return nil
	}

	var result = make([]PluginDependencyResponse, 0, len(dependencies))
	for _, d := range dependencies {
		result = append(result, PluginDependencyResponse{
			Name:         d.Name,
			InstallDir:   d.InstallDir,
			Version:      d.Version,
			DownloadURL:  d.DownloadURL,
			Dependencies: mapPluginDependencyToPluginDependencyResponses(d.Dependencies),
		})
	}

	return result
}

type InstallPluginRequest struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

// installPluginHandler
// @Summary				Install given plugin
// @Tags         		plugins
// @Param		 		plugin body InstallPluginRequest true "The plugin and version that should be installed"
// @Accept       		json
// @Success     		200
// @Failure				400  {object}  handlers.ErrorResponse
// @Failure				500  {object}  handlers.ErrorResponse
// @Router       		/plugins [post]
func installPluginHandler(c fiber.Ctx) error {
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
	if err := c.Bind().JSON(&installPluginRequest); err != nil {
		return NewErrorWithInternal(c, fiber.StatusBadRequest, "request is not valid", err)
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

// uninstallPluginHandler
// @Summary				Uninstall the currently installed plugin
// @Tags         		plugins
// @Success     		200
// @Failure				400  {object}  handlers.ErrorResponse
// @Failure				500  {object}  handlers.ErrorResponse
// @Router       		/plugins [delete]
func uninstallPluginHandler(c fiber.Ctx) error {
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
