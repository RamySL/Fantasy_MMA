package espn

import (
	"database/sql"
	"fantasy/database"
	"fmt"
	"log"
	"strings"
	"time"
)

/* Fichier Utilisé pour synchroniser régulièrement la base de données.
 */

// Lance la synchronisation de la base de données à partir de l'API
// à l'heure précisée et chaque period unité de temps
func SyncTicker(lunchHour int, period time.Duration) {
	now := time.Now()

	lunchHourDate := time.Date(now.Year(), now.Month(), now.Day(), lunchHour, 0, 0, 0, now.Location())
	// Si on éxecute après lunchHour dans la journée il faut ajusté pour avoir le bon calcul
	if now.After(lunchHourDate) {
		lunchHourDate = lunchHourDate.Add(24 * time.Hour)
	}
	sleepDuration := time.Until(lunchHourDate)
	log.Printf("En attente de %v. Prochaine synchro prévue à : %v\n", sleepDuration, sleepDuration)
	time.Sleep(sleepDuration)
 
	// L'heure voulue est atteinte
	log.Println("Début de la première synchronisation...") 
	if err := Sync(); err != nil { 
		log.Printf("[SyncTicker Erreur] : %v\n", err)
	}
	t := time.NewTicker(period)
	for range t.C {
		log.Println("Début de la synchronisation périodique...")
		if err := Sync(); err != nil {
			log.Printf("[SyncTicker Erreur] : %v\n", err)
		}
	}
}
// Fetch les cartes entières à travers l'api et insert dans la base : les cartes, les combats, les combattant.
//TODO: améliorer pour récup moins de cartes de manière redondante
func Sync() error {
	// Note: important de récupérer la prochaine avant d'actualiser la BDD.
	// Sinon le 'STATUS_SCHEDUDLED' sera écrasé
	nextCard, err := database.ScanCard(database.GetNextCard(database.DB)) 
	if err != nil {
		return err
	}

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
			log.Printf("Event trouvé pour : %s", eventDate)
		}

		//CAREFUL: Events est supposé être singleton
		event := scoreBoard2.Events[0]
		
		venue := emptyVenue()
		if len(event.Venues) != 0 {
			//CAREFUL
			venue = event.Venues[0]
		}
		
		// Insertion de carte
		card, err := database.UpsertCard(
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
		err = syncFights(tx, card, event.Competitions)
		if err != nil {
			return err
		}
	}

	// Cohérence de la base par rapport aux combats annulé
	cardsRows, err := database.GetCards(database.DB)
	if err != nil {
		//TODO: Faire partir tout de la tx ?
		return err
	}
	for cardsRows.Next() {
		card, err := database.ScanCard(cardsRows)
		if err != nil {
			log.Printf("[sync.Sync] : cardsRows Scan : %v", err)
		}

		assureCardCoherent(card)
	}
	
	// MAJ des scores par rapports aux prédicitions
	//TODO: erreur non catch
	syncPredictions(nextCard)

	return tx.Commit()
}

// Fetch les combats à travers l'api et remplie les tables `fighters` et `fights`
func syncFights(tx *sql.Tx, card database.Card, competitions []ESPNCompetition) error {

	for _, competition := range competitions {
		// Insertion combattant num 1
		c1 := competition.Competitors[0];
		c2 := competition.Competitors[1];  

		athlete1, err := fetchAthlete(c1.ID)
		if (err != nil) { return err }
		athlete2, err := fetchAthlete(c2.ID)
		if (err != nil) { return err }
	

		if 	strings.Contains(athlete1.FullName, TbaFighterName) || 
			strings.Contains(athlete2.FullName, TbaFighterName) {

			continue
		}

		c1ID, err := upsertFighter(tx, athlete1, c1.Records)
		if err != nil {
			return err
		}
		c2ID, err := upsertFighter(tx, athlete2, c2.Records)
		if err != nil {
			return err
		}

		// NOTE: Carte terminée && Combat non terminé => combat annulé.
		//  l'API ne donne pas explicitement un STATUS_CANCELED
		if !competition.Status.Type.Completed && card.Completed {
			competition.Status.Type.Name = "STATUS_CANCELED"
		}

		_, err = database.UpsertFight(
			tx, 
			competition.ID, 
			card.ID, 
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
//TODO: à simplifier en stockant la prochaine carte ?
//@precondition: BD traité par rapport aux combats annulés
func syncPredictions(nextCard database.Card) error{

	event, err := fetchEvent(nextCard.ExternalID)
	if err != nil {
		return err
	}
	// On regarde si la carte est terminée
	if event.Status.Type.Completed {
		
		var fightID, pointsGoodPrediction int 
		var fightExternalID, status string
		// TODO : délègue la requête à database
		fightRows, err := database.DB.Query(`
			SELECT 
				fights.id,
				fights.external_id,
				fights.points_good_prediction,
				fights.status
				
			FROM fights
			WHERE card_id = $1;
		`,
		nextCard.ID,
		)
		if err != nil {
			return err
		}
		defer fightRows.Close()

		for fightRows.Next(){
			if err := fightRows.Scan(&fightID, &fightExternalID, &pointsGoodPrediction, &status); err != nil {
				return err
			}
			if status == "STATUS_CANCELED" {
				database.DelPredictionByFightID(database.DB, fightID)
				continue
			}
			
			officialWinnerExternalID, ok := getWinnerByCompetID(event.Competitions, fightExternalID)
			if !ok {
				log.Println("[sync.syncPredictions : !ok]")
				continue
			}
			// Toutes les prédictions sur le combat 'fight_id' sont màj
			_, err := database.DB.Exec(`
				UPDATE predictions p
					SET points_obtained = $1
				FROM 	fighters f
				WHERE 	p.fight_id = $2 
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

// L'api d'après tests peut laisser des combats annulés pendant certains jours en tant que "STATUS_SCHEDULED"
// mais aussi peut les retirer après complètement, auquel cas l fonction syncFights() n'aura pas d'effet.
// D'où le rajout de cette fonction
func assureCardCoherent(card database.Card){

	fightsRows, err := database.GetCardFights(database.DB, card.ID)

	if err != nil {
		log.Printf("[sync.makeCardsCoherent] : Erreur dans la récupération de carte %v", err)
	}else{
		for fightsRows.Next() {
			fight, err := database.ScanFight(fightsRows)
			if err != nil {
				log.Printf("sync.makeCardsCoherent : Erreur Scan %v", err)
				continue
			}
			// NOTE: Carte terminée && Combat non terminé => combat annulé.
			//  l'API ne donne pas explicitement un STATUS_CANCELED
			if !fight.Completed && card.Completed {
				fight.Status = "STATUS_CANCELED"
				
				database.UpsertFight(
					database.DB, 
					fight.ExternalID, 
					card.ID, 
					fight.Fighter1.ID,
					fight.Fighter2.ID,
					fight.Winner,
					fight.Category,
					fight.Status,
					fight.Completed,
					10, //TODO: faire une logique pour les points à gagner
					)
			}

		}
	}

}
