package espn

import(
	"database/sql"
	"fantasy/database"
	"time"
	"log"
)

func upsertFighter (tx *sql.Tx, athlete ESPNAthlete, records []ESPNRecord) (int, error) {
	return database.UpsertFighter(
		tx,
		athlete.ID,
		athlete.FullName,
		getRecord(records),
		athlete.Flag.Alt,
		athlete.Flag.Href,
		athlete.HeadShot.Href,
	)
}

// TODO: Rendre générique la gestion de dates. Au lieu de couvrir tous les formats manipulés par le projet.
func GetDayDelta(d string, delta int) (string){
	t, err := time.Parse("2006-01-02T15:04Z", d)
	if err != nil {
		t, err = time.Parse("2006-01-02 15:04:05", d)

		if err != nil {
			// Forcément c'est l'autre format		
			t, err = time.Parse("2006-01-02T15:04:05Z", d)

			if err != nil{
					log.Printf("[utils.GetDayDelta] : Erreur de Parse dans getDayBefore : %v", err)
			}
		}
	}

	return t.AddDate(0, 0, delta).Format(DateLayout)
}

// Avec l'api il se peut que les données sont stocké avec un accès à la date
// exacte, le jours d'avant ou d'après, d'après la localisation de l'évènement.
func fetchRightDate(d string) (ESPNScoreboardResponse, error) {
	scoreBoard, err := fetchByDate(GetDayDelta(d, 0))
	if len(scoreBoard.Events) != 0 {
		return scoreBoard, err
	}
	// on teste avec le jour d'avant
	scoreBoard, err = fetchByDate(GetDayDelta(d, -1))
	if len(scoreBoard.Events) != 0 {
		return scoreBoard, err
	}
	// on teste avec le jour d'avant
	scoreBoard, err = fetchByDate(GetDayDelta(d, 1))
	if len(scoreBoard.Events) != 0 {
		return scoreBoard, err
	}

	// on doit pas arriver ici
	return scoreBoard, err
}

 
func getRecord(rs []ESPNRecord) string {
	if len(rs) == 0 {
		return ""
	}

	return rs[0].Summary
}

// L'id du gagnant retourné est l'id local à la base.
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

// L'id du gagnant retourné est l'id externe de l'API.
// Précondition: le combat est terminé.
func getWinnerByCompetID(competitions []ESPNCompetition, competID string) (string, bool) {
	for _, compet := range competitions{
		if compet.ID == competID{
			if compet.Competitors[0].Winner {
				return compet.Competitors[0].ID, true
			}else{
				return compet.Competitors[1].ID, true
			}
		}
	}
	return "", false
}


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

// L'api d'après tests peut laisser des combats annulés pendant certains jours en tant que "STATUS_SCHEDULED"
// et l'API ne donne pas explicitement un STATUS_CANCELED, donc cette fonction fait ce traitement
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
