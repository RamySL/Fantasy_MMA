import { useEffect, useMemo, useState } from "react";
import Badge from "../components/Badge.jsx";
import EmptyState from "../components/EmptyState.jsx";
import ErrorMessage from "../components/ErrorMessage.jsx";
import LoadingState from "../components/LoadingState.jsx";
import StatCard from "../components/StatCard.jsx";
import { predictionService } from "../services/predictionService.js";
import { formatShortDate } from "../utils/date.js";

function statusBadge(status) {
  if (status === "completed") {
    return <Badge variant="success">Terminee</Badge>;
  }

  return <Badge variant="info">En attente</Badge>;
}

function fightBadge(status) {
  if (status === "correct") {
    return <Badge variant="success">Correct</Badge>;
  }

  if (status === "wrong") {
    return <Badge variant="danger">Raté</Badge>;
  }

  return <Badge variant="info">En attente</Badge>;
}

export default function MyPredictionsPage() {
  const [predictionCards, setPredictionCards] = useState([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState("");

  useEffect(() => {
    let isMounted = true;

    async function loadPredictions() {
      try {
        const data = await predictionService.getMyPredictions();
        if (isMounted) {
          setPredictionCards(data || []);
        }
      } catch (err) {
        if (isMounted) {
          setError(err.message || "Impossible de charger tes predictions.");
        }
      } finally {
        if (isMounted) {
          setIsLoading(false);
        }
      }
    }

    loadPredictions();

    return () => {
      isMounted = false;
    };
  }, []);

  const totals = useMemo(() => {
    return predictionCards.reduce(
      (acc, card) => {
        acc.points += card.total_points || 0;
        acc.possiblePoints += card.possible_points || 0;
        acc.good += card.good_predictions || 0;
        acc.total += card.total_predictions || 0;
        return acc;
      },
      { points: 0, possiblePoints: 0, good: 0, total: 0 },
    );
  }, [predictionCards]);

  return (
    <div className="space-y-6">
      <div>
        <p className="text-sm font-semibold uppercase tracking-[0.25em] text-slate-500">Mon éspace</p>
        <h1 className="mt-2 text-3xl font-black tracking-tight text-slate-950">Mes prédictions</h1>
        <p className="mt-2 max-w-2xl text-sm leading-6 text-slate-600">
          Récapitulatif personnel de tes choix, points obtenus et prédictions encore en attente.
        </p>
      </div>

      {error && <ErrorMessage message={error} />}

      <section className="grid gap-4 sm:grid-cols-3">
        <StatCard label="Points" value={totals.points} helper={`${totals.possiblePoints} pts possibles`} />
        <StatCard label="Bons pronos" value={`${totals.good}/${totals.total}`} helper="Toutes cartes confondues" />
        <StatCard label="Cartes jouées" value={predictionCards.length} helper="Depuis tes prédictions sauvegardées" />
      </section>

      {isLoading ? (
        <LoadingState label="Chargement des prédictions..." />
      ) : predictionCards.length > 0 ? (
        <div className="space-y-5">
          {predictionCards.map((card) => (
            <article key={card.id} className="card-panel">
              <div className="flex flex-col gap-3 sm:flex-row sm:items-start sm:justify-between">
                <div>
                  <div className="flex flex-wrap items-center gap-2">
                    {statusBadge(card.status)}
                    <Badge>{formatShortDate(card.card_date)}</Badge>
                  </div>
                  <h2 className="mt-4 text-xl font-bold text-slate-950">{card.card_title}</h2>
                  <p className="mt-2 text-sm text-slate-500">
                    {card.total_points} pts obtenus sur {card.possible_points} possibles
                  </p>
                </div>
              </div>

              <div className="mt-5 overflow-hidden rounded-2xl border border-slate-200">
                <div className="hidden grid-cols-[1.1fr_1fr_1fr_auto] gap-4 bg-slate-100 px-4 py-3 text-xs font-bold uppercase tracking-wide text-slate-500 md:grid">
                  <span>Combat</span>
                  <span>Prédiction</span>
                  <span>Résultat</span>
                  <span>Points</span>
                </div>
                {card.fights.map((fight) => (
                  <div key={fight.fight_id} className="grid gap-3 border-t border-slate-200 px-4 py-4 text-sm md:grid-cols-[1.1fr_1fr_1fr_auto] md:items-center">
                    <div>
                      <p className="font-semibold text-slate-950">{fight.fighter1} vs {fight.fighter2}</p>
                      <p className="text-slate-500">{fight.category}</p>
                    </div>
                    <p className="text-slate-700">{fight.predicted_winner}</p>
                    <div className="flex flex-wrap items-center gap-2">
                      {fightBadge(fight.status)}
                      <span className="text-slate-500">{fight.official_winner || "A venir"}</span>
                    </div>
                    <p className="font-bold text-slate-950">{fight.points_obtained ?? 0}/{fight.points_available}</p>
                  </div>
                ))}
              </div>
            </article>
          ))}
        </div>
      ) : (
        <EmptyState title="Aucune prédiction" message="Tes prédictions sauvegardées apparaitront ici." actionLabel="Voir les cartes" actionTo="/cards" />
      )}
    </div>
  );
}
