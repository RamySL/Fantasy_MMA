package handlers

import (
	"encoding/json"
	"net/http"
	"fantasy_mma_backend/database"
)

type Card struct {
						// json friendly naming
	ID         int    `json:"id"`
	ExternalID string `json:"external_id"`
	Title      string `json:"title"`
	Date       string `json:"date"`
	Status     string `json:"status"`
	Completed  bool   `json:"completed"`
	VenueName  string `json:"venue_name"`
	City       string `json:"city"`
	Region     string `json:"region"`
	Country    string `json:"country"`
}

func GetCards(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
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
		http.Error(w, "Erreur lecture cards", http.StatusInternalServerError)
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
			http.Error(w, "Erreur scan card", http.StatusInternalServerError)
			return
		}

		cards = append(cards, card)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cards)
}
