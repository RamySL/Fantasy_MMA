import { mockMyPredictions } from "../data/mockData.js";

function sleep(ms) {
  return new Promise((resolve) => window.setTimeout(resolve, ms));
}

export const predictionService = {
  async getMyPredictions() {
    // TODO backend: remplacer par apiRequest("/predictions/me").
    await sleep(250);
    return mockMyPredictions;
  },

  async getPredictionsForCard(cardId) {
    // TODO backend: remplacer par apiRequest(`/cards/${cardId}/predictions/me`).
    await sleep(150);
    console.info("TODO backend: load predictions for card", cardId);
    return [];
  },

  async saveCardPredictions(cardId, selectedPredictions) {
    // On transforme l'objet local { fightId: fighterId } en tableau propre
    // pour coller a un futur payload backend simple a traiter en Go.
    const predictions = Object.entries(selectedPredictions).map(([fightId, fighterId]) => ({
      fight_id: Number(fightId),
      predicted_winner_id: Number(fighterId),
    }));

    // TODO backend: remplacer ce mock par un vrai appel, par exemple :
    // return apiRequest(`/cards/${cardId}/predictions`, {
    //   method: "POST",
    //   body: { predictions },
    // });
    await sleep(400);
    console.info("TODO backend: save predictions", {
      card_id: Number(cardId),
      predictions,
    });

    return {
      card_id: Number(cardId),
      predictions,
      mocked: true,
    };
  },
};
