package espn

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const siteBaseURL = "https://site.api.espn.com/apis/site/v2/sports/mma/ufc"
const coreBaseURL = "https://sports.core.api.espn.com/v2/sports/mma"
const leagueUFC = "/leagues/ufc"

// Nom de combatants que peut renvoyer l'API pour des cartes qui n'ont pas été finalisées.
// TBA : to be assigned
const TbaFighterName = "TBA"

// Récupere le scoreboard ESPN general de l'UFC.
func Fetch() (ESPNScoreboardResponse, error) {
	return fetchScoreboard(siteBaseURL + "/scoreboard")
}

// Récupère le scoreboard ESPN pour une date precise.
// Format attendu par ESPN : YYYYMMDD, par exemple "20260502".
func FetchByDate(date string) (ESPNScoreboardResponse, error) {
	return fetchScoreboard(siteBaseURL + "/scoreboard?dates=" + date)
}

func fetchScoreboard(endpoint string) (ESPNScoreboardResponse, error) {
	var scoreboard ESPNScoreboardResponse
	resp, err := http.Get(endpoint)
	if err != nil {
		return scoreboard, fmt.Errorf("appel ESPN impossible: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return scoreboard, fmt.Errorf("erreur ESPN: status HTTP %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(&scoreboard); err != nil {
		return scoreboard, fmt.Errorf("decodage JSON ESPN impossible: %w", err)
	}

	return scoreboard, nil
}

func fetchEvent (id string) (ESPNEvent, error){

	var event ESPNEvent
	url := coreBaseURL + leagueUFC + "/events/" + id
	resp, err := http.Get(url)
	if err != nil {
		return event, fmt.Errorf("appel ESPN impossible: %w", err)
	}
	
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return event, fmt.Errorf("erreur ESPN: status HTTP %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(&event); err != nil {
		return event, fmt.Errorf("decodage JSON ESPN impossible: %w", err)
	}

	return event, nil
}

func fetchAthlete (external_id string) (ESPNAthlete, error){

	url := coreBaseURL + leagueUFC  + "/athletes/" + external_id
	var athlete ESPNAthlete
	resp, err := http.Get(url)
	// TODO: factorise la gestion d'erreurs
	if err != nil {
		return athlete, fmt.Errorf("appel ESPN impossible: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return athlete, fmt.Errorf("erreur ESPN: status HTTP %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(&athlete); err != nil{
		return athlete, fmt.Errorf("decodage JSON ESPN impossible: %w", err)
	}

	return athlete, nil

}
