package database

import "database/sql"

/*
	- Type pour passage entre BD -> Json.
*/

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
	Fighter_country_name string			`json:"fighter_country_name"`
	Fighter_country_flag_url string 	`json:"fighter_country_flag_url"`
	Fighter_fighter_image_url string 	`json:"fighter_image_url"`
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