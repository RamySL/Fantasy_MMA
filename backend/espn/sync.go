package espn

import (
	"database/sql"
	"fantasy/database"
	"fantasy/handlers"
	"fmt"
	"strings"
)

/* Utilisé pour synchroniser régulièrement la base de données.
*/

// Fetch les cartes entières à travers l'api et insert dans la base : les cartes, les combats, les combattant.
func Sync() error {

	SyncPredictions()

	tx, err := database.DB.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	scoreBoard, err := Fetch()
	if err != nil {
		return err
	}
	if len(scoreBoard.Leagues) == 0 {
		return fmt.Errorf("Sync: aucune league trouvée dans la réponse ESPN")
	}
	//CAREFUL: Leagues est supposé être singleton
	calendars := scoreBoard.Leagues[0].Calendar
	for _, calendar := range calendars {
		eventDate := calendar.StartDate

		scoreBoard2, err := fetchRightDate(eventDate)
		if err != nil {
			return err
		}

		if len(scoreBoard2.Events) == 0 {
			return fmt.Errorf("Sync: aucun event trouvé pour la date %s", eventDate)
		}else{
			fmt.Println("Event trouvé pour : %s", eventDate)
		}

		//CAREFUL: Events est supposé être singleton
		event := scoreBoard2.Events[0]
		
		venue := emptyVenue()
		if len(event.Venues) != 0 {
			//CAREFUL
			venue = event.Venues[0]
		}
		
		// Insertion de carte
		cardID, err := database.UpsertCard(
			tx, 
			event.ID, 
			event.Name, 
			eventDate, 
			event.Status.Type.Name,
			event.Status.Type.Completed,
			venue.FullName,
			venue.Address.City,
			venue.Address.State,
			venue.Address.Country,
		)
		if err != nil {
			return err
		}
		// Insertion de combats et combattants
		err = syncFights(tx, cardID, event.Competitions)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

// Fetch les combats à travers l'api et remplie les tables `fighters` et `fights`
func syncFights(tx *sql.Tx, cardID int, competitions []ESPNCompetition) error {

	for _, competition := range competitions {
		// Insertion combattant num 1
		c1 := competition.Competitors[0]
		c2 := competition.Competitors[1]

		if 	strings.Contains(c1.Athlete.FullName, TbaFighterName) || 
			strings.Contains(c2.Athlete.FullName, TbaFighterName) {

			continue
		}

		c1ID, err := upsertFighter(tx, c1)
		if err != nil {
			return err
		}
		c2ID, err := upsertFighter(tx, c2)
		if err != nil {
			return err
		}

		_, err = database.UpsertFight(
			tx, 
			competition.ID, 
			cardID, 
			c1ID,
			c2ID,
			getWinner(competition, c1, c1ID, c2ID),
			competition.Type.Abbreviation,
			competition.Status.Type.Name,
			competition.Status.Type.Completed,
			10, //TODO: faire une logique pour les points à gagner
		)
		if err != nil {
			return err
		}
	}
	return nil
}

// Pour la prochaine carte qui a lieu, màj la table 'predictions' si la carte est terminée.
//TODO: à simplifier avec un fichier à coté pour stocker la prochaine carte ?
func SyncPredictions() error{

	var card handlers.Card
	// Prochaine carte non terminée dans la base
	err := database.DB.QueryRow(`
		SELECT id, external_id, title, date::text, status, completed,
		       COALESCE(venue_name, ''), COALESCE(city, ''), 
		       COALESCE(region, ''), COALESCE(country, '')
		 FROM cards WHERE status='STATUS_FINAL' ORDER BY date DESC LIMIT 1
	`).Scan(
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
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		return err
	}

	event, err := fetchEvent(card.ExternalID)
	if err != nil {
		return err
	}
	// On regarde si la carte est terminée
	if event.Status.Type.Completed {

		var fightID, pointsGoodPrediction int 
		var fightExternalID string
		// A cette étape la base n'est pas màj avec winner_fighter_id
		fightRows, err := database.DB.Query(`
			SELECT 
				fights.id,
				fights.external_id,
				fights.points_good_prediction
				
			FROM fights
			WHERE card_id = $1;
		`,
		card.ID,
		)
		if err != nil {
			return err
		}
		defer fightRows.Close()

		for fightRows.Next(){
			if err := fightRows.Scan(&fightID, &fightExternalID, &pointsGoodPrediction); err != nil {
				return err
			}
			
			officialWinnerExternalID, ok := getWinnerByCompetID(event.Competitions, fightExternalID)
			if !ok {
				continue
			}
			// Toutes les prédictions sur le combat 'fight_id' sont màj
			_, err := database.DB.Exec(`
				UPDATE predictions p
				SET points_obtained = $1
				FROM fighters f
				WHERE p.fight_id = $2 
				AND f.external_id = $3
				AND p.predicted_winner_id = f.id
				;
			`,
			pointsGoodPrediction, fightID, officialWinnerExternalID,
			)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

