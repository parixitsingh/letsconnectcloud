package filemanager

import "net/http"

// Routes registers all the application routes.
func Routes() http.Handler {
	router := NewRouter()
	fileManager := NewFileManager()
	router.Register(http.MethodGet, "/listfiles", HandlerFunc(fileManager.ListFiles))
	router.Register(http.MethodPost, "/addfiles", HandlerFunc(fileManager.AddFiles))
	return router.RouteHandler
}
