package handlers

import (
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

// MigrationHandler keeps the existing application as the public entry point
// while device and telemetry functionality is moved to microservices.
type MigrationHandler struct {
	deviceProxy  *httputil.ReverseProxy
	historyProxy *httputil.ReverseProxy
}

func NewMigrationHandler(deviceURL, historyURL string) (*MigrationHandler, error) {
	deviceTarget, err := url.Parse(deviceURL)
	if err != nil {
		return nil, err
	}
	historyTarget, err := url.Parse(historyURL)
	if err != nil {
		return nil, err
	}
	return &MigrationHandler{
		deviceProxy:  httputil.NewSingleHostReverseProxy(deviceTarget),
		historyProxy: httputil.NewSingleHostReverseProxy(historyTarget),
	}, nil
}

func (h *MigrationHandler) RegisterRoutes(router *gin.RouterGroup) {
	router.Any("/devices", h.proxy(h.deviceProxy))
	router.Any("/devices/*path", h.proxy(h.deviceProxy))
	router.Any("/telemetry", h.proxy(h.historyProxy))
}

func (h *MigrationHandler) proxy(target *httputil.ReverseProxy) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request.URL.Path = strings.TrimPrefix(c.Request.URL.Path, "/api")
		target.ServeHTTP(c.Writer, c.Request)
	}
}
