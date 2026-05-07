const variants = {
  error: {
    wrapper: "border-rose-200 bg-rose-50 text-rose-800",
    title: "text-rose-950",
  },
  info: {
    wrapper: "border-sky-200 bg-sky-50 text-sky-800",
    title: "text-sky-950",
  },
  success: {
    wrapper: "border-emerald-200 bg-emerald-50 text-emerald-800",
    title: "text-emerald-950",
  },
};

export default function ErrorMessage({ title = "Erreur", message, type = "error" }) {
  const theme = variants[type] || variants.error;

  return (
    <div className={`rounded-2xl border p-4 ${theme.wrapper}`}>
      <p className={`text-sm font-bold ${theme.title}`}>{title}</p>
      {message && <p className="mt-1 text-sm">{message}</p>}
    </div>
  );
}
