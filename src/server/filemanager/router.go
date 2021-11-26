package filemanager

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type Handler interface {
	ServeHTTP(*http.Request) (interface{}, error)
}

type HandlerFunc func(*http.Request) (interface{}, error)

func (hf HandlerFunc) ServeHTTP(r *http.Request) (interface{}, error) {
	return hf(r)
}

type controller struct {
	handler Handler
}

func (ctr controller) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	res, err := ctr.handler.ServeHTTP(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(data); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

type Router struct {
	RouteHandler *mux.Router
}

func NewRouter() *Router {
	return &Router{
		RouteHandler: mux.NewRouter(),
	}
}

func (r *Router) Register(method string, url string, handler Handler) {
	c := controller{
		handler: handler,
	}
	r.RouteHandler.Handle(url, c).Methods(method)
}
