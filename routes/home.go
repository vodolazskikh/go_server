package routes

import (
	"fmt"
	"net/http"
)

// HomeRoute - главный роут
func HomeRoute(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello world")
}
