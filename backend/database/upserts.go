package database

import (
	"database/sql"
)

/* 
Ce fichier contient des fonctions qui font des insertions dans la base si les données n'existaient
pas avant, sinon mettent à jour la ligne sans insertion de nouvelle entrée dans la table.
*/

func UpsertCard(
	tx *sql.Tx,
	externalID string,
	title string,
	date string,
	status string,
	completed bool,
	venueName string,
	city string,
	region string,
	country string,
) (int, error) {

	var cardID int

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
		RETURNING id;
	`

	err := tx.QueryRow(
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
	).Scan(&cardID)

	return cardID, err
}

func UpsertFighter(
	tx *sql.Tx,
	externalID string,
	fullName string,
	record string,
) (int, error) {
	var fighterID int

	query := `
		INSERT INTO fighters (
			external_id, full_name, record
		)
		VALUES ($1, $2, $3)
		ON CONFLICT (external_id)
		DO UPDATE SET
			full_name = EXCLUDED.full_name,
			record = EXCLUDED.record
		RETURNING id;
	`

	err := tx.QueryRow(
		query,
		externalID,
		fullName,
		record,
	).Scan(&fighterID)

	return fighterID, err
}

func UpsertFight(
	tx *sql.Tx,
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

	err := tx.QueryRow(
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