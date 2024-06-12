package main

import (
	"log"
	"net/http"

	"github.com/franchesko/assembly-labyrinth/src/internal/api"
	"github.com/franchesko/assembly-labyrinth/src/internal/config"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	cfg := config.MustLoad()

	muxRouter := mux.NewRouter()
	muxRouter.HandleFunc("/levels", api.GetLevelsHandler(cfg)).Methods("GET")
	muxRouter.HandleFunc("/levels/{level}", api.GetLevelInfoHandler(cfg)).Methods("GET")
	muxRouter.HandleFunc("/levels/{level}", api.GetRunLevelHandler(cfg)).Methods("POST")

	router := cors.Default().Handler(muxRouter)
	log.Fatal(http.ListenAndServe(cfg.Address+":"+cfg.Port, router))
}
