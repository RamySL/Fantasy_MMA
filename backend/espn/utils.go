package espn

import(
	"database/sql"
	"fantasy/database"
	"time"
	"log"
)

func upsertFighter (tx *sql.Tx, c ESPNCompetitor) (int, error) {
	return database.UpsertFighter(
		tx,
		c.ID,
		c.Athlete.FullName,
		getRecord(c),
	)
}

func getDayDelta(d string, delta int) (string){
	t, err := time.Parse("2006-01-02T15:04Z", d)
	if err != nil {
		log.Printf("Erreur de Parse dans getDayBefore : ", err)
	}

	return t.AddDate(0, 0, delta).Format("20060102")
}

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

 
func getRecord(c ESPNCompetitor) string {
	if len(c.Records) == 0 {
		return ""
	}

	return c.Records[0].Summary
}

 
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