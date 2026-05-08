import { apiRequest } from "./api.js";

export const leaderboardService = {
  getLeaderboard() {
    return apiRequest("/ranking");
  },
};