import Badge from "./Badge.jsx";
import { getFightWinnerName, getWinnerId } from "../utils/fights.js";

function FighterButton({ fighter, isSelected, isWinner, disabled, onClick }) {
  const {
    full_name = "Combattant inconnu",
    record = "Record inconnu",
    fighter_image_url = null,
    fighter_country_name = null,
    fighter_country_flag_url = null,
  } = fighter;

  return (
    <button
      type="button"
      disabled={disabled}
      aria-pressed={isSelected}
      onClick={onClick}
      className={`rounded-2xl border p-4 text-left transition ${
        isSelected
          ? "border-slate-950 bg-slate-950 text-white shadow-md"
          : "border-slate-200 bg-white text-slate-900 hover:border-slate-400 hover:bg-slate-50"
      } ${disabled ? "cursor-not-allowed opacity-70" : ""}`}
    >
      <div className="flex items-start gap-3">
        {/* Image du combattant */}
        {fighter_image_url ? (
          <img
            src={fighter_image_url}
            alt={full_name}
            className="h-20 w-20 rounded-full object-cover border border-slate-200"
          />
        ) : (
          <div className="h-14 w-14 rounded-full bg-slate-200 flex items-center justify-center text-slate-500 text-xs">
            Pas d'image
          </div>
        )}

        <div className="flex-1">
          <span className="block text-base font-bold">{full_name}</span>

          {/* Pays + drapeau */}
          <div className="flex items-center gap-1.5 mt-1">
            {fighter_country_flag_url && (
              <img
                src={fighter_country_flag_url}
                alt={fighter_country_name || "Drapeau"}
                className="h-4 w-6 object-cover rounded-sm"
              />
            )}
            <span
              className={`text-xs ${isSelected ? "text-slate-200" : "text-slate-500"}`}
            >
              {fighter_country_name || "Pays inconnu"}
            </span>
          </div>

          {/* Record */}
          <span
            className={`block text-sm mt-1 ${isSelected ? "text-slate-200" : "text-slate-500"}`}
          >
            {record}
          </span>
        </div>
      </div>

      {/* Badge gagnant officiel */}
      {isWinner && (
        <span
          className={`mt-3 inline-flex rounded-full px-2.5 py-1 text-xs font-semibold ${
            isSelected
              ? "bg-white/15 text-white"
              : "bg-emerald-50 text-emerald-700"
          }`}
        >
          Gagnant
        </span>
      )}
    </button>
  );
}

export default function FightPredictionCard({
  fight,
  selectedWinnerId,
  onSelect,
  disabled = false,
}) {
  const winnerId = getWinnerId(fight.winner_id);
  const officialWinnerName = getFightWinnerName(fight);
  const isCompleted = Boolean(fight.completed);
  const isCanceled = fight.status == "STATUS_CANCELED"
  const selectionDisabled = disabled || isCanceled || isCompleted;

  return (
    <article className="card-panel">
      <div className="flex flex-col gap-4 sm:flex-row sm:items-start sm:justify-between">
        <div>
          <div className="flex flex-wrap items-center gap-2">
            <Badge variant={isCanceled ? "warning" : isCompleted ? "success" : "info"}>
              {
                isCanceled ? "Combat annulé" : 
                  isCompleted ? "Combat terminé" :  "À venir"
              }
            </Badge>
            {fight.category && <Badge>{fight.category}</Badge>}
          </div>
          <h3 className="mt-4 text-lg font-bold text-slate-950">
            {fight.fighter1.full_name}{" "}
            <span className="text-slate-400">vs</span>{" "}
            {fight.fighter2.full_name}
          </h3>
        </div>

        <div className="rounded-2xl bg-slate-100 px-4 py-3 text-sm font-semibold text-slate-700">
          {fight.points_good_prediction || 0} pts
        </div>
      </div>

      <div className="mt-5 grid gap-3 md:grid-cols-[1fr_auto_1fr] md:items-stretch">
        <FighterButton
          fighter={fight.fighter1}
          isSelected={selectedWinnerId === fight.fighter1.id}
          isWinner={winnerId === fight.fighter1.id}
          disabled={selectionDisabled}
          onClick={() => onSelect(fight.id, fight.fighter1.id)}
        />

        <div className="flex items-center justify-center text-xs font-black uppercase tracking-[0.35em] text-slate-400">
          vs
        </div>

        <FighterButton
          fighter={fight.fighter2}
          isSelected={selectedWinnerId === fight.fighter2.id}
          isWinner={winnerId === fight.fighter2.id}
          disabled={selectionDisabled}
          onClick={() => onSelect(fight.id, fight.fighter2.id)}
        />
      </div>

      <div className="mt-4 flex flex-col gap-2 text-sm text-slate-500 sm:flex-row sm:items-center sm:justify-between">
        <span>
          {selectionDisabled
            ? ""
            : "Choisis le gagnant avant le début du combat."}
        </span>
      </div>
    </article>
  );
}