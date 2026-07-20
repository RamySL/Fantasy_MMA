package main

import (
	"fantasy/database"
	"fantasy/espn_api"
	"fantasy/server"
	"log"
	"time"
)

func main(){

	database.InitDB()

	// Remplissage initale de la base
	go func(){
		err := espn.InitialSync()
		if err != nil{
			log.Printf("[main] : erreur synchronisation init : %v ", err)
		}
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