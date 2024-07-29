package handlers

import (
	"cs-server-controller/httpex/errorwrp"
	"cs-server-controller/logwrt"
	"cs-server-controller/middleware"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

func LogsHandler(r *http.Request) (errorwrp.HttpResponse, *errorwrp.HttpError) {
	logWriter, ok := r.Context().Value(middleware.UserLogWriterKey).(*logwrt.LogWriter)
	if logWriter == nil || !ok {
		return errorwrp.NewHttpErrorInternalServerError("internal error", errors.New("failed to get user logs writer from context"))
	}

	countStr := r.URL.Query().Get("count")
	if countStr == "" {
		return errorwrp.NewHttpError2(http.StatusBadRequest, "count parameter missing")
	}

	count, err := strconv.Atoi(countStr)
	if err != nil || count < 0 || count > logWriter.GetLogsLimit() {
		return errorwrp.NewHttpError2(
			http.StatusBadRequest,
			fmt.Sprintf("count parameter is not a valid number between 1 and %v", logWriter.GetLogsLimit()),
		)
	}

	logs, err := logWriter.GetLogs(count)
	if err != nil {
		return errorwrp.NewHttpErrorInternalServerError("internal error", err)
	}
	return errorwrp.NewOkJsonHttpResponse(logs)

}

func LogsSinceHandler(r *http.Request) (errorwrp.HttpResponse, *errorwrp.HttpError) {
	logWriter, ok := r.Context().Value(middleware.UserLogWriterKey).(*logwrt.LogWriter)
	if logWriter == nil || !ok {
		return errorwrp.NewHttpErrorInternalServerError("internal error", errors.New("failed to get user logs writer from context"))
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
	logWriter, ok := r.Context().Value(middleware.UserLogWriterKey).(*logwrt.LogWriter)
	if logWriter == nil || !ok {
		return errorwrp.NewHttpErrorInternalServerError("internal error", errors.New("failed to get user logs writer from context"))
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
	logWriter, ok := r.Context().Value(middleware.UserLogWriterKey).(*logwrt.LogWriter)
	if logWriter == nil || !ok {
		return errorwrp.NewHttpErrorInternalServerError("internal error", errors.New("failed to get user logs writer from context"))
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
