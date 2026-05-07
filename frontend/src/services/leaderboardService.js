import { mockLeaderboard } from "../data/mockData.js";

function sleep(ms) {
  return new Promise((resolve) => window.setTimeout(resolve, ms));
}

export const leaderboardService = {
  async getLeaderboard() {
    // TODO backend: remplacer par apiRequest("/leaderboard").
    await sleep(250);
    return mockLeaderboard;
  },
};
