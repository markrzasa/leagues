package main

import (
	"net/http"
	"os"

	//	"go.mongodb.org/mongo-driver/bson"
	//	"go.mongodb.org/mongo-driver/bson/primitive"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/kamva/mgm/v3"

	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/markrzasa/leagues/rest"
)

func main() {
	err := mgm.SetDefaultConfig(nil, "leagues", options.Client().ApplyURI(os.Getenv("MONGO_URL")))
	if err != nil {
		panic(err)
	}

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)	

	r.Mount(rest.LeagueURI, rest.LeagueRouter())
	r.Mount(rest.TeamsURI, rest.TeamRouter())
	http.ListenAndServe(":3000", r)
}
