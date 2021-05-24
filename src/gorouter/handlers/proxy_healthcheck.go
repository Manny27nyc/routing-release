package handlers

import (
	"code.cloudfoundry.org/routing-release/gorouter/common/health"
	"code.cloudfoundry.org/routing-release/gorouter/logger"
	"github.com/urfave/negroni"
	"net/http"
)

type proxyHealthcheck struct {
	userAgent string
	health    *health.Health
	logger    logger.Logger
}

// NewHealthcheck creates a handler that responds to healthcheck requests.
// If userAgent is set to a non-empty string, it will use that user agent to
// differentiate between healthcheck requests and non-healthcheck requests.
// Otherwise, it will treat all requests as healthcheck requests.
func NewProxyHealthcheck(userAgent string, health *health.Health, logger logger.Logger) negroni.Handler {
	return &proxyHealthcheck{
		userAgent: userAgent,
		health:    health,
		logger:    logger,
	}
}

func (h *proxyHealthcheck) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	// If reqeust is not intended for healthcheck
	if r.Header.Get("User-Agent") != h.userAgent {
		next(rw, r)
		return
	}

	rw.Header().Set("Cache-Control", "private, max-age=0")
	rw.Header().Set("Expires", "0")

	if h.health.Health() != health.Healthy {
		rw.WriteHeader(http.StatusServiceUnavailable)
		r.Close = true
		return
	}

	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("ok\n"))
	r.Close = true
}