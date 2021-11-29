package filemanager

import "net/http"

// Routes registers all the application routes.
func Routes() http.Handler {
	router := NewRouter()
	fileManager := NewFileManager()
	router.Register(http.MethodGet, "/listfiles", HandlerFunc(fileManager.ListFiles))
	router.Register(http.MethodGet, "/wordscount", HandlerFunc(fileManager.WordCounts))
	router.Register(http.MethodGet, "/wordsfrequency", HandlerFunc(fileManager.WordFrequency))
	router.Register(http.MethodPost, "/addfiles", HandlerFunc(fileManager.AddFiles))
	router.Register(http.MethodPut, "/updatefiles", HandlerFunc(fileManager.UpdateFiles))
	router.Register(http.MethodDelete, "/removefile", HandlerFunc(fileManager.RemoveFile))
	return router.RouteHandler
}
