import { apiRequest } from "./api.js";

export const predictionService = {
  async getMyPredictions() {
    return apiRequest("/predictions/me");
  },

  async getPredictionsForCard(cardId) {
    return apiRequest(`/cards/${cardId}/predictions/me`);
  },

  async saveCardPredictions(cardId, selectedPredictions) {
    // selectedPredictions a cette forme :
    // {
    //   [fightId]: fighterId
    // }
    //
    // Le backend actuel n'a pas de route bulk, donc on envoie une requete
    // POST /predictions par combat selectionne.
    const predictions = Object.entries(selectedPredictions)
      .filter(([, fighterId]) => Boolean(fighterId))
      .map(([fightId, fighterId]) => ({
        fight_id: Number(fightId),
        predicted_winner_id: Number(fighterId),
      }));
    //TODO: faire dans le backend une route qui prend une liste de prédictions.
    await Promise.all(
      predictions.map((prediction) =>
        apiRequest("/predictions", {
          method: "POST",
          body: prediction,
        }),
      ),
    );

    return {
      card_id: Number(cardId),
      saved_count: predictions.length,
    };
  },
};