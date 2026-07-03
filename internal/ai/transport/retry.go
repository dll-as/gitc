package transport

import (
	"context"
	"errors"
	"net"
	"net/http"
)

func shouldRetry(resp *http.Response, err error) bool {
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return false
		}

		if errors.Is(err, context.DeadlineExceeded) {
			return true
		}

		var netErr net.Error
		if errors.As(err, &netErr) {
			return true
		}

		return true
	}

	switch resp.StatusCode {

	case http.StatusTooManyRequests,
		http.StatusInternalServerError,
		http.StatusBadGateway,
		http.StatusServiceUnavailable,
		http.StatusGatewayTimeout:

		return true
	}

	return false
}
