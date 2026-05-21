import { useEffect, useMemo, useState } from "react";
import CardEventCard from "../components/CardEventCard.jsx";
import EmptyState from "../components/EmptyState.jsx";
import ErrorMessage from "../components/ErrorMessage.jsx";
import LoadingState from "../components/LoadingState.jsx";
import { cardService } from "../services/cardService.js";
import { isUpcomingCard, sortCardsByDate } from "../utils/date.js";

const filters = [
  { id: "all", label: "Toutes" },
  { id: "upcoming", label: "A venir" },
  { id: "completed", label: "Terminées" },
];

export default function CardsPage() {
  const [cards, setCards] = useState([]);
  const [filter, setFilter] = useState("all");
  const [search, setSearch] = useState("");
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

  const visibleCards = useMemo(() => {
    const query = search.trim().toLowerCase();

    return sortCardsByDate(cards).filter((card) => {
      const matchesFilter = filter === "all" || (filter === "upcoming" && isUpcomingCard(card)) || (filter === "completed" && card.completed);
      const searchableText = [card.title, card.city, card.country, card.status].filter(Boolean).join(" ").toLowerCase();
      const matchesSearch = !query || searchableText.includes(query);

      return matchesFilter && matchesSearch;
    });
  }, [cards, filter, search]);

  return (
    <div className="space-y-6">
      <div className="flex flex-col gap-4 lg:flex-row lg:items-end lg:justify-between">
        <div>
          <p className="text-sm font-semibold uppercase tracking-[0.25em] text-slate-500">Cartes</p>
          <h1 className="mt-2 text-3xl font-black tracking-tight text-slate-950">Toutes les cartes de combats</h1>
          <p className="mt-2 max-w-2xl text-sm leading-6 text-slate-600">
            Selectionne une carte pour voir les combats et faire tes predictions.
          </p>
        </div>

        <div className="flex flex-col gap-3 sm:flex-row sm:items-center">
          <input className="form-input mt-0 sm:w-72" type="search" placeholder="Rechercher une carte" value={search} onChange={(event) => setSearch(event.target.value)} />
          <div className="flex rounded-2xl border border-slate-200 bg-white p-1">
            {filters.map((item) => (
              <button
                key={item.id}
                type="button"
                onClick={() => setFilter(item.id)}
                className={`rounded-xl px-3 py-2 text-sm font-semibold transition ${filter === item.id ? "bg-slate-950 text-white" : "text-slate-600 hover:bg-slate-100"}`}
              >
                {item.label}
              </button>
            ))}
          </div>
        </div>
      </div>

      {error && <ErrorMessage message={error} />}

      {isLoading ? (
        <LoadingState label="Chargement des cartes..." />
      ) : visibleCards.length > 0 ? (
        <div className="grid gap-5 md:grid-cols-2 xl:grid-cols-3">
          {visibleCards.map((card) => (
            <CardEventCard key={card.id} card={card} />
          ))}
        </div>
      ) : (
        <EmptyState title="Aucune carte trouvee" message="Essaie un autre filtre ou verifie que le backend retourne bien des cartes." />
      )}
    </div>
  );
}
