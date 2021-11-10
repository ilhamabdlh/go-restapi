package routers

import (
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

var Route = mux.NewRouter()
func Router(){
	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization", ""})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"})
	origin := handlers.AllowedOrigins([]string{"*"})
	http.ListenAndServe(":4001", handlers.CORS(headers, methods, origin)(Route))
}