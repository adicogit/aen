package server

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"aen.it/poolmanager/config"
	"aen.it/poolmanager/log"
)

const (
	DefaultStaticPath = "/opt/baas-server/website/"
	DefaultIndexPath  = "index.html"
)

// webuiHandler implements the http.Handler interface, so we can use it
// to respond to HTTP requests.
type webuiHandler struct {
	staticPath string
	indexPath  string
}

type pathError struct {
	statusCode int
	message    string
}

// WebuiHandler is the handler for mricrofrontend UI
var WebuiHandler *webuiHandler

func init() {
	WebuiHandler = &webuiHandler{}
	WebuiHandler.staticPath = DefaultStaticPath
	WebuiHandler.indexPath = DefaultIndexPath
}

func (h webuiHandler) SetStaticPath(staticPath string) {
	log.Log.Debug("Entering SetStaticPath")
	if len(staticPath) > 0 {
		WebuiHandler.staticPath = staticPath
	} else {
		log.Log.Warn("Requested to change UI static path to empty value. Still using current value", "currentValue", WebuiHandler.staticPath)
	}
	log.Log.Debug("Exiting SetStaticPath")
}

// ServeHTTP inspects the URL path to locate a file within the static dir
// on the SPA handler. If a file is found, it will be served. If not, the
// file located at the index path on the SPA handler will be served. This
// is suitable behavior for serving an SPA (single page application).
func (h webuiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Log.Debug("Entering ServeHTTP")

	cleanPath := filepath.Clean(r.URL.Path)
	log.Log.Debug("Retrieve path from HTTP request", "path", cleanPath)

	path, err := h.validatePath(cleanPath, r.URL.Path)
	if err != nil {
		h.handlePathError(w, err)
		return
	}

	log.Log.Debug("Validated path within static directory", "path", path)

	fileToServe := h.resolveFile(path)
	log.Log.Info("serving requested resource for UI", "resource", fileToServe)

	h.setSecurityHeaders(w)
	http.ServeFile(w, r, fileToServe)

	log.Log.Debug("Exiting ServeHTTP")
}

func (h webuiHandler) validatePath(cleanPath, originalPath string) (string, error) {
	path := filepath.Join(h.staticPath, cleanPath)
	absPath, err := filepath.Abs(path)
	if err != nil {
		log.Log.Error("Unable to retrieve absolute path from HTTP request because ", "error", err)
		return "", &pathError{statusCode: http.StatusBadRequest, message: err.Error()}
	}

	if h.isPathTraversal(absPath) {
		log.Log.Warn("Path traversal attempt detected", "requested", originalPath, "resolved", absPath)
		return "", &pathError{statusCode: http.StatusForbidden, message: "Forbidden"}
	}

	return absPath, nil
}

func (h webuiHandler) isPathTraversal(absPath string) bool {
	rel, err := filepath.Rel(h.staticPath, absPath)
	return err != nil || strings.HasPrefix(rel, "..")
}

func (h webuiHandler) resolveFile(path string) string {
	if h.fileExists(path) {
		return path
	}
	return filepath.Join(h.staticPath, h.indexPath)
}

func (h webuiHandler) fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func (h webuiHandler) setSecurityHeaders(w http.ResponseWriter) {
	if len(config.UIConfig.PortalFrontendURL) > 0 {
		w.Header().Set("X-Frame-Options", "allow-from "+config.UIConfig.PortalFrontendURL)
	}
	w.Header().Set("Content-Security-Policy", "default-src *; style-src 'self' 'unsafe-inline'; script-src 'self' 'unsafe-inline' 'unsafe-eval' http://www.google.com")
}

func (h webuiHandler) handlePathError(w http.ResponseWriter, err error) {
	if pe, ok := err.(*pathError); ok {
		http.Error(w, pe.message, pe.statusCode)
		return
	}
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

func (e *pathError) Error() string {
	return e.message
}
