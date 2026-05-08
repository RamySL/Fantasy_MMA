package server

/*
 * ce Squellette est pris de : https://gobyexample.com/http-servers
 */

import (
	"fantasy/espn"
	"fantasy/handlers"
	"fmt"
	"net/http"
)

func syncPreds(w http.ResponseWriter, req *http.Request) {
	err := espn.SyncPredictions()
	if err != nil {
		fmt.Printf("Erreur sync preds %v", err)
	}
}

func Start() {
	// Cards
	http.HandleFunc("/cards", handlers.GetCards)
	http.HandleFunc("/cards/", handlers.GetCardsRoutes)

	// Predictions
	http.HandleFunc("/predictions", handlers.Predictions)
	http.HandleFunc("/predictions/", handlers.PredictionsRoutes)

	// Auth
	http.HandleFunc("/auth/me", handlers.Me)
	http.HandleFunc("/auth/", handlers.PostAuthRoutes)

	// Ranking
	http.HandleFunc("/ranking", handlers.Ranking)

	// tests
	http.HandleFunc("/syncPreds", syncPreds)

	http.ListenAndServe(":8090", nil)
}
