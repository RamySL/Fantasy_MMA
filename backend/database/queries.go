package database

import (
	"database/sql"
	"time"
)

/*
Des fonctions qui font des insertions dans la base si les données n'existaient
pas avant, sinon mettent à jour la ligne sans insertion de nouvelle entrée dans la table.
*/

func UpsertCard(
	queryMaker QueryMaker,
	externalID string,
	title string,
	date string,
	status string,
	completed bool,
	venueName string,
	city string,
	region string,
	country string,
) (Card, error) {

	var card Card

	query := `
		INSERT INTO cards (
			external_id, title, date, status, completed,
			venue_name, city, region, country
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT (external_id)
		DO UPDATE SET
			title = EXCLUDED.title,
			date = EXCLUDED.date,
			status = EXCLUDED.status,
			completed = EXCLUDED.completed,
			venue_name = EXCLUDED.venue_name,
			city = EXCLUDED.city,
			region = EXCLUDED.region,
			country = EXCLUDED.country
		RETURNING *;
	`
	
	err := queryMaker.QueryRow(
		query,
		externalID,
		title,
		date,
		status,
		completed,
		venueName,
		city,
		region,
		country,
	).Scan(
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

	return card, err
}

func UpsertFighter(
	queryMaker QueryMaker,
	externalID string,
	fullName string,
	record string,
	country_name string,
	country_flag string,
	fighter_image string,
) (int, error) {
	var fighterID int

	query := `
		INSERT INTO fighters (
			external_id, full_name, record, country_name, country_flag, fighter_image
		)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (external_id)
		DO UPDATE SET
			full_name = EXCLUDED.full_name,
			record = EXCLUDED.record,
			country_name = EXCLUDED.country_name,
			country_flag = EXCLUDED.country_flag,
			fighter_image = EXCLUDED.fighter_image
		RETURNING id;
	`

	err := queryMaker.QueryRow(
		query,
		externalID,
		fullName,
		record,
		country_name, 
		country_flag, 
		fighter_image,
	).Scan(&fighterID)

	return fighterID, err
}

func UpsertFight(
	queryMaker QueryMaker,
	externalID string,
	cardID int,
	fighter1ID int,
	fighter2ID int,
	winnerFighterID sql.NullInt64,
	category string,
	status string,
	completed bool,
	pointsGoodPrediction int,
) (int, error) {
	var fightID int

	query := `
		INSERT INTO fights (
			external_id, card_id, fighter1_id, fighter2_id,winner_fighter_id,
			category, status, completed, points_good_prediction
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT (external_id)
		DO UPDATE SET
			card_id = EXCLUDED.card_id,
			fighter1_id = EXCLUDED.fighter1_id,
			fighter2_id = EXCLUDED.fighter2_id,
			winner_fighter_id = EXCLUDED.winner_fighter_id,
			category = EXCLUDED.category,
			status = EXCLUDED.status,
			completed = EXCLUDED.completed,
			points_good_prediction = EXCLUDED.points_good_prediction
		RETURNING id;
	`

	err := queryMaker.QueryRow(
		query,
		externalID,
		cardID,
		fighter1ID,
		fighter2ID,
		winnerFighterID,
		category,
		status,
		completed,
		pointsGoodPrediction,
	).Scan(&fightID)
	
	return fightID, err
}


/** 
 Cards
*/

func GetCards(queryMaker QueryMaker) (*sql.Rows, error) {
	return queryMaker.Query(`
		SELECT id, external_id, title, date::text, status, completed,
		       COALESCE(venue_name, ''), COALESCE(city, ''), 
		       COALESCE(region, ''), COALESCE(country, '')
		FROM cards
		ORDER BY date ASC
	`)
}

func GetCardByID(queryMaker QueryMaker, id int) (*sql.Row) {
	return queryMaker.QueryRow(`
		SELECT id, external_id, title, date::text, status, completed,
		       COALESCE(venue_name, ''), COALESCE(city, ''), 
		       COALESCE(region, ''), COALESCE(country, '')
		FROM cards
		WHERE id = $1;
	`,
	id,
	)
}

func GetCardFights(queryMaker QueryMaker, id int) (*sql.Rows, error){
	return queryMaker.Query(`
		SELECT 
		    fights.id,
			fights.external_id,
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
}

func GetNextCard(qM QueryMaker) (*sql.Row) {
	return qM.QueryRow(`
		SELECT id, external_id, title, date::text, status, completed,
		       COALESCE(venue_name, ''), COALESCE(city, ''), 
		       COALESCE(region, ''), COALESCE(country, '')
		 FROM cards WHERE status='STATUS_SCHEDULED' ORDER BY date ASC LIMIT 1
	`)
}

/**
Predictions
*/

func DelPredictionByFightID(qM QueryMaker, fightID int) {

	qM.Exec(`
		DELETE FROM predictions
		WHERE predictions.fight_id = $1
	`,
	fightID)
}

/**
Auth
*/

func InsertUser(qM QueryMaker, pseudo string, email string, pwd string)(*sql.Row){
	return qM.QueryRow(`
		INSERT INTO users (pseudo, email, password_hash)
		VALUES ($1, $2, $3)
		RETURNING id, pseudo, email
	`,
	pseudo,
	email,
	pwd,
	)
}

func GetUserByEmail(qM QueryMaker, email string)(*sql.Row){
	return qM.QueryRow(`
		SELECT id, pseudo, email, password_hash
		FROM users
		WHERE email = $1
	`, 
	email)
}

func GetUserByID(qM QueryMaker, id int)(*sql.Row){
	return qM.QueryRow(`
		SELECT id, pseudo, email
		FROM users
		WHERE id = $1
	`, id)
}

func InsertSession(qM QueryMaker, userID int, tokenHash string, expiresAt time.Time)(sql.Result, error){
	return qM.Exec(`
		INSERT INTO sessions (user_id, token_hash, expires_at)
		VALUES ($1, $2, $3)
	`, userID, tokenHash, expiresAt)
}

func GetSessionByTokenH(qM QueryMaker, tokenHash string)(*sql.Row){
	return qM.QueryRow(`
		SELECT user_id
		FROM sessions
		WHERE token_hash = $1
	`, tokenHash)
}

func GetValidSessionByTokenH(qM QueryMaker, tokenHash string)(*sql.Row){
	return qM.QueryRow(`
		SELECT user_id
		FROM sessions
		WHERE token_hash = $1
		  AND expires_at > NOW()
	`, tokenHash)
}

func DeleteSessionByTokenH(qM QueryMaker, tokenHash string)(sql.Result, error){
	return qM.Exec(`
		DELETE FROM sessions
		WHERE token_hash = $1
	`, tokenHash)
}
