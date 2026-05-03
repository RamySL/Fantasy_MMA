package handlers

import (
	"encoding/json"
	"net/http"
	"fantasy/database"
	"strings"
)

func GetCards(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(w, "cards.GetCards : Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	rows, err := database.DB.Query(`
		SELECT id, external_id, title, date::text, status, completed,
		       COALESCE(venue_name, ''), COALESCE(city, ''), 
		       COALESCE(region, ''), COALESCE(country, '')
		FROM cards
		ORDER BY date ASC
	`)
	if err != nil {
		http.Error(w, "cards.GetCards : Erreur lecture cards", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	cards := []Card{}

	for rows.Next() {
		var card Card

		err := rows.Scan(
			&card.ID,
			&card.ExternalID,
			&card.Title,
			&card.Date,
			&card.Status,
			&card.Completed,
			&card.VenueName,
			&card.City,
			&card.Region,
			&card.Country,
		)
		if err != nil {
			http.Error(w, "cards.GetCards : Erreur scan card", http.StatusInternalServerError)
			return
		}

		cards = append(cards, card)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cards)
}

func GetCardsRoutes(w http.ResponseWriter, req *http.Request){
	path := strings.Trim(req.URL.Path, "/")
	parts := strings.Split(path, "/")

	// cards/{id}
	if(len(parts) == 2){
		getCardByID(w, req, parts[1])
		return
	}
	// cards/{id}/fights
	if(len(parts) == 3 && parts[2] == "fights"){
		getCardFights(w, req, parts[1])
		return
	}

	http.NotFound(w, req)
}

func getCardByID(w http.ResponseWriter, req *http.Request, id string) {
	if req.Method != http.MethodGet {
		http.Error(w, "cards.getCardByID : Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	row := database.DB.QueryRow(`
		SELECT id, external_id, title, date::text, status, completed,
		       COALESCE(venue_name, ''), COALESCE(city, ''), 
		       COALESCE(region, ''), COALESCE(country, '')
		FROM cards
		WHERE id = $1;
	`,
	id,
	)

	var card Card
	err := row.Scan(
		&card.ID,
		&card.ExternalID,
		&card.Title,
		&card.Date,
		&card.Status,
		&card.Completed,
		&card.VenueName,
		&card.City,
		&card.Region,
		&card.Country,
	)
	if err != nil {
		http.Error(w, "cards.getCardByID : Erreur scan card", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(card)

}

func getCardFights(w http.ResponseWriter, req *http.Request, id string){
	if req.Method != http.MethodGet {
		http.Error(w, "cards.getCardByID : Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	rows, err := database.DB.Query(`
		SELECT 
		    fights.id,
			fights.category,
			fights.status,
			fights.completed,
			fights.points_good_prediction,

			fighter1.id AS fighter1_id,
			fighter1.full_name AS fighter1_full_name,
			fighter1.record AS fighter1_record,

			fighter2.id AS fighter2_id,
			fighter2.full_name AS fighter2_full_name,
			fighter2.record AS fighter2_record,

			fights.winner_fighter_id AS winner_id

		FROM fights
		JOIN fighters fighter1 ON fighter1.id = fights.fighter1_id
		JOIN fighters fighter2 ON fighter2.id = fights.fighter2_id

		WHERE card_id = $1;
	`,
	id,
	)

	if err != nil {
		http.Error(w, "cards.getCardFights : Erreur lecture ", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	fights := []Fight{}
	for rows.Next(){
		var fight Fight

		err := rows.Scan(
			&fight.ID,
			&fight.Category,
			&fight.Status,
			&fight.Completed,
			&fight.PointsGoodPrediction,

			&fight.Fighter1.ID,
			&fight.Fighter1.FullName,
			&fight.Fighter1.Record,

			&fight.Fighter2.ID,
			&fight.Fighter2.FullName,
			&fight.Fighter2.Record,

			&fight.Winner,
		)
		if err != nil {
			http.Error(w, "cards.getCardFights : Erreur scan card", http.StatusInternalServerError)
			return
		}

		fights = append(fights, fight)

	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(fights)

}
