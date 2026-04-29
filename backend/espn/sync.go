package espn

import (
	"database/sql"
	"fantasy/database"
	"log"
	"time"
	"fmt"
)

/* Utilisé pour synchroniser régulièrement la base de données.

- Chaque carte est réalisée la nuit du samedi à dimanche. Donc par exemple Dimanche
12h serait un bon temps pour actuliser les résultat d'une carte.

- Par exemple les annonces de cartes sont faitent en millieu de semaine. Faire un fetch le jeudi.

*/

// Fetch les cartes entières à travers l'api et insert dans la base : les cartes, les combats, les combattant.
func Sync() error {
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

func upsertFighter (tx *sql.Tx, c ESPNCompetitor) (int, error) {
	return database.UpsertFighter(
		tx,
		c.ID,
		c.Athlete.FullName,
		getRecord(c),
	)
}

//TODO:utils
func getDayDelta(d string, delta int) (string){
	t, err := time.Parse("2006-01-02T15:04Z", d)
	if err != nil {
		log.Printf("Erreur de Parse dans getDayBefore : ", err)
	}

	return t.AddDate(0, 0, delta).Format("20060102")
}

//TODO: utils
// Avec l'api il se peut que les données sont stocké avec un accès à la date
// exacte, le jours d'avant ou d'après, d'après la localisation de l'évènement.
func fetchRightDate(d string) (ESPNScoreboardResponse, error) {
	scoreBoard, err := FetchByDate(getDayDelta(d, 0))
	if len(scoreBoard.Events) != 0 {
		return scoreBoard, err
	}
	// on teste avec le jour d'avant
	scoreBoard, err = FetchByDate(getDayDelta(d, -1))
	if len(scoreBoard.Events) != 0 {
		return scoreBoard, err
	}
	// on teste avec le jour d'avant
	scoreBoard, err = FetchByDate(getDayDelta(d, 1))
	if len(scoreBoard.Events) != 0 {
		return scoreBoard, err
	}

	// on doit pas arriver ici
	return scoreBoard, err
}

//TODO: utils 
func getRecord(c ESPNCompetitor) string {
	if len(c.Records) == 0 {
		return ""
	}

	return c.Records[0].Summary
}

//TODO: utils 
func getWinner(compet ESPNCompetition, c1 ESPNCompetitor, c1ID int, c2ID int) (sql.NullInt64){
	if (compet.Status.Type.Completed){
		if c1.Winner{
			return sql.NullInt64{Valid: true, Int64: int64(c1ID)}
		}else{
			return sql.NullInt64{Valid: true, Int64: int64(c2ID)}
		}
	}else{
		return sql.NullInt64{Valid: false}
	}
}

//TODO: utils
func emptyVenue() ESPNVenue {
	return ESPNVenue{
		FullName: "",
		Address: ESPNAddress{
			City:    "",
			State:   "",
			Country: "",
		},
	}
}