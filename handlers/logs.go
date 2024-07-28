package handlers

import (
	"cs-server-controller/httpex/errorwrp"
	"cs-server-controller/middleware"
	"cs-server-controller/user_logs"
	"errors"
	"net/http"
	"strconv"
	"time"
)

func LogHandler(r *http.Request) (errorwrp.HttpResponse, *errorwrp.HttpError) {
	userLogsWriter, ok := r.Context().Value(middleware.UserLogWriterKey).(*user_logs.LogWriter)
	if userLogsWriter == nil || !ok {
		return errorwrp.NewHttpErrorInternalServerError("internal error", errors.New("failed to get user logs writer from context"))
	}

	countStr := r.URL.Query().Get("count")
	if countStr == "" {
		return errorwrp.NewHttpErrorInternalServerError2("count parameter missing")
	}

	count, err := strconv.Atoi(countStr)
	if err != nil || count < 0 || count > 200 {
		return errorwrp.NewHttpErrorInternalServerError2("count parameter is not a valid number between 1 and 200")
	}

	sinceStr := r.URL.Query().Get("since")
	if sinceStr != "" {
		t, err := time.Parse(time.RFC3339, sinceStr)
		if err != nil {
			return errorwrp.NewHttpErrorInternalServerError2("since parameter is not of valid RFC3339 timestamp")
		}

		logs := userLogsWriter.GetLogsSince(t, count)
		return errorwrp.NewOkJsonHttpResponse(logs)
	} else {
		logs := userLogsWriter.GetLogs(count)
		return errorwrp.NewOkJsonHttpResponse(logs)
	}
}
