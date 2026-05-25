package handlers

import (
	"database/sql" // TODO: enlève la dépendance à sql (passage par fantasy/database)
	"encoding/json"
	"fantasy/database"
	"fantasy/espn"
	"log"
	"net/http"
	"strings"
	"time"
)

func Predictions(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/predictions" {
		http.NotFound(w, req)
		return
	}

	if req.Method != http.MethodPost {
		writeJsonError(w, http.StatusMethodNotAllowed, "Méthode non autorisée")
		return
	}

	createOrUpdatePrediction(w, req)
}

func PredictionsRoutes(w http.ResponseWriter, req *http.Request) {
	path := strings.Trim(req.URL.Path, "/")
	parts := strings.Split(path, "/")

	// predictions/me
	if len(parts) == 2 && parts[1] == "me" {
		if req.Method != http.MethodGet {
			writeJsonError(w, http.StatusMethodNotAllowed, "Méthode non autorisée")
			return
		}

		getMyPredictions(w, req)
		return
	}

	http.NotFound(w, req)
}

// Reçoit une liste de prédictions et les crée/met à jour en une seule requête.
// Si une prédiction est invalide (combat terminé, prédictions fermées, combattant absent…),
// toute la requête est rejetée
func PredictionsBulk(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		writeJsonError(w, http.StatusMethodNotAllowed, "Méthode non autorisée")
		return
	}

	userID, ok := getAuthenticatedUserID(w, req)
	if !ok {
		return
	}

	var body BulkPredictionRequestBody
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		writeJsonError(w, http.StatusBadRequest, "JSON invalide")
		return
	}

	if len(body.Predictions) == 0 {
		writeJsonError(w, http.StatusBadRequest, "La liste de prédictions est vide")
		return
	}

	// validation de toutes les prédictions avant les écritures
	type validatedPrediction struct {
		fightID           int
		predictedWinnerID int
	}
	validated := make([]validatedPrediction, 0, len(body.Predictions))

	for _, p := range body.Predictions {
		if p.PredictedWinnerID <= 0 {
			writeJsonError(w, http.StatusBadRequest, "predicted_winner_id invalide")
			return
		}

		var fighter1ID, fighter2ID int
		var fightCompleted, cardCompleted bool
		var cardDate string

		err := database.GetFightAndCardStatus(database.DB, p.FightID).Scan(
			&fighter1ID,
			&fighter2ID,
			&fightCompleted,
			&cardCompleted,
			&cardDate,
		)
		if err != nil {
			if err == sql.ErrNoRows {
				writeJsonError(w, http.StatusNotFound, "Combat introuvable")
				return
			}
			log.Printf("predictions.PredictionsBulk: erreur lecture combat %d: %v", p.FightID, err)
			writeJsonError(w, http.StatusInternalServerError, "Erreur serveur interne")
			return
		}

		parsedCardDate, _ := time.Parse(espn.DateLayout, espn.GetDayDelta(cardDate, 0))
		if time.Now().After(parsedCardDate.AddDate(0, 0, -closePredictionsDeadline)) {
			writeJsonError(w, http.StatusConflict, "Prédictions fermées")
			return
		}

		if fightCompleted || cardCompleted {
			writeJsonError(w, http.StatusConflict, "Impossible de modifier une prédiction sur un combat terminé")
			return
		}

		if p.PredictedWinnerID != fighter1ID && p.PredictedWinnerID != fighter2ID {
			writeJsonError(w, http.StatusBadRequest, "Le combattant choisi ne participe pas à ce combat")
			return
		}

		validated = append(validated, validatedPrediction{
			fightID:           p.FightID,
			predictedWinnerID: p.PredictedWinnerID,
		})
	}

	// Écriture en base
	for _, p := range validated {
		if _, err := database.UpsertPrediction(database.DB, userID, p.fightID, p.predictedWinnerID); err != nil {
			log.Printf("[predictions.PredictionsBulk]: erreur upsert fight %d: %v", p.fightID, err)
			writeJsonError(w, http.StatusInternalServerError, "Erreur sauvegarde prédiction")
			return
		}
	}

	writeJsonResponse(w, http.StatusOK, BulkPredictionResult{SavedCount: len(validated)})
}

