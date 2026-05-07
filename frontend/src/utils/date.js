export function parseDate(value) {
  if (!value) {
    return null;
  }

  const normalized = typeof value === "string" && value.includes(" ") && !value.includes("T")
    ? value.replace(" ", "T")
    : value;

  const date = new Date(normalized);
  return Number.isNaN(date.getTime()) ? null : date;
}

export function formatDate(value, options = {}) {
  const date = parseDate(value);

  if (!date) {
    return "Date a confirmer";
  }

  return new Intl.DateTimeFormat("fr-FR", {
    dateStyle: "medium",
    timeStyle: "short",
    ...options,
  }).format(date);
}

export function formatShortDate(value) {
  const date = parseDate(value);

  if (!date) {
    return "A confirmer";
  }

  return new Intl.DateTimeFormat("fr-FR", {
    day: "2-digit",
    month: "short",
    year: "numeric",
  }).format(date);
}

export function sortCardsByDate(cards) {
  return [...cards].sort((cardA, cardB) => {
    const dateA = parseDate(cardA.date)?.getTime() ?? Number.MAX_SAFE_INTEGER;
    const dateB = parseDate(cardB.date)?.getTime() ?? Number.MAX_SAFE_INTEGER;
    return dateA - dateB;
  });
}

export function isUpcomingCard(card) {
  if (card.completed) {
    return false;
  }

  const cardDate = parseDate(card.date);
  return !cardDate || cardDate >= new Date();
}
