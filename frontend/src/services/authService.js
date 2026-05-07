import { apiRequest } from "./api.js";

export const authService = {
  register(payload) {
    return apiRequest("/auth/register", {
      method: "POST",
      body: payload,
    });
  },

  login(payload) {
    return apiRequest("/auth/login", {
      method: "POST",
      body: payload,
    });
  },

  logout() {
    return apiRequest("/auth/logout", {
      method: "POST",
    });
  },

  getMe() {
    return apiRequest("/auth/me");
  },
};
