package handlers

import ( 
	"net/http"
	"encoding/json"
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