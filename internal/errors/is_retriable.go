package errors

import (
	"net"
	"net/http"

	"github.com/go-resty/resty/v2"
	"github.com/jackc/pgerrcode"
)

func IsRetriableError(err error) bool {
	if err == nil {
		return false
	}
	if netErr, ok := err.(net.Error); ok {
		return netErr.Timeout()
	}
	if err.Error() == pgerrcode.ConnectionException {
		return true
	}
	if restyErr, ok := err.(*resty.ResponseError); ok {
		if restyErr.Response != nil {
			if restyErr.Response.StatusCode() == http.StatusRequestTimeout || restyErr.Response.StatusCode() == http.StatusGatewayTimeout {
				return true
			}
		}
	}
	return false
}
