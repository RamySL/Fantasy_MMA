package handlers

import (
	"encoding/json"
	"fantasy/database"
	"log"
	"net/http"
	"strings"

	"github.com/lib/pq"
	"github.com/lib/pq/pqerror"
	"golang.org/x/crypto/bcrypt"
	"database/sql"
)

func PostAuthRoutes(w http.ResponseWriter, req *http.Request){
	if req.Method != http.MethodPost {
		 writeJsonError(w, http.StatusMethodNotAllowed, "Méthode non autorisée")
		 return
	}

	path := strings.Trim(req.URL.Path, "/")
	parts := strings.Split(path, "/")

	
	if len(parts) == 2 {
		// auth/register
		if parts[1] == "register"{
			registerUser(w, req)
			return
		}
		// auth/login
		if parts[1] == "login"{
			logInUser(w, req)
			return
		}
		// auth/logout
		if parts[1] == "logout"{
			logOutUser(w, req)
			return
		}

	}

	http.NotFound(w, req)

}

func registerUser(w http.ResponseWriter, req *http.Request){

	var body RegisterRequestBody
	err := json.NewDecoder(req.Body).Decode(&body)
	if err != nil {
		writeJsonError(w, http.StatusBadRequest, "JSON invalide")
		return
	}

	pwdHash, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		writeJsonError(w, http.StatusInternalServerError, "Erreur serveur interne")
		return
	}
	var userResp AuthUserResponse 
	// insertion dans la base
	err = database.DB.QueryRow(`
		INSERT INTO users (pseudo, email, password_hash)
		VALUES ($1, $2, $3)
		RETURNING id, pseudo, email
	`,
	body.Pseudo,
	body.Email,
	string(pwdHash),
	).Scan(&userResp.ID, &userResp.Pseudo, &userResp.Email)

	if err != nil {
		pqError := pq.As(err, pqerror.UniqueViolation)

		if pqError != nil {
			
			if pqError.Constraint == "users_pseudo_key" {
				writeJsonError(w, http.StatusConflict, "Pseudo déjà pris")
				return
			}

			if pqError.Constraint == "users_email_key" {
				writeJsonError(w, http.StatusConflict, "Email déjà utilisé")
				return
			}
		}
		log.Printf("Erreur Scan : %v", err)
		writeJsonError(w, http.StatusInternalServerError, "Erreur création utilisateur")
		return
	}

	token, err := createSession(userResp.ID)
	if(err != nil){
		log.Printf("Erreur création session %v", err)
		writeJsonError(w, http.StatusInternalServerError, "Erreur création session")
		return
	}

	setSessionCookie(w, token)
	writeJsonResponse(w, http.StatusCreated, userResp)
}

func logInUser(w http.ResponseWriter, req *http.Request) {
	var body LoginRequest

	err := json.NewDecoder(req.Body).Decode(&body)
	if err != nil {
		writeJsonError(w, http.StatusBadRequest, "JSON invalide")
		return
	}

	var user AuthUserResponse
	var passwordHash string

	err = database.DB.QueryRow(`
		SELECT id, pseudo, email, password_hash
		FROM users
		WHERE email = $1
	`, body.Email).Scan(
		&user.ID,
		&user.Pseudo,
		&user.Email,
		&passwordHash,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			writeJsonError(w, http.StatusUnauthorized, "Identifiants invalides")
			return
		}

		writeJsonError(w, http.StatusInternalServerError, "Erreur lecture utilisateur")
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(body.Password))
	if err != nil {
		writeJsonError(w, http.StatusUnauthorized, "mote de passe invalide")
		return
	}

	token, err := createSession(user.ID)
	if err != nil {
		log.Printf("Erreur création session %v", err)
		writeJsonError(w, http.StatusInternalServerError, "Erreur création session")
		return
	}

	setSessionCookie(w, token)
	writeJsonResponse(w, http.StatusOK, user)
}

func logOutUser(w http.ResponseWriter, req *http.Request) {
	cookie, err := req.Cookie("session_token")

	if err == nil {
		err = deleteSession(cookie.Value)
		if err != nil {
			writeJsonError(w, http.StatusInternalServerError, "Erreur suppression session")
			return
		}
	}

	clearSessionCookie(w)
	writeJsonResponse(w, http.StatusOK, map[string]string{
		"message": "Déconnexion réussie",
	})
}

func Me(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		writeJsonError(w, http.StatusMethodNotAllowed, "Méthode non autorisée")
		return
	}

	cookie, err := req.Cookie("session_token")
	if err != nil {
		writeJsonError(w, http.StatusUnauthorized, "Non authentifié")
		return
	}

	userID, err := findUserIDBySessionToken(cookie.Value)
	if err != nil {
		writeJsonError(w, http.StatusUnauthorized, "Session invalide ou expirée")
		return
	}

	var user AuthUserResponse

	err = database.DB.QueryRow(`
		SELECT id, pseudo, email
		FROM users
		WHERE id = $1
	`, userID).Scan(
		&user.ID,
		&user.Pseudo,
		&user.Email,
	)

	if err != nil {
		writeJsonError(w, http.StatusInternalServerError, "Erreur lecture utilisateur")
		return
	}

	writeJsonResponse(w, http.StatusOK, user)
}