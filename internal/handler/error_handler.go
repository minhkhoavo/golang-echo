package handler

import (
	"errors"
	"fmt"
	"golang-echo/pkg/response"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

func CustomHTTPErrorHandler(err error, c echo.Context) {
	// 1. If response is already committed, we cannot send any more data
	if c.Response().Committed {
		return
	}

	var (
		code = http.StatusInternalServerError
		key  = "SERVER_INTERNAL_ERROR"
		msg  = "Internal Server Error"
	)

	// 2. Classify the error (Type Assertion & Error Wrapping Check)
	var appErr *response.AppError
	var echoErr *echo.HTTPError

	if errors.As(err, &appErr) {
		// This is an application error (AppError)
		code = appErr.Code
		key = appErr.Key
		msg = appErr.Message
		// Log the underlying error
		// if appErr.Err != nil {
		// 	log.Errorf("App Error occurred: %v", appErr.Err)
		// }

	} else if errors.As(err, &echoErr) {
		// This is an Echo error (e.g., 404 Route not found, 405 Method not allowed)
		code = echoErr.Code
		key = "ECHO_HTTP_ERROR"
		msg = fmt.Sprintf("%v", echoErr.Message)
	} else {
		// Unknown error (Unknown panic or third-party library error)
		// Only log this error, DO NOT return details to client for security reasons
		log.Errorf("Unhandled error occurred: %v", err)
	}

	// 3. Build standard response
	// Get Request ID if available
	reqID := c.Response().Header().Get(echo.HeaderXRequestID)

	response := map[string]interface{}{
		"status":     code,
		"error_code": key,
		"message":    msg,
		"request_id": reqID,
	}

	// 4. Send response to client
	if c.Request().Method == http.MethodHead {
		err = c.NoContent(code)
	} else {
		err = c.JSON(code, response)
	}

	// Fallback if sending JSON also fails
	if err != nil {
		c.Logger().Error(err)
	}
}
