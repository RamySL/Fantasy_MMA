package handlers

import (
	"database/sql" // TODO: séparation de tâches : passer que par fantasy/database
	"encoding/json"
	"fantasy/database"
	"log"
	"net/http"
	"strings"
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

	var fighter1ID int
	var fighter2ID int
	var fightCompleted bool
	var cardCompleted bool

	err = database.DB.QueryRow(`
		SELECT
			fights.fighter1_id,
			fights.fighter2_id,
			fights.completed,
			cards.completed
		FROM fights
		JOIN cards ON cards.id = fights.card_id
		WHERE fights.id = $1
	`, body.FightID).Scan(
		&fighter1ID,
		&fighter2ID,
		&fightCompleted,
		&cardCompleted,
	)

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

	if fightCompleted || cardCompleted {
		writeJsonError(w, http.StatusConflict, "Impossible de modifier une prédiction sur un combat terminé")
		return
	}

	if body.PredictedWinnerID != fighter1ID && body.PredictedWinnerID != fighter2ID {
		writeJsonError(w, http.StatusBadRequest, "Le combattant choisi ne participe pas à ce combat")
		return
	}


	// Grâce au UNIQUE (user_id, fight_id), cette requête crée la prédiction
	// si elle n'existe pas encore, ou met à jour le choix de l'utilisateur sinon.
	_, err = database.DB.Exec(`
		INSERT INTO predictions (
			user_id,
			fight_id,
			predicted_winner_id,
			points_obtained
		)
		VALUES ($1, $2, $3, 0)
		ON CONFLICT (user_id, fight_id)
		DO UPDATE SET
			predicted_winner_id = EXCLUDED.predicted_winner_id,
			points_obtained = 0
	`,
	// FIXME: (VALUES ($1, $2, $3, 0)) quand tu modifies la base change le 0 pour un type nullable.
		userID,
		body.FightID,
		body.PredictedWinnerID,
	)

	if err != nil {
		log.Printf("predictions.createOrUpdatePrediction: erreur upsert prédiction: %v", err)
		writeJsonError(w, http.StatusInternalServerError, "Erreur sauvegarde prédiction")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Pour toutes les cartes dont l'utilisateur connecté possède des prédictions, retourne ces prédictions
func getMyPredictions(w http.ResponseWriter, req *http.Request) {
	userID, ok := getAuthenticatedUserID(w, req)
	if !ok {
		return
	}

	rows, err := database.DB.Query(`
		SELECT
			cards.id,
			cards.title,
			cards.date::text,
			cards.status,
			cards.completed,

			fights.id,
			COALESCE(fights.category, ''),
			fights.completed,
			fights.points_good_prediction,

			fighter1.full_name,
			fighter2.full_name,

			predictions.predicted_winner_id,
			predicted_fighter.full_name,

			fights.winner_fighter_id,
			winner_fighter.full_name,

			predictions.points_obtained

		FROM predictions
		JOIN fights ON fights.id = predictions.fight_id
		JOIN cards ON cards.id = fights.card_id

		JOIN fighters fighter1 ON fighter1.id = fights.fighter1_id
		JOIN fighters fighter2 ON fighter2.id = fights.fighter2_id
		JOIN fighters predicted_fighter ON predicted_fighter.id = predictions.predicted_winner_id
		LEFT JOIN fighters winner_fighter ON winner_fighter.id = fights.winner_fighter_id

		WHERE predictions.user_id = $1
		ORDER BY cards.date DESC, fights.id ASC
	`, userID)

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
			// NOTE: points obtained sera toujours à 0, ce n'est pas encore màj
			&pointsObtained,
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
			FightID:         fightID,
			Category:        category,
			Fighter1:         fighter1Name,
			Fighter2:         fighter2Name,
			PredictedWinner: predictedWinnerName,
			OfficialWinner:  officialWinner,
			PointsObtained:  pointsObtainedResp,
			PointsGoodPrediction: pointsGoodPrediction,
			Status:          fightStatus,
		})
	}

	if err := rows.Err(); err != nil {
		log.Printf("predictions.getMyPredictions: erreur rows: %v", err)
		writeJsonError(w, http.StatusInternalServerError, "Erreur serveur interne")
		return
	}

	// On garde l'ordre SQL, car une map Go ne garantit pas l'ordre
	response := []MyPredictionCardResponse{}
	for _, cardID := range cardsOrder {
		response = append(response, *cardsByID[cardID])
	}

	writeJsonResponse(w, http.StatusOK, response)
}

// Retourne les prédictions de l'utilisateur connecté pour une carte précise.
func GetCardPredictionsMe(w http.ResponseWriter, req *http.Request, cardID int) {
	if req.Method != http.MethodGet {
		writeJsonError(w, http.StatusMethodNotAllowed, "Méthode non autorisée")
		return
	}

	userID, ok := getAuthenticatedUserID(w, req)
	if !ok {
		return
	}

	rows, err := database.DB.Query(`
		SELECT
			predictions.id,
			predictions.fight_id,
			predictions.predicted_winner_id,
			predictions.points_obtained
		FROM predictions
		JOIN fights ON fights.id = predictions.fight_id
		WHERE predictions.user_id = $1
		  AND fights.card_id = $2
		ORDER BY fights.id ASC
	`, userID, cardID)

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