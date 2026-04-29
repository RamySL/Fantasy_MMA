package main

//import "fantasy/server"
import (
	"fantasy/database"
	"fantasy/espn"
	"log"
)

func main(){
	database.InitDB()
	//server.Start()
	err := espn.Sync()
	if (err != nil){
		log.Printf("Erreur lors de sync : %s" , err)
	}
}