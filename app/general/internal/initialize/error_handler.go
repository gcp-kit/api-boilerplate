package initialize

import (
	"encoding/json"
	"net/http"

	"github.com/cockroachdb/errors"
	"github.com/getkin/kin-openapi/routers"
	"github.com/labstack/echo/v4"
)

type Map map[string]interface{}

// ErrorHandler - error handler
func ErrorHandler() echo.HTTPErrorHandler {
	return func(err error, c echo.Context) {
		if c.Response().Committed {
			return
		}

		var (
			httpErr  *echo.HTTPError
			routeErr *routers.RouteError
		)

		if errors.As(err, &httpErr) {
			if httpErr.Internal != nil {
				if herr, ok := httpErr.Internal.(*echo.HTTPError); ok {
					httpErr = herr
				}
			}
			// Send response
			if c.Request().Method == http.MethodHead {
				err = c.NoContent(httpErr.Code)
			} else {
				err = c.JSON(httpErr.Code, httpErr.Message)
			}
			if err != nil {
				c.Logger().Error(err)
			}
			return
		}

		if errors.As(err, &routeErr) {
			var he *echo.HTTPError
			switch {
			case errors.Is(routeErr, routers.ErrPathNotFound):
				he = &echo.HTTPError{
					Code:    http.StatusNotFound,
					Message: http.StatusText(http.StatusNotFound),
				}
			case errors.Is(routeErr, routers.ErrMethodNotAllowed):
				he = &echo.HTTPError{
					Code:    http.StatusMethodNotAllowed,
					Message: http.StatusText(http.StatusMethodNotAllowed),
				}
			default:
				he = &echo.HTTPError{
					Code:    http.StatusInternalServerError,
					Message: http.StatusText(http.StatusInternalServerError),
				}
			}

			// Send response
			if c.Request().Method == http.MethodHead {
				err = c.NoContent(he.Code)
			} else {
				err = c.JSON(he.Code, he.Message)
			}
			if err != nil {
				c.Logger().Error(err)
			}
			return
		}

		he := &echo.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: http.StatusText(http.StatusInternalServerError),
		}

		code := he.Code
		message := he.Message

		switch m := he.Message.(type) {
		case string:
			if c.Echo().Debug {
				message = Map{"message": m, "error": err.Error()}
			} else {
				message = Map{"message": m}
			}
		case json.Marshaler:
			// do nothing - this type knows how to format itself to JSON
		case error:
			message = Map{"message": m.Error()}
		}

		// Send response
		if c.Request().Method == http.MethodHead { // Issue #608
			err = c.NoContent(he.Code)
		} else {
			err = c.JSON(code, message)
		}
		if err != nil {
			c.Logger().Error(err)
		}
	}
}
