const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || "/api";

export class ApiError extends Error {
  constructor(message, status, payload) {
    super(message);
    this.name = "ApiError";
    this.status = status;
    this.payload = payload;
  }
}

function buildUrl(path) {
  if (path.startsWith("http")) {
    return path;
  }

  return `${API_BASE_URL}${path.startsWith("/") ? path : `/${path}`}`;
}

async function parseJsonResponse(response) {
  const rawText = await response.text();

  if (!rawText) {
    return null;
  }

  try {
    return JSON.parse(rawText);
  } catch {
    return rawText;
  }
}

export async function apiRequest(path, options = {}) {
  const config = {
    credentials: "include",
    ...options,
    headers: {
      Accept: "application/json",
      ...(options.headers || {}),
    },
  };

  // Les handlers Go lisent du JSON. On serialise ici les objets JS pour eviter
  // de repeter la meme logique dans chaque service API.
  if (config.body && !(config.body instanceof FormData)) {
    config.headers = {
      "Content-Type": "application/json",
      ...config.headers,
    };

    if (typeof config.body !== "string") {
      config.body = JSON.stringify(config.body);
    }
  }

  const response = await fetch(buildUrl(path), config);
  const payload = await parseJsonResponse(response);

  if (!response.ok) {
    const message = payload?.error || payload?.message || "Erreur lors de la requete";
    throw new ApiError(message, response.status, payload);
  }

  return payload;
}
