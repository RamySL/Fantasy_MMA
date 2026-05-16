package main

import (
	"fantasy/database"
	"fantasy/espn"
	"time"
	"fantasy/server"
)

func main(){

	database.InitDB()

	go func ()  {
		espn.SyncTicker(12, time.Hour * 24)
	}()

	server.Start()
}