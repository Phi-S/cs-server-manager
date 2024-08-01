package errorwrp

import (
    "errors"
    "net/http"
)

type HttpResponse struct {
    Status   int
    Response any
}

type HttpError struct {
    Status          int
    ResponseMessage string
    InternalError   error
}

func NewJsonHttpResponse(status int, resp any) (HttpResponse, *HttpError) {
    return HttpResponse{
        Status:   status,
        Response: resp,
    }, nil
}

func NewOkJsonHttpResponse(resp any) (HttpResponse, *HttpError) {
    return HttpResponse{
        Status:   http.StatusOK,
        Response: resp,
    }, nil
}

func NewHttpResponse(status int) (HttpResponse, *HttpError) {
    return HttpResponse{
        Status:   status,
        Response: nil,
    }, nil
}

func NewOkHttpResponse() (HttpResponse, *HttpError) {
    return HttpResponse{
        Status:   http.StatusOK,
        Response: nil,
    }, nil
}

func NewHttpError(status int, msg string, err error) (HttpResponse, *HttpError) {
    return HttpResponse{}, &HttpError{
        Status:          status,
        ResponseMessage: msg,
        InternalError:   err,
    }
}

func NewHttpError2(status int, msg string) (HttpResponse, *HttpError) {
    return NewHttpError(status, msg, errors.New(msg))
}

func NewHttpErrorInternalServerError(msg string, err error) (HttpResponse, *HttpError) {
    return NewHttpError(http.StatusInternalServerError, msg, err)
}

func NewHttpErrorInternalServerError2(msg string) (HttpResponse, *HttpError) {
    return NewHttpError2(http.StatusInternalServerError, msg)
}
