package main

import (
	"github.com/ilhamabdlh/go-restapi/collections"
	"github.com/gorilla/handlers"
)

type Route struct{
	Router := mux.NewRouter()
}

func Main(){
	collections.MainProtocols(Router)
	collections.MainStatus(Router)
	collections.MainConfigs(Router)
	collections.MainDescriptors(Router)
	collections.MainStatus(Router)

	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization", ""})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"})
	origin := handlers.AllowedOrigins([]string{"*"})
	http.ListenAndServe(":4001", handlers.CORS(headers, methods, origin)(Router))
}