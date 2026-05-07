import EmptyState from "../components/EmptyState.jsx";

export default function NotFoundPage() {
  return (
    <EmptyState
      title="Page introuvable"
      message="La page demandee n'existe pas ou a ete deplacee."
      actionLabel="Retour a l'accueil"
      actionTo="/"
    />
  );
}
