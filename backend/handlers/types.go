package handlers

import (
	"fantasy/database"
	"time"
)

/* TODO : Pour les autres réponses comme pour getCards il y'a une fonction identité implicite
entre database.types et handlers.types
*/

type CardResponse struct {
	database.Card
	PredictionsClosed bool `json:"predictions_closed"`
	EndPredictionsDate time.Time `json:"end_predictions_date"`
}

/* AUTH */

type RegisterRequestBody struct {
	Pseudo		string		`json:"pseudo"`
	Email		string		`json:"email"`
	Password	string		`json:"password"`
}

type AuthUserResponse struct {
	ID     int    `json:"id"`
	Pseudo string `json:"pseudo"`
	Email  string `json:"email"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

/* PREDICTIONS */

type PredictionRequestBody struct {
	FightID           int `json:"fight_id"`
	PredictedWinnerID int `json:"predicted_winner_id"`
}

type CardPredictionResponse struct {
	ID                int `json:"id"`
	FightID           int `json:"fight_id"`
	PredictedWinnerID int `json:"predicted_winner_id"`
	PointsObtained    int `json:"points_obtained"`
}

type MyPredictionCardResponse struct {
	ID               int                         `json:"id"`
	CardTitle        string                      `json:"card_title"`
	CardDate         string                      `json:"card_date"`
	Status           string                      `json:"status"`
	TotalPoints      int                         `json:"total_points"`
	PossiblePoints   int                         `json:"possible_points"`
	GoodPredictions  int                         `json:"good_predictions"`
	TotalPredictions int                         `json:"total_predictions"`
	Fights           []MyPredictionFightResponse `json:"fights"`
}

type RankingEntryResponse struct {
	Rank             int    `json:"rank"`
	Pseudo           string `json:"pseudo"`
	TotalPoints      int    `json:"total_points"`
	GoodPredictions  int    `json:"good_predictions"`
	TotalPredictions int    `json:"total_predictions"`
}

type MyPredictionFightResponse struct {
	FightID          int     `json:"fight_id"`
	Category         string  `json:"category"`
	Fighter1         string  `json:"fighter1"`
	Fighter2         string  `json:"fighter2"`
	PredictedWinner  string  `json:"predicted_winner"`
	OfficialWinner   *string `json:"official_winner"`
	PointsObtained   *int    `json:"points_obtained"`
	PointsGoodPrediction  int     `json:"points_available"`
	Status           string  `json:"status"`
}

// Nombre de jours avant les combats pour fermer les prédictions
const closePredictionsDeadline = 1
