package handlers

import (
	"cs-server-manager/constants"
	"cs-server-manager/logwrt"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"
)

// LogsHandler
// @Summary				Gets logs
// @Tags         		logs
// @Produce     		json
// @Param countOrSince path int false "Gets the last x logs or all logs since date" 100
// @Success     		200  {object}  []logwrt.LogEntry
// @Failure				400  {object}  handlers.ErrorResponse
// @Failure				500  {object}  handlers.ErrorResponse
// @Router       		/logs/{countOrSince} [get]
func LogsHandler(c fiber.Ctx) error {
	logWriter, err := GetFromLocals[*logwrt.LogWriter](c, constants.UserLogWriterKey)
	if err != nil {
		return NewInternalServerErrorWithInternal(c, err)
	}

	countOrSince := strings.TrimSpace(c.Params("countOrSince"))
	if countOrSince == "" {
		return fiber.NewError(fiber.StatusBadGateway, "expected parameter is empty")
	}

	count, countErr := strconv.ParseInt(countOrSince, 10, 64)
	since, sinceErr := time.Parse(time.RFC3339Nano, countOrSince)
	if countErr != nil && sinceErr != nil {
		return fiber.NewError(fiber.StatusBadGateway, "expected parameter can only be a number of timestamp")
	}

	var result []logwrt.LogEntry
	if countErr == nil {
		result, err = logWriter.GetLogs(int(count))
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to get logs. %v", err))
		}
	} else {
		result, err = logWriter.GetLogsSince(since)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to get logs. %v", err))
		}
	}

	return c.Status(fiber.StatusOK).JSON(result)
}

/*


func LogsSinceHandler(r *http.Request) (errorwrp.HttpResponse, *errorwrp.HttpError) {
	logWriter, err := ctxex.GetUserLogWriter(r.Context())
	if err != nil {
		return errorwrp.NewHttpErrorInternalServerError("internal error", err)
	}

	sinceStr := r.URL.Query().Get("since")
	if sinceStr == "" {
		return errorwrp.NewHttpError2(http.StatusBadRequest, "since parameter missing")
	}

	since, err := time.Parse(time.RFC3339Nano, sinceStr)
	if err != nil {
		return errorwrp.NewHttpError2(http.StatusBadRequest, "since parameter is not of valid RFC3339Nano format")
	}

	logs, err := logWriter.GetLogsSince(since)
	if err != nil {
		return errorwrp.NewHttpErrorInternalServerError("internal error", err)
	}

	return errorwrp.NewOkJsonHttpResponse(logs)
}

func LogFilesHandler(r *http.Request) (errorwrp.HttpResponse, *errorwrp.HttpError) {
	logWriter, err := ctxex.GetUserLogWriter(r.Context())
	if err != nil {
		return errorwrp.NewHttpErrorInternalServerError("internal error", err)
	}

	files, err := logWriter.GetPastLogFiles()
	if err != nil {
		return errorwrp.NewHttpErrorInternalServerError("internal error", err)
	}

	resp := struct {
		LogFiles []string `json:"log-files"`
	}{
		LogFiles: files,
	}
	return errorwrp.NewOkJsonHttpResponse(resp)
}

func LogFileContentHandler(r *http.Request) (errorwrp.HttpResponse, *errorwrp.HttpError) {
	logWriter, err := ctxex.GetUserLogWriter(r.Context())
	if err != nil {
		return errorwrp.NewHttpErrorInternalServerError("internal error", err)
	}

	fileName := r.URL.Query().Get("name")
	if fileName == "" {
		return errorwrp.NewHttpError2(http.StatusBadRequest, "name parameter missing")
	}

	logs, err := logWriter.GetContentOfPastLogFile(fileName)
	if err != nil {
		return errorwrp.NewHttpErrorInternalServerError("internal error", err)
	}

	return errorwrp.NewOkJsonHttpResponse(logs)
}
*/
