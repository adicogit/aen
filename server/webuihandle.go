package server

import (
	"net/http"
	"os"
	"path/filepath"
	"regexp"

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
	// get the absolute path to prevent directory traversal
	path, err := filepath.Abs(r.URL.Path)
	if err != nil {
		log.Log.Error("Unable to retrieve absolute path from HTTP request becasue ", "error", err)
		// if we failed to get the absolute path respond with a 400 bad request
		// and stop
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Log.Debug("Retrieve path from HTTP request", "path", path)

	//remove c:\ and d:\ in case of windows platform
	re := regexp.MustCompile(`(?i)c:\\`)
	path = re.ReplaceAllString(path, "\\\\")
	re = regexp.MustCompile(`(?i)d:\\`)
	path = re.ReplaceAllString(path, "\\\\")
	log.Log.Debug("Removing disk letter from path in case of Windows platform", "path", path)
	// prepend the path with the path to the static directory
	path = filepath.Join(h.staticPath, path)
	log.Log.Debug("Prepending the path with the path to the static directory", "path", path)

	// check whether a file exists at the given path
	_, err = os.Stat(path)
	if !os.IsNotExist(err) {
		log.Log.Info("serving requested resource for UI", "resource", path)
		// file does not exist, serve index.html
		//http.ServeFile(w, r, filepath.Join(h.staticPath, h.indexPath))
		http.ServeFile(w, r, path)
		return
	} else if err != nil {
		log.Log.Error("Unable to serve UI becasue ", "error", err)
		// if we got an error (that wasn't that the file doesn't exist) stating the
		// file, return a 500 internal server error and stop
		http.Error(w, "Unable to serve UI", http.StatusNotFound)
		return
	}
	if len(config.UIConfig.PortalFrontendURL) > 0 {
		// PortalFrontendURL property has been specified and we must set X-Frame-Options
		// allowing Portal plugin to use micro frontend
		allowFrom := "allow-from " + config.UIConfig.PortalFrontendURL
		w.Header().Set("X-Frame-Options", allowFrom)
	}
	w.Header().Set("Content-Security-Policy", "default-src *; style-src 'self' 'unsafe-inline'; script-src 'self' 'unsafe-inline' 'unsafe-eval' http://www.google.com")

	// otherwise, use http.FileServer to serve the static dir
	http.FileServer(http.Dir(h.staticPath)).ServeHTTP(w, r)
	log.Log.Debug("Exiting ServeHTTP")
}
