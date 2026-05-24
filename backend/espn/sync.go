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

	log.Printf("Attente de %v avant prochaine synchro de la base. Prochaine synchro prévue à : %v\n", sleepDuration, lunchHourDate)
	time.Sleep(sleepDuration)
 
	// L'heure voulue est atteinte
	if err := UpdateSync(); err != nil { 
		log.Printf("[SyncTicker Erreur] : %v\n", err)
	}
	
	t := time.NewTicker(period)
	for range t.C {
		log.Println("Début de la synchronisation périodique...")
		if err := UpdateSync(); err != nil {
			log.Printf("[SyncTicker Erreur] : %v\n", err)
		}
	}
}

// InitialSync télécharge toutes les cartes de l'année et remplit une base de données vierge.
func InitialSync() error {
	log.Println("[InitialSync] Démarrage de la synchronisation initiale...")

	scoreBoard, err := fetch()
	if err != nil {
		return err
	}
	if len(scoreBoard.Leagues) == 0 {
		return fmt.Errorf("InitialSync: aucune league trouvée dans la réponse ESPN")
	}

	calendars := scoreBoard.Leagues[0].Calendar
	for _, calendar := range calendars {
		eventDate := calendar.StartDate
		
		scoreBoard2, err := fetchRightDate(eventDate)
		if err != nil {
			log.Printf("[InitialSync] Erreur fetch %s: %v\n", eventDate, err)
			continue
		}
		if len(scoreBoard2.Events) == 0 {
			continue
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
		
		log.Printf("Event initialisé : %s\n", eventDate)
	}

	log.Println("[InitialSync] Synchronisation initiale terminée avec succès.")
	return nil
}

// UpdateSync actualise la base de données à partir de la prochaine carte prévue.
// Gère les prédictions et les combats annulés.
func UpdateSync() error {

	nextCard, err := database.ScanCard(database.GetNextCard(database.DB)) 
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("[UpdateSync] Base vide détectée (Erreur ErrNoRows). Lancement d'une InitialSync à la place.")
			return InitialSync()
		}
		return err
	}

	parsedNextDate, err := time.Parse(DateLayout, GetDayDelta(nextCard.Date, 0))
	if err != nil {
		fmt.Println("[UpdateSync] next card parse erreur:", err)
		return err
	}

	scoreBoard, err := fetch()
	if err != nil {
		return err
	}
	if len(scoreBoard.Leagues) == 0 {
		return fmt.Errorf("UpdateSync: aucune league trouvée dans la réponse ESPN")
	}

	calendars := scoreBoard.Leagues[0].Calendar
	for _, calendar := range calendars {

		eventDate := calendar.StartDate
		parsedDate, err := time.Parse(DateLayout, GetDayDelta(eventDate, 0))
		if err != nil {
			fmt.Println("[UpdateSync] parse erreur:", err)
			continue
		}
		
		// On ignore les événements passés
		if parsedDate.Before(parsedNextDate) {
			continue
		}

		scoreBoard2, err := fetchRightDate(eventDate)
		if err != nil {
			return err
		}

		if len(scoreBoard2.Events) == 0 {
			continue
		} else {
			log.Printf("Event mis à jour pour : %s", eventDate)
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

	// 4. Traitement des annulations et prédictions
	assureCardCoherent(nextCard)	
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