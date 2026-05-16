package server

/*
 * ce Squellette est pris de : https://gobyexample.com/http-servers
 */

import (
	"fantasy/handlers"
	"net/http"
)

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

	http.ListenAndServe(":8090", nil)
}
