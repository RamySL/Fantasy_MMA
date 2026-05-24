package espn

import (
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

// fonction principale d'actualisation des tables de base de données.
// Fetch les prochaines cartes à travers l'api et insert dans la base : les cartes, les combats, les combattant.
// et actualise les résultat de prédictions
func Sync() error {
	// Note: important de récupérer la prochaine avant d'actualiser la BDD.
	// Sinon le 'STATUS_SCHEDUDLED' sera écrasé
	nextCard, err := database.ScanCard(database.GetNextCard(database.DB)) 
	if err != nil {
		return err
	}

	scoreBoard, err := fetch()
	if err != nil {
		return err
	}
	if len(scoreBoard.Leagues) == 0 {
		return fmt.Errorf("Sync: aucune league trouvée dans la réponse ESPN")
	}

    parsedNextDate, err := time.Parse(DateLayout, GetDayDelta(nextCard.Date, 0))
    if err != nil {
        fmt.Println("[sync.sync] next card parse:", err)
        return err
    }
	//CAREFUL: Leagues est supposé être singleton
	calendars := scoreBoard.Leagues[0].Calendar
	// On parcourt les carte de combat pour le clendrier de l'année en cours
	for _, calendar := range calendars {

		eventDate := calendar.StartDate
		parsedDate, err := time.Parse(DateLayout, GetDayDelta(eventDate, 0))
		if err != nil {
			fmt.Println("[sync.sync] parse:", err)
			continue
		}
		// évenement terminé
		if parsedDate.Before(parsedNextDate){
			continue
		}

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
			database.DB, 
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
		err = syncFights(card, event.Competitions)
		if err != nil {
			return err
		}
	}

	// traitement des combats annulés
	assureCardCoherent(nextCard)	

	// MAJ des scores par rapports aux prédicitions
	// Précondition respectée
	return syncPredictions(nextCard)

}

// remplie les tables `fighters` et `fights`.
// Si une des insertions échoue toutes les autres sont annulées.
func syncFight(compet ESPNCompetition, card database.Card) error {
	// Insertion combattant num 1
	c1 := compet.Competitors[0];
	c2 := compet.Competitors[1];  

	athlete1, err := fetchAthlete(c1.ID)
	if (err != nil) { return err }
	athlete2, err := fetchAthlete(c2.ID)
	if (err != nil) { return err }

	// Combat vide 
	if 	strings.Contains(athlete1.FullName, TbaFighterName) || 
		strings.Contains(athlete2.FullName, TbaFighterName) {

		return nil
	}

	tx, err := database.DB.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	c1ID, err := upsertFighter(tx, athlete1, c1.Records)
	if err != nil {
		return err
	}
	c2ID, err := upsertFighter(tx, athlete2, c2.Records)
	if err != nil {
		return err
	}

	_, err = database.UpsertFight(
		tx, 
		compet.ID, 
		card.ID, 
		c1ID,
		c2ID,
		getWinner(compet, c1, c1ID, c2ID),
		compet.Type.Abbreviation,
		compet.Status.Type.Name,
		compet.Status.Type.Completed,
		10, //TODO: faire une logique pour les points à gagner
	)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// Fetch les combats à travers l'api et remplie les tables `fighters` et `fights`
func syncFights(card database.Card, competitions []ESPNCompetition) error {
	// competition ~= fight
	for _, competition := range competitions {
		err := syncFight(competition, card)
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

		fightRows, err := database.GetFightByCardID(database.DB, nextCard.ID)
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
			_, err := database.UpdatePrediction(database.DB, pointsGoodPrediction, fightID, officialWinnerExternalID)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
