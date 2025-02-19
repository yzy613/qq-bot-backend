package errcode

import (
	"github.com/gogf/gf/v2/errors/gcode"
	"net/http"
)

var (
	Unauthorized  = gcode.New(http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized), nil)
	Forbidden     = gcode.New(http.StatusForbidden, http.StatusText(http.StatusForbidden), nil)
	Conflict      = gcode.New(http.StatusConflict, http.StatusText(http.StatusConflict), nil)
	TooEarly      = gcode.New(http.StatusTooEarly, http.StatusText(http.StatusTooEarly), nil)
	TooMany       = gcode.New(http.StatusTooManyRequests, http.StatusText(http.StatusTooManyRequests), nil)
	InternalError = gcode.New(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), nil)
)
