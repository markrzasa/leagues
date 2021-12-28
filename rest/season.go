package rest

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/markrzasa/leagues/datastore"
)

type Season struct {
	LeagueURI string `json:"league"`
	Name      string `json:"name"`
	URI       string `json:"uri"`
}

func uriForSeason(leagueId string, season *datastore.Season) (string, string) {
	leagueURI := fmt.Sprintf("%s/%s", LeagueURI, leagueId)
	seasonURI := fmt.Sprintf("%s/season/%s", leagueURI, season.ID.Hex())
	return leagueURI, seasonURI
}

func seasonFromDatastore(leagueId string, season *datastore.Season) *Season {
	leagueURI, seasonURI := uriForSeason(leagueId, season)
	return &Season{
		LeagueURI: leagueURI,
		Name:      season.Name,
		URI:       seasonURI,
	}
}

func postSeason(w http.ResponseWriter, r *http.Request) {
	leagueId := chi.URLParam(r, "leagueId")
	var season Season 
    json.NewDecoder(r.Body).Decode(&season)
	s, err := datastore.NewSeason(season.Name, leagueId)
	if err == nil {
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(seasonFromDatastore(leagueId, s))
	} else {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(NewError(err))
	}
}

func getSeason(w http.ResponseWriter, r *http.Request) {
	leagueId := chi.URLParam(r, "leagueId")
	_, err := datastore.GetLeague(leagueId)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(NewError(err))
	} else {
		s, err := datastore.GetSeason(chi.URLParam(r, "seasonId"))
		if err == nil {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(seasonFromDatastore(leagueId, s))	
		} else {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(NewError(err))	
		}
	}
}

func deleteSeason(w http.ResponseWriter, r *http.Request) {
	leagueId := chi.URLParam(r, "leagueId")
	_, err := datastore.GetLeague(leagueId)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(NewError(err))
	} else {
		err := datastore.DeleteSeason(chi.URLParam(r, "seasonId"))
		if err == nil {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(NewError(err))
		}
	}
}
