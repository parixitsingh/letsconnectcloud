package filemanager

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// Handler is representing a custom handler for route
type Handler interface {
	ServeHTTP(*http.Request) (interface{}, error)
}

// HandlerFunc is implementing ServeHTTP
type HandlerFunc func(*http.Request) (interface{}, error)

// HandlerFunc's ServeHTTP
func (hf HandlerFunc) ServeHTTP(r *http.Request) (interface{}, error) {
	return hf(r)
}

// controller is a wrapper holding custom Handler
type controller struct {
	handler Handler
}

// ServeHTTP over controller making it acts http.Handler
func (ctr controller) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	res, err := ctr.handler.ServeHTTP(r)
	if err != nil {
		fmt.Printf("error occured from serveHTTP : %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(res)
	if err != nil {
		fmt.Printf("error occured from marshal : %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(data); err != nil {
		fmt.Printf("error occured from write : %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// Router is wrapping *mux.Router
type Router struct {
	RouteHandler *mux.Router
}

// Constructor creating Router instance
func NewRouter() *Router {
	return &Router{
		RouteHandler: mux.NewRouter(),
	}
}

// Register method is registering routes
func (r *Router) Register(method string, url string, handler Handler) {
	c := controller{
		handler: handler,
	}
	r.RouteHandler.Handle(url, c).Methods(method)
}
