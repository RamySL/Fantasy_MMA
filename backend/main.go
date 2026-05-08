package main

//import "fantasy/server"
import (
	"fantasy/database"
	//"fantasy/espn"
	"fantasy/server"
	//"log"
)

func main(){

	database.InitDB()

	/*err := espn.Sync()
	if (err != nil){
		log.Printf("Erreur lors de sync : %s" , err)
	}*/

	server.Start()

}