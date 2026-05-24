package main

import (
	"fantasy/database"
	"fantasy/espn"
	"time"
	"fantasy/server"
)

func main(){

	database.InitDB()

	// Remplissage initale de la base
	go func(){
		espn.Sync(true)
	}()
	
	//Actualisation de la base de données avec les résultats 
	go func () {
		espn.SyncTicker(12, time.Hour * 1)
	}()

	// Suppression des sessions expirées
	go func() {
		ticker := time.NewTicker(1 * time.Hour)
		defer ticker.Stop()
		for range ticker.C {
			_, _ = database.DeleteExpiredSessions(database.DB)
		}
	}()

	server.Start()
}