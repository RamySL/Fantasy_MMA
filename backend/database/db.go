package database

/*
Inspiré de : https://github.com/lib/pq
*/

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() {
	connStr := "host=localhost port=5432 user=postgres password=postgres dbname=fantasy_mma connect_timeout=5 sslmode=disable"

	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Erreur ouverture DB:", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatal("Erreur connexion DB:", err)
	}

	log.Println("Connexion PostgreSQL réussie")
}