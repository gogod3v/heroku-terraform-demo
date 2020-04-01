package main

import (
	"log"
	"net/http"

	"github.com/justinas/alice"
)

// MiddlewareProvider ...
type MiddlewareProvider struct {
	appName    string
	appVersion string
}

// NewMiddlewareProvider ...
func NewMiddlewareProvider(appName, version string) MiddlewareProvider {
	return MiddlewareProvider{
		appName:    appName,
		appVersion: version,
	}
}

func createRequestLoggerMiddleware(appName, version string) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Printf("-> %s (%s) [%s] %s", appName, version, r.Method, r.RequestURI)
			h.ServeHTTP(w, r)
		})
	}
}

// CommonMiddleware ...
func (m MiddlewareProvider) CommonMiddleware() alice.Chain {
	commonMiddleware := alice.New()

	return commonMiddleware.Append(
		createRequestLoggerMiddleware(m.appName, m.appVersion),
	)
}
