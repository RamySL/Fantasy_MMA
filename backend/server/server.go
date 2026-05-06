package server

/*
 * ce Squellette est pris de : https://gobyexample.com/http-servers 
*/

import (
    //"fmt"
    "net/http"
    "fantasy/handlers"
)


func Start() {
    // Handlers 
    http.HandleFunc("/cards", handlers.GetCards)
    http.HandleFunc("/cards/", handlers.GetCardsRoutes)
    http.HandleFunc("/auth/", handlers.PostAuthRoutes) 
    http.HandleFunc("/auth/me", handlers.Me)
    http.ListenAndServe(":8090", nil)
}