func createOrUpdatePrediction(w http.ResponseWriter, req *http.Request) {
	userID, ok := getAuthenticatedUserID(w, req)
	if !ok {
		return
	}

	var body PredictionRequestBody
	err := json.NewDecoder(req.Body).Decode(&body)
	if err != nil {
		writeJsonError(w, http.StatusBadRequest, "JSON invalide")
		return
	}

	if body.PredictedWinnerID <= 0 {
		writeJsonError(w, http.StatusBadRequest, "predicted_winner_id invalide")
		return
	}

	var fighter1ID, fighter2ID int
	var fightCompleted, cardCompleted bool
	var cardDate string

	err = database.GetFightAndCardStatus(database.DB, body.FightID).Scan(
		&fighter1ID,
		&fighter2ID,
		&fightCompleted,
		&cardCompleted,
		&cardDate,
	)
	parsedCardDate, _ := time.Parse(espn.DateLayout, espn.GetDayDelta(cardDate, 0))

	// Note : le race qui peut arriver avec le fait que la carte soit terminée juste après cette reqête
	// est évité parceque on ferme les prédiction avant que la carte soit terminé
	if err != nil {
		if err == sql.ErrNoRows {
			writeJsonError(w, http.StatusNotFound, "Combat introuvable")
			return
		}

		log.Printf("predictions.createOrUpdatePrediction: erreur lecture combat: %v", err)
		writeJsonError(w, http.StatusInternalServerError, "Erreur serveur interne")
		return
	}

	// Les prédictions se ferment un jour avant les combats
	if time.Now().After(parsedCardDate.AddDate(0, 0, -closePredictionsDeadline)) {
		writeJsonError(w, http.StatusConflict, "Prédictions fermées")
		return
	}

	if fightCompleted || cardCompleted {
		writeJsonError(w, http.StatusConflict, "Impossible de modifier une prédiction sur un combat terminé")
		return
	}

	if body.PredictedWinnerID != fighter1ID && body.PredictedWinnerID != fighter2ID {
		writeJsonError(w, http.StatusBadRequest, "Le combattant choisi ne participe pas à ce combat")
		return
	}

	_, err = database.UpsertPrediction(database.DB, userID, body.FightID, body.PredictedWinnerID)

	if err != nil {
		log.Printf("predictions.createOrUpdatePrediction: erreur upsert prédiction: %v", err)
		writeJsonError(w, http.StatusInternalServerError, "Erreur sauvegarde prédiction")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func getMyPredictions(w http.ResponseWriter, req *http.Request) {
	userID, ok := getAuthenticatedUserID(w, req)
	if !ok {
		return
	}

	rows, err := database.GetMyPredictions(database.DB, userID)

	if err != nil {
		log.Printf("predictions.getMyPredictions: erreur query: %v", err)
		writeJsonError(w, http.StatusInternalServerError, "Erreur serveur interne")
		return
	}
	defer rows.Close()

	cardsByID := map[int]*MyPredictionCardResponse{}
	cardsOrder := []int{}

	for rows.Next() {
		var cardID int
		var cardTitle string
		var cardDate string
		var cardStatus string
		var cardCompleted bool

		var fightID int
		var category string
		var fightCompleted bool
		var pointsGoodPrediction int

		var fighter1Name string
		var fighter2Name string

		var predictedWinnerID int
		var predictedWinnerName string

		var winnerID sql.NullInt64
		var winnerName sql.NullString

		var pointsObtained int

		err := rows.Scan(
			&cardID,
			&cardTitle,
			&cardDate,
			&cardStatus,
			&cardCompleted,

			&fightID,
			&category,
			&fightCompleted,
			&pointsGoodPrediction,

			&fighter1Name,
			&fighter2Name,

			&predictedWinnerID,
			&predictedWinnerName,

			&winnerID,
			&winnerName,

			&pointsObtained, // NOTE: points obtained sera toujours à 0, ce n'est pas encore màj
		)

		if err != nil {
			log.Printf("predictions.getMyPredictions: erreur scan: %v", err)
			writeJsonError(w, http.StatusInternalServerError, "Erreur serveur interne")
			return
		}

		card, exists := cardsByID[cardID]
		if !exists {
			status := cardStatus
			if cardCompleted {
				status = "completed"
			}

			card = &MyPredictionCardResponse{
				ID:        cardID,
				CardTitle: cardTitle,
				CardDate:  cardDate,
				Status:    status,
				Fights:    []MyPredictionFightResponse{},
			}

			cardsByID[cardID] = card
			cardsOrder = append(cardsOrder, cardID)
		}

		fightStatus := "pending"
		var officialWinner *string
		var pointsObtainedResp *int

		if winnerName.Valid {
			winner := winnerName.String
			officialWinner = &winner
		}

		if fightCompleted && winnerID.Valid {
			points := pointsObtained
			pointsObtainedResp = &points

			if predictedWinnerID == int(winnerID.Int64) {
				fightStatus = "correct"
				card.GoodPredictions++
			} else {
				fightStatus = "wrong"
			}

			card.TotalPoints += pointsObtained
		}

		card.PossiblePoints += pointsGoodPrediction
		card.TotalPredictions++

		card.Fights = append(card.Fights, MyPredictionFightResponse{
			FightID:              fightID,
			Category:             category,
			Fighter1:             fighter1Name,
			Fighter2:             fighter2Name,
			PredictedWinner:      predictedWinnerName,
			OfficialWinner:       officialWinner,
			PointsObtained:       pointsObtainedResp,
			PointsGoodPrediction: pointsGoodPrediction,
			Status:               fightStatus,
		})
	}

	if err := rows.Err(); err != nil {
		log.Printf("predictions.getMyPredictions: erreur rows: %v", err)
		writeJsonError(w, http.StatusInternalServerError, "Erreur serveur interne")
		return
	}

	response := []MyPredictionCardResponse{}
	for _, cardID := range cardsOrder {
		response = append(response, *cardsByID[cardID])
	}

	writeJsonResponse(w, http.StatusOK, response)
}

func GetCardPredictionsMe(w http.ResponseWriter, req *http.Request, cardID int) {
	if req.Method != http.MethodGet {
		writeJsonError(w, http.StatusMethodNotAllowed, "Méthode non autorisée")
		return
	}

	userID, ok := getAuthenticatedUserID(w, req)
	if !ok {
		return
	}

	rows, err := database.GetCardPredictionsMe(database.DB, userID, cardID)

	if err != nil {
		log.Printf("predictions.GetCardPredictionsMe: erreur query: %v", err)
		writeJsonError(w, http.StatusInternalServerError, "Erreur serveur interne")
		return
	}
	defer rows.Close()

	predictions := []CardPredictionResponse{}

	for rows.Next() {
		var prediction CardPredictionResponse

		err := rows.Scan(
			&prediction.ID,
			&prediction.FightID,
			&prediction.PredictedWinnerID,
			&prediction.PointsObtained,
		)

		if err != nil {
			log.Printf("predictions.GetCardPredictionsMe: erreur scan: %v", err)
			writeJsonError(w, http.StatusInternalServerError, "Erreur serveur interne")
			return
		}

		predictions = append(predictions, prediction)
	}

	if err := rows.Err(); err != nil {
		log.Printf("predictions.GetCardPredictionsMe: erreur rows: %v", err)
		writeJsonError(w, http.StatusInternalServerError, "Erreur serveur interne")
		return
	}

	writeJsonResponse(w, http.StatusOK, predictions)
}