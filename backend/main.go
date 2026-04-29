package main

//import "fantasy/server"
import (
	"fantasy/espn"
	"fmt"
	"log"
)

func main(){
	//database.InitDB()
	//server.Start()
	scoreboard, err := espn.Fetch()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(scoreboard.Season.Year)
	fmt.Println(scoreboard.Events[0].Name)
}