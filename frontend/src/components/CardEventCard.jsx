import { Link } from "react-router-dom";
import Badge from "./Badge.jsx";
import { formatDate } from "../utils/date.js";

function getLocation(card) {
  return [card.venue_name, card.city, card.region, card.country].filter(Boolean).join(", ");
}

export default function CardEventCard({ card }) {
  const location = getLocation(card);
  const isCompleted = Boolean(card.completed);

  return (
    <article className="card-panel flex h-full flex-col justify-between gap-5 transition hover:-translate-y-0.5 hover:shadow-lg">
      <div>
        <div className="flex flex-wrap items-center gap-2">
          <Badge variant={isCompleted ? "success" : "warning"}>{isCompleted ? "Terminée" : card.status || "A venir"}</Badge>
          <span className="text-xs font-medium text-slate-500">#{card.id}</span>
        </div>

        <h3 className="mt-4 text-xl font-bold tracking-tight text-slate-950">{card.title}</h3>
        <p className="mt-3 text-sm text-slate-600">{formatDate(card.date)}</p>
        <p className="mt-2 text-sm text-slate-500">{location || "Lieu a confirmer"}</p>
      </div>

      <div className="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
        <p className="text-xs text-slate-500">Predictions {isCompleted ? "fermees" : "ouvertes selon statut"}</p>
        <Link to={`/cards/${card.id}`} className="btn-secondary">
          Voir la carte
        </Link>
      </div>
    </article>
  );
}
