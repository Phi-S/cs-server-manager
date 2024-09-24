package handlers

import (
	"net/url"

	"github.com/Phi-S/cs-server-manager/constants"
	"github.com/Phi-S/cs-server-manager/editor"
	"github.com/Phi-S/cs-server-manager/gvalidator"

	"github.com/gofiber/fiber/v3"
)

func RegisterFiles(r fiber.Router) {
	r.Get("/files", getAllEditableFilesHandler)
	r.Get("/files/:file", getFileContent)
	r.Patch("/files/:file", setFileContent)
}

type FilesResponse struct {
	Files []string `json:"files"`
}

// @Summary				Get editable files
// @Tags         		files
// @Produce     		json
// @Success     		200  	{object}	[]FilesResponse
// @Failure				400  	{object}	handlers.ErrorResponse
// @Failure				500  	{object}	handlers.ErrorResponse
// @Router       		/files [get]
func getAllEditableFilesHandler(c fiber.Ctx) error {
	editorInstance, err := GetFromLocals[*editor.Instance](c, constants.EditorKey)
	if err != nil {
		return NewInternalServerErrorWithInternal(c, err)
	}

	editableFiles, err := editorInstance.GetAllEditableFiles()
	if err != nil {
		return NewErrorWithInternal(c, fiber.StatusInternalServerError, "failed to get editable files", err)
	}

	if editableFiles == nil {
		editableFiles = make([]string, 0)
	}

	result := FilesResponse{Files: editableFiles}
	return c.Status(fiber.StatusOK).JSON(result)
}

// @Summary				Get files content
// @Tags         		files
// @Produce     		plain
// @Param 				file	path		string true "file to get content for"
// @Success     		200  	{string}	string
// @Failure				400  	{object}	handlers.ErrorResponse
// @Failure				500  	{object}	handlers.ErrorResponse
// @Router       		/files/{file} [get]
func getFileContent(c fiber.Ctx) error {
	fileParam := c.Params("file")

	fileParam, err := url.QueryUnescape(fileParam)
	if err != nil {
		return NewInternalServerErrorWithInternal(c, err)
	}

	err = gvalidator.Instance().Var(fileParam, "required,filepath")
	if err != nil {
		return NewErrorValidation(c, err)
	}

	editorInstance, err := GetFromLocals[*editor.Instance](c, constants.EditorKey)
	if err != nil {
		return NewInternalServerErrorWithInternal(c, err)
	}

	content, err := editorInstance.GetFileContent(fileParam)
	if err != nil {
		return NewErrorWithInternal(c, fiber.StatusInternalServerError, "failed to get file content", err)
	}

	return c.Status(fiber.StatusOK).SendString(string(content))
}

// @Summary				Set files content
// @Tags         		files
// @Accept 				plain
// @Param 				file 	path string true "file to set content for"
// @Param				content body string true "file content"
// @Success     		200
// @Failure				400  	{object}	handlers.ErrorResponse
// @Failure				500  	{object}	handlers.ErrorResponse
// @Router       		/files{file} [PATCH]
func setFileContent(c fiber.Ctx) error {
	fileParam := c.Params("file")
	fileParam, err := url.QueryUnescape(fileParam)
	if err != nil {
		return NewInternalServerErrorWithInternal(c, err)
	}

	err = gvalidator.Instance().Var(fileParam, "required,filepath")
	if err != nil {
		return NewErrorValidation(c, err)
	}

	editorInstance, err := GetFromLocals[*editor.Instance](c, constants.EditorKey)
	if err != nil {
		return NewInternalServerErrorWithInternal(c, err)
	}

	if len(c.Body())/1024/1024 > 2 {
		return NewErrorWithMessage(c, fiber.StatusRequestEntityTooLarge, "file content can not be bigger then 2 MB")
	}

	if err := editorInstance.SetFileContent(fileParam, c.BodyRaw()); err != nil {
		return NewErrorWithInternal(c, fiber.StatusInternalServerError, "failed to write file content", err)
	}

	return c.SendStatus(fiber.StatusOK)
}
