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
    const predictions = Object.entries(selectedPredictions)
      .filter(([, fighterId]) => Boolean(fighterId))
      .map(([fightId, fighterId]) => ({
        fight_id: Number(fightId),
        predicted_winner_id: Number(fighterId),
      }));

    if (predictions.length === 0) {
      return { card_id: Number(cardId), saved_count: 0 };
    }

    const result = await apiRequest("/predictions/bulk", {
      method: "POST",
      body: { predictions },
    });

    return {
      card_id: Number(cardId),
      saved_count: result.saved_count,
    };
  },
};