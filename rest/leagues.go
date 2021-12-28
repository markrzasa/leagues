package rest

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/markrzasa/leagues/datastore"
)

const LeagueURI = "/league"

type League struct {
	Name string `json:"name"`
	URI  string `json:"uri"`
}

func uriForLeague(league *datastore.League) string {
	return fmt.Sprintf("%s/%s", LeagueURI, league.ID.Hex())
}

func leagueFromDatastore(league *datastore.League) *League {
	return &League{
		Name: league.Name,
		URI: uriForLeague(league),
	}
}

func leaguesFromDatastore(leagues *[]datastore.League) *[]League {
	var restLeagues []League
	for _, l := range(*leagues) {
		restLeagues = append(restLeagues, *leagueFromDatastore(&l))
	}

	return &restLeagues
}

func listLeagues(w http.ResponseWriter, r *http.Request) {
	leagues, err := datastore.ListLeagues()
	if err == nil {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(leaguesFromDatastore(leagues))
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(NewError(err))
	}
}

func getLeague(w http.ResponseWriter, r *http.Request) {
	league, err := datastore.GetLeague(chi.URLParam(r, "leagueId"))
	if err == nil {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(league)
	} else {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(NewError(err))
	}
}

func postLeague(w http.ResponseWriter, r *http.Request) {
	var league League 
    json.NewDecoder(r.Body).Decode(&league)  
	l, err := datastore.NewLeague(league.Name)
	league.URI = uriForLeague(l)
	if err == nil {
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(leagueFromDatastore(l))
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(NewError(err))
	}
}

func deleteLeague(w http.ResponseWriter, r *http.Request) {
	err := datastore.DeleteLeague(chi.URLParam(r, "leagueId"))
	if err == nil {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(NewError(err))
	}
}

func LeagueRouter() http.Handler {
	r := chi.NewRouter()
	r.Get("/", listLeagues)
	r.Get("/{leagueId}", getLeague)
	r.Get("/{leagueId}/season/{seasonId}", getSeason)
	r.Post("/", postLeague)
	r.Post("/{leagueId}/season", postSeason)
	r.Delete("/{leagueId}", deleteLeague)
	r.Delete("/{leagueId}/season/{seasonId}", deleteSeason)
	return r
}
