package handlers

import (
	"log"
	"net/http"
	"fantasy/database"
	"strings"
)

func GetCards(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		writeJsonError(w, http.StatusMethodNotAllowed, "Méthode non autorisée")
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
		log.Printf("handlers.GetCards: erreur query cards: %v", err)
		writeJsonError(w, http.StatusInternalServerError, "Erreur serveur interne")
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
			log.Printf("handlers.GetCards: erreur scan card: %v", err)
			writeJsonError(w, http.StatusInternalServerError, "Erreur serveur interne")
			return
		}

		cards = append(cards, card)
	}

	writeJsonResponse(w, http.StatusOK, cards)
}

func GetCardsRoutes(w http.ResponseWriter, req *http.Request) {

	if req.Method != http.MethodGet {
		writeJsonError(w, http.StatusMethodNotAllowed, "Méthode non autorisée")
		return
	}

	path := strings.Trim(req.URL.Path, "/")
	parts := strings.Split(path, "/")

	// cards/{id}
	if len(parts) == 2 {
		getCardByID(w, req, parts[1])
		return
	}

	// cards/{id}/fights
	if len(parts) == 3 && parts[2] == "fights" {
		getCardFights(w, req, parts[1])
		return
	}

	// cards/{id}/predictions/me
	if len(parts) == 4 && parts[2] == "predictions" && parts[3] == "me" {
		GetCardPredictionsMe(w, req, parts[1])
		return
	}

	http.NotFound(w, req)
}

func getCardByID(w http.ResponseWriter, req *http.Request, id string) {

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
		log.Printf("cards.getCardByID : Erreur scan card %v", err)
		writeJsonError(w, http.StatusInternalServerError, "Erreur serveur interne")
		return
	}

	writeJsonResponse(w, http.StatusOK, card)

}

func getCardFights(w http.ResponseWriter, req *http.Request, id string){

	rows, err := database.DB.Query(`
		SELECT 
		    fights.id,
			COALESCE(fights.category, ''),
			COALESCE(fights.status, ''),
			fights.completed,
			fights.points_good_prediction,

			fighter1.id AS fighter1_id,
			fighter1.full_name AS fighter1_full_name,
			fighter1.record AS fighter1_record,
			COALESCE(fighter1.country_name, '') AS fighter1_country_name,
			COALESCE(fighter1.country_flag, '') AS fighter1_country_flag_url,
			COALESCE(fighter1.fighter_image, '') AS fighter1_fighter_image_url,
			

			fighter2.id AS fighter2_id,
			fighter2.full_name AS fighter2_full_name,
			fighter2.record AS fighter2_record,
			COALESCE(fighter2.country_name, '') AS fighter2_country_name,
			COALESCE(fighter2.country_flag, '') AS fighter2_country_flag_url,
			COALESCE(fighter2.fighter_image, '') AS fighter2_fighter_image_url,
			
			fights.winner_fighter_id AS winner_id

		FROM fights
		JOIN fighters fighter1 ON fighter1.id = fights.fighter1_id
		JOIN fighters fighter2 ON fighter2.id = fights.fighter2_id

		WHERE card_id = $1;
	`,
	id,
	)

	if err != nil {
		log.Printf("cards.getCardFights : Erreur lecture %v", err)
		writeJsonError(w, http.StatusInternalServerError, "Erreur serveur interne")
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
			&fight.Fighter1.Fighter_country_name,
			&fight.Fighter1.Fighter_country_flag_url,
			&fight.Fighter1.Fighter_fighter_image_url,

			&fight.Fighter2.ID,
			&fight.Fighter2.FullName,
			&fight.Fighter2.Record,
			&fight.Fighter2.Fighter_country_name,
			&fight.Fighter2.Fighter_country_flag_url,
			&fight.Fighter2.Fighter_fighter_image_url,

			&fight.Winner,
		)
		if err != nil {
			log.Printf("cards.getCardFights : Erreur scan card %v", err)
			writeJsonError(w, http.StatusInternalServerError, "Erreur serveur interne")
			return
		}

		fights = append(fights, fight)

	}

	writeJsonResponse(w, http.StatusOK, fights)

}
