package handlers

import "database/sql"

/*
	- Type pour passage entre BD -> Json pour les requêtes.
*/

// TODO: à mettre dans database ?
type Card struct {
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

type Fighter struct {
    ID       int    `json:"id"`
    FullName string `json:"full_name"`
    Record   string `json:"record"`
}

type Fight struct {
    ID                   int            `json:"id"`
    Category             string         `json:"category"`
    Status               string         `json:"status"`
    Completed            bool           `json:"completed"`
    PointsGoodPrediction int            `json:"points_good_prediction"`
    Fighter1             Fighter        `json:"fighter1"`
    Fighter2             Fighter        `json:"fighter2"`
    Winner               sql.NullInt64 	 `json:"winner_id"`
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
