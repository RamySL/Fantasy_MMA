package espn

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const baseURL = "https://site.api.espn.com/apis/site/v2/sports/mma/ufc"
const leaguesBaseURL = "https://sports.core.api.espn.com/v2/sports/mma/leagues/ufc"
// Nom de combatants que peut renvoyer l'API pour des cartes qui n'ont pas été finalisées.
// TBA : to be assigned
const TbaFighterName = "TBA"

// Récupere le scoreboard ESPN general de l'UFC.
func Fetch() (ESPNScoreboardResponse, error) {
	return fetchScoreboard(baseURL + "/scoreboard")
}

// FetchByDate recupere le scoreboard ESPN pour une date precise.
// Format attendu par ESPN : YYYYMMDD, par exemple "20260502".
func FetchByDate(date string) (ESPNScoreboardResponse, error) {
	return fetchScoreboard(baseURL + "/scoreboard?dates=" + date)
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
	resp, err := http.Get(leaguesBaseURL + "/events/" + id)
	fmt.Println("endpoint fetch event : " + baseURL + "/events/" + id )
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
