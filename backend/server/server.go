package server

/*
 * ce Squellette est pris de : https://gobyexample.com/http-servers
*	Cors prblem : https://www.stackhawk.com/blog/golang-cors-guide-what-it-is-and-how-to-enable-it/
*/

import (
	"fantasy/handlers"
	"log"
	"net/http"
	"os"
)

// corsMiddleware autorise le service d'hébergement du client (Vercel) à discuter avec le serveur (Render)
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// ex: https://****.vercel.app)
		origin := r.Header.Get("Origin")

		if origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		} else {
			w.Header().Set("Access-Control-Allow-Origin", "*")
		}

		// pour : credentials: "include"
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func Start() {

	mux := http.NewServeMux()

	// Cards
	mux.HandleFunc("/cards", handlers.GetCards)
	mux.HandleFunc("/cards/", handlers.GetCardsRoutes)

	// Predictions
	mux.HandleFunc("/predictions", handlers.Predictions)
	mux.HandleFunc("/predictions/bulk", handlers.PredictionsBulk)
	mux.HandleFunc("/predictions/", handlers.PredictionsRoutes)

	// Auth
	mux.HandleFunc("/auth/me", handlers.Me)
	mux.HandleFunc("/auth/", handlers.PostAuthRoutes)

	// Ranking
	mux.HandleFunc("/ranking", handlers.Ranking)

	handlerCORS := corsMiddleware(mux)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8090"
	}

	log.Printf("Serveur démarré sur le port %s", port)

	err := http.ListenAndServe(":"+port, handlerCORS)
	if err != nil {
		log.Fatal("Erreur serveur : ", err)
	}
}