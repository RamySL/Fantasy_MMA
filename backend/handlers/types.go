package handlers

import "database/sql"

/*
	- Type pour passage entre BD -> Json pour les requêtes.
*/

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