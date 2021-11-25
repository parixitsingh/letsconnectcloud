package route

import (
	"fmt"
	"html"
	"net/http"

	"github.com/gorilla/mux"
)

// RegisterRoutes registers all the application routes.
func RegisterRoutes() http.Handler {
	router := mux.NewRouter()
	rh := new(routeHandler)
	rh1 := new(routeHandler)
	router.Handle("/", rh).Methods(http.MethodPost)
	router.Handle("/", rh1).Methods(http.MethodGet)
	return router
}

type routeHandler struct{}

func (rh routeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
}
