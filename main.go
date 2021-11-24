package main

import (
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/ilhamabdlh/go-restapi/collections"
	"github.com/ilhamabdlh/go-restapi/helper"
)

var Routers = helper.Routes

func main() {
	collections.MainProtocols()
	collections.MainStatusConfigs()
	collections.MainItems()
	collections.MainDescriptors()

	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization", ""})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT"})
	origin := handlers.AllowedOrigins([]string{"*"})
	http.ListenAndServe(":4001", handlers.CORS(headers, methods, origin)(Routers))
}
