package handlers

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fantasy/database"
	"net/http"
	"time"
)

type ErrorResp struct{
	Error string	`json:"error"`
}

func writeJsonResponse(w http.ResponseWriter, status int, data any){
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeJsonError(w http.ResponseWriter, status int, message string){
	writeJsonResponse(w, status, ErrorResp{Error: message})
}

/* User */

// Génère un token aléatoirement.
func generateSessionToken() (string) {
	b := make([]byte, 32)
	rand.Read(b)
	// Très important pour avoir des chaines correct à utiliser pour les url, cookies etc
	return base64.RawURLEncoding.EncodeToString(b)
}

func hashSessionToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

// Crée une session pour l'utilisateur avec l'id désigné, et insère dans la session dans la base.
// Retourne le token de la session
func createSession(userID int) (string, error) {

	token := generateSessionToken()
	tokenHash := hashSessionToken(token)
	expiresAt := time.Now().Add(7 * 24 * time.Hour)

	_, err := database.InsertSession(database.DB, userID, tokenHash, expiresAt)

	if err != nil {
		return "", err
	}

	return token, nil
}

func setSessionCookie(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    token,
		Path:     "/", // Le cookie sera envoyé pour toutes les routes proposées
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
		Secure:   true,   // HTTPS
		MaxAge:   7 * 24 * 60 * 60,
	})
}

func findUserIDBySessionToken(token string) (int, error) {
	tokenHash := hashSessionToken(token)
	var userID int

	err := database.GetValidSessionByTokenH(database.DB, tokenHash).Scan(&userID)

	if err != nil {
		return -1, err
	}

	return userID, nil
}

// Envoi un cookie qui va faire supprimer le cookie courant coté client.
func clearSessionCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     	"session_token",
		Value:    	"",
		Path: 		"/",
		HttpOnly: true,
		Secure:   true,
    	SameSite: http.SameSiteNoneMode,
		MaxAge:   -1,	// suppresion immédiate par le client
	})
}

func getAuthenticatedUserID(w http.ResponseWriter, req *http.Request) (int, bool) {
	cookie, err := req.Cookie("session_token")
	if err != nil {
		writeJsonError(w, http.StatusUnauthorized, "Non authentifié")
		return -1, false
	}

	userID, err := findUserIDBySessionToken(cookie.Value)
	if err != nil {
		writeJsonError(w, http.StatusUnauthorized, "Session invalide ou expirée")
		return -1, false
	}

	return userID, true
}