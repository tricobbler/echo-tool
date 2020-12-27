package httpError

import (
	"github.com/labstack/echo/v4"
	"github.com/maybgit/glog"
	"net/http"
)

type httpError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func NewHTTPError(code int, msg string) *httpError {
	return &httpError{
		Code:    code,
		Message: msg,
	}
}

func NewDefaultHttpError() *httpError {
	return NewHTTPError(http.StatusInternalServerError, DefaultError)
}

// Error makes it compatible with `error` interface.
func (e *httpError) Error() string {
	return e.Message
}

const DefaultError = "网络开小差了，请稍后重试"

// httpErrorHandler customize echo's HTTP error handler.
func HttpErrorHandler(err error, c echo.Context) {
	var (
		code = http.StatusInternalServerError
		msg  = DefaultError
	)

	if he, ok := err.(*httpError); ok {
		code = he.Code
		msg = he.Message
	} else if ee, ok := err.(*echo.HTTPError); ok {
		code = ee.Code
		msg = http.StatusText(code)
	} else if c.Echo().Debug {
		msg = err.Error()
	} else {
		msg = http.StatusText(code)
	}

	if !c.Response().Committed {
		if c.Request().Method == echo.HEAD {
			err := c.NoContent(code)
			if err != nil {
				glog.Error(err)
			}
		} else {
			err := c.JSON(200, NewHTTPError(code, msg))
			if err != nil {
				glog.Error(err)
			}
		}
	}
}
