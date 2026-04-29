package espn

import (
	"database/sql"
	"fantasy/database"
)

/* Utilisé pour synchroniser régulièrement la base de données. 

- Chaque carte est réalisée la nuit du samedi à dimanche. Donc par exemple Dimanche 
12h serait un bon temps pour actuliser les résultat d'une carte.

- Par exemple les annonces de cartes son faite millieu de semaine. Faire un fetch le jeudi.

*/

// Fetch les cartes à travers l'api et insert dans la base
func SyncCards() {
	tx, err := database.DB.Begin()

	scoreBoard, err := Fetch()
	//CAREFUL: Leagues est supposé être singleton
	
}

func syncCard(tx *sql.Tx, scoreBoard ESPNScoreboardResponse, i int) {

	cardID := extractID(scoreBoard.Leagues[0].Calendar[i].Event.Ref)
	database.UpsertCard(tx, cardID, scoreBoard. )
}

// Fetch les combats à travers l'api et insert dans la base
func syncFights() {

}

//TODO:helper
func extractID(ref string) (string){
	return ""
}