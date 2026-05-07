export default function LoadingState({ label = "Chargement..." }) {
  return (
    <div className="flex min-h-40 items-center justify-center rounded-3xl border border-dashed border-slate-300 bg-white/80 p-8 text-center">
      <div>
        <div className="mx-auto h-10 w-10 animate-spin rounded-full border-4 border-slate-200 border-t-slate-950" />
        <p className="mt-4 text-sm font-medium text-slate-600">{label}</p>
      </div>
    </div>
  );
}
