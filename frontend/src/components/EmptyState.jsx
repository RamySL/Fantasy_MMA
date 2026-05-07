import { Link } from "react-router-dom";

export default function EmptyState({ title, message, actionLabel, actionTo }) {
  return (
    <div className="rounded-3xl border border-dashed border-slate-300 bg-white/80 p-8 text-center">
      <p className="text-lg font-bold text-slate-950">{title}</p>
      {message && <p className="mx-auto mt-2 max-w-xl text-sm leading-6 text-slate-600">{message}</p>}
      {actionLabel && actionTo && (
        <Link to={actionTo} className="btn-primary mt-5">
          {actionLabel}
        </Link>
      )}
    </div>
  );
}
