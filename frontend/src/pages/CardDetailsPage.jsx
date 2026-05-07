import { useEffect, useMemo, useState } from "react";
import { Link, useNavigate, useParams } from "react-router-dom";
import Badge from "../components/Badge.jsx";
import EmptyState from "../components/EmptyState.jsx";
import ErrorMessage from "../components/ErrorMessage.jsx";
import FightPredictionCard from "../components/FightPredictionCard.jsx";
import LoadingState from "../components/LoadingState.jsx";
import StatCard from "../components/StatCard.jsx";
import { useAuth } from "../context/AuthContext.jsx";
import { cardService } from "../services/cardService.js";
import { predictionService } from "../services/predictionService.js";
import { formatDate } from "../utils/date.js";

function getLocation(card) {
  return [card.venue_name, card.city, card.region, card.country].filter(Boolean).join(", ");
}

export default function CardDetailsPage() {
  const { id } = useParams();
  const navigate = useNavigate();
  const { isAuthenticated } = useAuth();
  const [card, setCard] = useState(null);
  const [fights, setFights] = useState([]);
  const [selectedPredictions, setSelectedPredictions] = useState({});
  const [isLoading, setIsLoading] = useState(true);
  const [isSaving, setIsSaving] = useState(false);
  const [error, setError] = useState("");
  const [saveMessage, setSaveMessage] = useState("");

  useEffect(() => {
    let isMounted = true;

    async function loadCardDetails() {
      try {
        const [cardData, fightsData, savedPredictions] = await Promise.all([
          cardService.getCard(id),
          cardService.getCardFights(id),
          predictionService.getPredictionsForCard(id),
        ]);

        if (!isMounted) {
          return;
        }

        setCard(cardData);
        setFights(fightsData || []);

        // TODO backend: quand GET /cards/:id/predictions/me existera,
        // adapter la transformation ci-dessous selon la vraie reponse.
        const selected = (savedPredictions || []).reduce((acc, prediction) => {
          acc[prediction.fight_id] = prediction.predicted_winner_id;
          return acc;
        }, {});
        setSelectedPredictions(selected);
      } catch (err) {
        if (isMounted) {
          setError(err.message || "Impossible de charger cette carte.");
        }
      } finally {
        if (isMounted) {
          setIsLoading(false);
        }
      }
    }

    loadCardDetails();

    return () => {
      isMounted = false;
    };
  }, [id]);

  const summary = useMemo(() => {
    const selectedFightIds = Object.keys(selectedPredictions);
    const selectedPoints = fights.reduce((total, fight) => {
      return selectedPredictions[fight.id] ? total + (fight.points_good_prediction || 0) : total;
    }, 0);

    return {
      selectedCount: selectedFightIds.length,
      totalFights: fights.length,
      selectedPoints,
      availableFights: fights.filter((fight) => !fight.completed).length,
    };
  }, [fights, selectedPredictions]);

  function handleSelectPrediction(fightId, fighterId) {
    // Stockage sous forme d'objet indexe par fight_id pour rendre chaque
    // selection instantanee, meme si une carte contient beaucoup de combats.
    setSelectedPredictions((current) => ({
      ...current,
      [fightId]: fighterId,
    }));
    setSaveMessage("");
  }

  async function handleSavePredictions() {
    if (!isAuthenticated) {
      navigate("/login", { state: { from: `/cards/${id}` } });
      return;
    }

    setIsSaving(true);
    setSaveMessage("");
    setError("");

    try {
      await predictionService.saveCardPredictions(id, selectedPredictions);
      setSaveMessage("Predictions preparees cote front. TODO: brancher la sauvegarde backend.");
    } catch (err) {
      setError(err.message || "Impossible de sauvegarder les predictions.");
    } finally {
      setIsSaving(false);
    }
  }

  if (isLoading) {
    return <LoadingState label="Chargement de la carte..." />;
  }

  if (error && !card) {
    return <ErrorMessage message={error} />;
  }

  if (!card) {
    return <EmptyState title="Carte introuvable" message="Aucune carte ne correspond a cet identifiant." actionLabel="Retour aux cartes" actionTo="/cards" />;
  }

  const location = getLocation(card);

  return (
    <div className="space-y-6">
      <Link to="/cards" className="btn-ghost px-0">
        Retour aux cartes
      </Link>

      <section className="card-panel">
        <div className="flex flex-col gap-5 lg:flex-row lg:items-start lg:justify-between">
          <div>
            <div className="flex flex-wrap items-center gap-2">
              <Badge variant={card.completed ? "success" : "warning"}>{card.completed ? "Terminee" : card.status || "A venir"}</Badge>
              <Badge>Carte #{card.id}</Badge>
            </div>
            <h1 className="mt-4 text-3xl font-black tracking-tight text-slate-950">{card.title}</h1>
            <p className="mt-3 text-sm text-slate-600">{formatDate(card.date)}</p>
            <p className="mt-2 text-sm text-slate-500">{location || "Lieu a confirmer"}</p>
          </div>

          <button type="button" className="btn-primary" disabled={summary.selectedCount === 0 || isSaving} onClick={handleSavePredictions}>
            {isSaving ? "Sauvegarde..." : "Enregistrer mes predictions"}
          </button>
        </div>
      </section>

      {!isAuthenticated && (
        <ErrorMessage
          type="info"
          title="Connexion requise"
          message="Tu peux consulter la carte, mais il faudra te connecter pour sauvegarder tes predictions."
        />
      )}

      <ErrorMessage
        type="info"
        title="Vue predictions modelisee"
        message="TODO backend: brancher le chargement des predictions existantes et la sauvegarde reelle. Les selections fonctionnent deja cote front."
      />

      {error && <ErrorMessage message={error} />}
      {saveMessage && <ErrorMessage type="success" title="Pret pour le branchement" message={saveMessage} />}

      <section className="grid gap-4 sm:grid-cols-3">
        <StatCard label="Combats" value={summary.totalFights} helper={`${summary.availableFights} encore selectionnables`} />
        <StatCard label="Predictions choisies" value={`${summary.selectedCount}/${summary.totalFights}`} helper="Selection locale" />
        <StatCard label="Points potentiels" value={summary.selectedPoints} helper="Selon tes choix actuels" />
      </section>

      <section className="space-y-4">
        <div>
          <h2 className="section-title">Combats</h2>
          <p className="mt-1 text-sm text-slate-500">Selectionne un combattant par combat.</p>
        </div>

        {fights.length > 0 ? (
          <div className="space-y-4">
            {fights.map((fight) => (
              <FightPredictionCard
                key={fight.id}
                fight={fight}
                selectedWinnerId={selectedPredictions[fight.id]}
                disabled={!isAuthenticated}
                onSelect={handleSelectPrediction}
              />
            ))}
          </div>
        ) : (
          <EmptyState title="Aucun combat trouve" message="Les combats apparaitront ici quand /cards/:id/fights retournera des donnees." />
        )}
      </section>
    </div>
  );
}
