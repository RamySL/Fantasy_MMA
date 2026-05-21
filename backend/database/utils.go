package database

/** SCANS */

func ScanCard(row RowScanner)(Card, error){

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
	return card, err
}

func ScanFight(rows RowScanner)(Fight, error){
	var fight Fight
	err := rows.Scan(
		&fight.ID,
		&fight.ExternalID,
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

	return fight, err
}