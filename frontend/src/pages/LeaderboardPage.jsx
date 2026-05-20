import { useEffect, useState } from "react";
import Badge from "../components/Badge.jsx";
import EmptyState from "../components/EmptyState.jsx";
import ErrorMessage from "../components/ErrorMessage.jsx";
import LoadingState from "../components/LoadingState.jsx";
import { leaderboardService } from "../services/leaderboardService.js";

function accuracy(entry) {
  if (!entry.total_predictions) {
    return "0%";
  }

  return `${Math.round((entry.good_predictions / entry.total_predictions) * 100)}%`;
}

function rankVariant(rank) {
  if (rank === 1) return "success";
  if (rank === 2 || rank === 3) return "warning";
  return "default";
}

export default function LeaderboardPage() {
  const [entries, setEntries] = useState([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState("");

  useEffect(() => {
    let isMounted = true;

    async function loadLeaderboard() {
      try {
        const data = await leaderboardService.getLeaderboard();
        if (isMounted) {
          setEntries(data || []);
        }
      } catch (err) {
        if (isMounted) {
          setError(err.message || "Impossible de charger le classement.");
        }
      } finally {
        if (isMounted) {
          setIsLoading(false);
        }
      }
    }

    loadLeaderboard();

    return () => {
      isMounted = false;
    };
  }, []);

  return (
    <div className="space-y-6">
      <div>
        <p className="mt-2 text-3xl font-black tracking-tight text-slate-950">Classement</p>
        <p className="mt-2 max-w-2xl text-sm leading-6 text-slate-600">
          Compare les points des participants et suis les meilleurs pronostiqueurs.
        </p>
      </div>

      {error && <ErrorMessage message={error} />}

      {isLoading ? (
        <LoadingState label="Chargement du classement..." />
      ) : entries.length > 0 ? (
        <div className="card-panel overflow-hidden p-0">
          <div className="hidden grid-cols-[90px_1fr_130px_150px_110px] gap-4 border-b border-slate-200 bg-slate-100 px-5 py-4 text-xs font-black uppercase tracking-wide text-slate-500 md:grid">
            <span>Rang</span>
            <span>Participant</span>
            <span>Points</span>
            <span>Bons pronos</span>
            <span>Reussite</span>
          </div>

          {entries.map((entry) => (
            <div key={entry.pseudo} className="grid gap-3 border-b border-slate-100 px-5 py-5 text-sm last:border-b-0 md:grid-cols-[90px_1fr_130px_150px_110px] md:items-center">
              <div>
                <Badge variant={rankVariant(entry.rank)}>#{entry.rank}</Badge>
              </div>
              <div>
                <p className="font-bold text-slate-950">{entry.pseudo}</p>
                <p className="text-slate-500 md:hidden">{entry.total_points} points</p>
              </div>
              <p className="font-black text-slate-950">{entry.total_points}</p>
              <p className="text-slate-700">{entry.good_predictions}/{entry.total_predictions}</p>
              <p className="font-semibold text-slate-950">{accuracy(entry)}</p>
            </div>
          ))}
        </div>
      ) : (
        <EmptyState title="Classement vide" message="Le classement apparaitra ici quand le backend retournera les scores." />
      )}
    </div>
  );
}
