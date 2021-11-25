package route

import (
	"fmt"
	"html"
	"net/http"

	"github.com/gorilla/mux"
)

// Routes registers all the application routes.
func Routes() http.Handler {
	router := mux.NewRouter()
	rh := new(routeHandler)
	rh1 := new(routeHandler1)
	router.Handle("/", rh).Methods(http.MethodPost)
	router.Handle("/", rh1).Methods(http.MethodGet)
	return router
}

type routeHandler struct{}

func (rh routeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
}

type routeHandler1 struct{}

func (rh routeHandler1) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
}
