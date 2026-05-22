package handlers

import (
	"fantasy/database"
	"log"
	"net/http"
)

// Retourne le classement global des utilisateurs.
// Le score est basé sur le totale des points puis sur le nombre de bonne prédiction.
func Ranking(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		writeJsonError(w, http.StatusMethodNotAllowed, "Méthode non autorisée")
		return
	}

	rows, err := database.GetGlobalRanking(database.DB)

	if err != nil {
		log.Printf("ranking.Ranking: erreur query: %v", err)
		writeJsonError(w, http.StatusInternalServerError, "Erreur serveur interne")
		return
	}
	defer rows.Close()

	entries := []RankingEntryResponse{}
	rank := 1

	for rows.Next() {
		var entry RankingEntryResponse

		entry.Rank = rank

		err := rows.Scan(
			&entry.Pseudo,
			&entry.TotalPoints,
			&entry.GoodPredictions,
			&entry.TotalPredictions,
		)

		if err != nil {
			log.Printf("ranking.Ranking: erreur scan: %v", err)
			writeJsonError(w, http.StatusInternalServerError, "Erreur serveur interne")
			return
		}

		entries = append(entries, entry)
		rank++
	}

	writeJsonResponse(w, http.StatusOK, entries)
}