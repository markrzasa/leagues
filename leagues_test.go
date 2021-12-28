package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/go-resty/resty/v2"

	"github.com/markrzasa/leagues/rest"
)

func getenv(variable, fallback string) string {
    value := os.Getenv(variable)
    if len(value) == 0 {
        return fallback
    }
    return value
}

func baseURL() string {
	return getenv("BASE_URL", "http://localhost:3000")
}

func assertStatus(t *testing.T, err error, url string, expectedStatus int, resp *resty.Response) {
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode() != expectedStatus {
		t.Errorf("%s: expected: %d, actual: %d", url, expectedStatus, resp.StatusCode())
		t.Fail()
	}
}

func TestSeason(t *testing.T) {
	leaguesURL := fmt.Sprintf("%s%s", baseURL(), rest.LeagueURI)
	client := resty.New()
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(`{"name": "test-league"}`).
		Post(leaguesURL)
	assertStatus(t, err, leaguesURL, http.StatusCreated, resp)
	league := &rest.League{}
	json.Unmarshal(resp.Body(), &league)
	leagueURL := fmt.Sprintf("%s%s", baseURL(), league.URI)
	resp, err = client.R().Get(leagueURL)
	assertStatus(t, err, leaguesURL, http.StatusOK, resp)

	// create a season
	seasonURL := fmt.Sprintf("%s/season", leagueURL)
	resp, err = client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(`{"name": "2021-2022"}`).
		Post(seasonURL)
	assertStatus(t, err, seasonURL, http.StatusCreated, resp)
	season := &rest.Season{}
	json.Unmarshal(resp.Body(), &season)
	seasonURL = fmt.Sprintf("%s%s", baseURL(), season.URI)
	resp, err = client.R().Get(seasonURL)
	assertStatus(t, err, seasonURL, http.StatusOK, resp)
	resp, err = client.R().Delete(seasonURL)
	assertStatus(t, err, seasonURL, http.StatusOK, resp)

	// delete the league used for testing
	t.Cleanup(func() {
		resp, err := client.R().Delete(leagueURL)
		assertStatus(t, err, leagueURL, http.StatusOK, resp)
	})
}
