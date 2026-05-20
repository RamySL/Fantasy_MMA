import { useEffect, useMemo, useState } from "react";
import { Link } from "react-router-dom";
import CardEventCard from "../components/CardEventCard.jsx";
import EmptyState from "../components/EmptyState.jsx";
import ErrorMessage from "../components/ErrorMessage.jsx";
import LoadingState from "../components/LoadingState.jsx";
import StatCard from "../components/StatCard.jsx";
import { useAuth } from "../context/AuthContext.jsx";
import { cardService } from "../services/cardService.js";
import { isUpcomingCard, sortCardsByDate } from "../utils/date.js";

export default function HomePage() {
  const { user, isAuthenticated } = useAuth();
  const [cards, setCards] = useState([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState("");

  useEffect(() => {
    let isMounted = true;

    async function loadCards() {
      try {
        const data = await cardService.getCards();
        if (isMounted) {
          setCards(data || []);
        }
      } catch (err) {
        if (isMounted) {
          setError(err.message || "Impossible de charger les cartes.");
        }
      } finally {
        if (isMounted) {
          setIsLoading(false);
        }
      }
    }

    loadCards();

    return () => {
      isMounted = false;
    };
  }, []);

  const dashboard = useMemo(() => {
    const sorted = sortCardsByDate(cards);
    const upcoming = sorted.filter(isUpcomingCard);
    const completed = sorted.filter((card) => card.completed);

    return {
      nextCard: upcoming[0],
      upcomingCount: upcoming.length,
      completedCount: completed.length,
      totalCount: sorted.length,
    };
  }, [cards]);

  return (
    <div className="space-y-8">
      <section className="grid gap-6 lg:grid-cols-[1.3fr_0.7fr] lg:items-stretch">
        <div className="card-panel flex flex-col justify-between bg-slate-950 text-white">
          <div>
            <p className="text-sm font-semibold uppercase tracking-[0.3em] text-slate-300">Fantasy MMA</p>
            <h1 className="mt-4 max-w-3xl text-4xl font-black tracking-tight sm:text-5xl">
              Fais tes predictions, suis tes points, grimpe au classement.
            </h1>
            <p className="mt-5 max-w-2xl text-base leading-7 text-slate-300">
              Une interface simple pour choisir les gagnants de chaque combat et comparer tes resultats avec les autres participants.
            </p>
          </div>

          <div className="mt-8 flex flex-col gap-3 sm:flex-row">
            <Link to="/cards" className="btn-primary bg-white text-slate-950 hover:bg-slate-100">
              Voir les cartes
            </Link>
            <Link to="/leaderboard" className="btn-secondary border-white/20 bg-white/10 text-white hover:bg-white/15">
              Voir le classement
            </Link>
          </div>
        </div>

        <div className="card-panel">
          <p className="text-sm font-semibold text-slate-500">Bienvenue</p>
          <h2 className="mt-3 text-2xl font-bold text-slate-950">
            {isAuthenticated ? `Salut ${user.pseudo}` : "Connecte-toi pour sauvegarder tes choix"}
          </h2>
          {!isAuthenticated && (
            <Link to="/login" className="btn-primary mt-6 w-full">
              Se connecter
            </Link>
          )}
        </div>
      </section>

      <section className="grid gap-4 sm:grid-cols-3">
        <StatCard label="Cartes totales" value={dashboard.totalCount} />
        <StatCard label="Cartes à venir" value={dashboard.upcomingCount} helper="Predictions ouvertes selon statut" />
        <StatCard label="Cartes terminées" value={dashboard.completedCount} helper="Resultats disponibles" />
      </section>

      {error && <ErrorMessage message={error} />}

      <section className="space-y-4">
        <div className="flex flex-col gap-2 sm:flex-row sm:items-center sm:justify-between">
          <div>
            <h2 className="section-title">Prochaine carte</h2>
            <p className="mt-1 text-sm text-slate-500">Acces rapide a la carte la plus proche.</p>
          </div>
          <Link to="/cards" className="btn-secondary">
            Toutes les cartes
          </Link>
        </div>

        {isLoading ? (
          <LoadingState label="Chargement des cartes..." />
        ) : dashboard.nextCard ? (
          <CardEventCard card={dashboard.nextCard} />
        ) : (
          <EmptyState title="Aucune carte a venir" message="Les prochaines cartes apparaitront ici quand le backend les retournera." actionLabel="Voir toutes les cartes" actionTo="/cards" />
        )}
      </section>
    </div>
  );
}
