package handlers

import (
	"fantasy/database"
	"log"
	"net/http"
)

// Retourne le classement global des utilisateurs.
// Le score est basé sur predictions.points_obtained
func Ranking(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		writeJsonError(w, http.StatusMethodNotAllowed, "Méthode non autorisée")
		return
	}
	rows, err := database.DB.Query(`
		SELECT
			users.pseudo,
			COALESCE(SUM(predictions.points_obtained), 0)::int AS total_points,
			COALESCE(SUM(
				CASE
					WHEN predictions.points_obtained > 0 THEN 1
					ELSE 0
				END
			), 0)::int AS good_predictions,
			COUNT(predictions.id)::int AS total_predictions
		FROM users
		LEFT JOIN predictions ON predictions.user_id = users.id
		GROUP BY users.id, users.pseudo
		ORDER BY total_points DESC, good_predictions DESC, total_predictions DESC, users.pseudo ASC
	`)

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