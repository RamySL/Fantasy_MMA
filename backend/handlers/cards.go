package handlers

import (
	"fantasy/database"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func GetCards(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		writeJsonError(w, http.StatusMethodNotAllowed, "Méthode non autorisée")
		return
	}

	rows, err := database.GetCards(database.DB)
	if err != nil {
		log.Printf("handlers.GetCards: erreur query cards: %v", err)
		writeJsonError(w, http.StatusInternalServerError, "Erreur serveur interne")
		return
	}
	defer rows.Close()

	cards := []database.Card{}

	for rows.Next() {
		card, err := database.ScanCard(rows)
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

			
	id, err := strconv.Atoi(parts[1])
	if err != nil || id < 0 {
		writeJsonError(w, http.StatusBadRequest, "ID invalide")
	}
	// cards/{id}
	if len(parts) == 2 {
		getCardByID(w, req, id)
		return
	}

	// cards/{id}/fights
	if len(parts) == 3 && parts[2] == "fights" {
		getCardFights(w, req, id)
		return
	}

	// cards/{id}/predictions/me
	if len(parts) == 4 && parts[2] == "predictions" && parts[3] == "me" {
		GetCardPredictionsMe(w, req, id)
		return
	}

	http.NotFound(w, req)
}

func getCardByID(w http.ResponseWriter, req *http.Request, id int) {

	row := database.GetCardByID(database.DB, id)

	var card database.Card
	card, err := database.ScanCard(row)

	if err != nil {
		log.Printf("cards.getCardByID : Erreur scan card %v", err)
		writeJsonError(w, http.StatusInternalServerError, "Erreur serveur interne")
		return
	}

	writeJsonResponse(w, http.StatusOK, card)

}

func getCardFights(w http.ResponseWriter, req *http.Request, id int){

	rows, err := database.GetCardFights(database.DB, id)

	if err != nil {
		log.Printf("cards.getCardFights : Erreur lecture %v", err)
		writeJsonError(w, http.StatusInternalServerError, "Erreur serveur interne")
		return
	}
	defer rows.Close()

	fights := []database.Fight{}
	for rows.Next(){

		fight, err := database.ScanFight(rows)
		if err != nil {
			log.Printf("cards.getCardFights : Erreur scan card %v", err)
			writeJsonError(w, http.StatusInternalServerError, "Erreur serveur interne")
			return
		}
		fights = append(fights, fight)

	}

	writeJsonResponse(w, http.StatusOK, fights)

}
