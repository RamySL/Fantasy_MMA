package main

import (
	"fantasy/database"
	"fantasy/espn"
	//"time"
	"fantasy/server"
)

func main(){

	database.InitDB()

	go func ()  {
		//espn.SyncTicker(21, time.Hour * 1)
		espn.Sync()
	}()
	

	server.Start()
}