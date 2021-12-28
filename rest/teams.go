package rest

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/markrzasa/leagues/datastore"
)

const TeamsURI = "/team"

type Team struct {
	Name string `json:"name"`
	City string `json:"city"`
	URI  string `json:"uri"`
}

func uriForTeam(team *datastore.Team) string {
	return fmt.Sprintf("%s/%s", TeamsURI, team.ID.Hex())
}

func teamFromDatastore(team *datastore.Team) *Team {
	return &Team{
		Name: team.Name,
		City: team.City,
		URI: uriForTeam(team),
	}
}

func teamsFromDatastore(teams *[]datastore.Team) *[]Team {
	var restTeams []Team
	for _, t := range(*teams) {
		restTeams = append(restTeams, *teamFromDatastore(&t))
	}

	return &restTeams
}

func listTeams(w http.ResponseWriter, r *http.Request) {
	teams, err := datastore.ListTeams()
	if err == nil {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(teamsFromDatastore(teams))
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(NewError(err))
	}
}

func getTeam(w http.ResponseWriter, r *http.Request) {
	team, err := datastore.GetTeam(chi.URLParam(r, "teamId"))
	if err == nil {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(team)
	} else {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(NewError(err))
	}
}

func postTeam(w http.ResponseWriter, r *http.Request) {
	var team Team 
    json.NewDecoder(r.Body).Decode(&team)  
	l, err := datastore.NewTeam(team.Name)
	team.URI = uriForTeam(l)
	if err == nil {
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(l)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(NewError(err))
	}
}

func TeamRouter() http.Handler {
	r := chi.NewRouter()
	r.Get("/", listTeams)
	r.Get("/{teamId}", getTeam)
	r.Post("/", postTeam)
//	r.Delete("/{leagueId}", deleteLeague)
	return r
}
