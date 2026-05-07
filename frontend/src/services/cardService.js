import { apiRequest } from "./api.js";

export const cardService = {
  getCards() {
    return apiRequest("/cards");
  },

  getCard(cardId) {
    return apiRequest(`/cards/${cardId}`);
  },

  getCardFights(cardId) {
    return apiRequest(`/cards/${cardId}/fights`);
  },
};
