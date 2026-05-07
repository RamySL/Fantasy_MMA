import { useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import ErrorMessage from "../components/ErrorMessage.jsx";
import { useAuth } from "../context/AuthContext.jsx";

export default function RegisterPage() {
  const navigate = useNavigate();
  const { register } = useAuth();
  const [form, setForm] = useState({ pseudo: "", email: "", password: "" });
  const [error, setError] = useState("");
  const [isSubmitting, setIsSubmitting] = useState(false);

  function updateField(event) {
    const { name, value } = event.target;
    setForm((current) => ({ ...current, [name]: value }));
  }

  async function handleSubmit(event) {
    event.preventDefault();
    setError("");

    if (form.password.length < 6) {
      setError("Le mot de passe doit contenir au moins 6 caracteres.");
      return;
    }

    setIsSubmitting(true);

    try {
      await register(form);
      navigate("/", { replace: true });
    } catch (err) {
      setError(err.message || "Inscription impossible.");
    } finally {
      setIsSubmitting(false);
    }
  }

  return (
    <div className="mx-auto max-w-md">
      <div className="card-panel">
        <div>
          <p className="text-sm font-semibold uppercase tracking-[0.25em] text-slate-500">Inscription</p>
          <h1 className="mt-2 text-3xl font-black tracking-tight text-slate-950">Creer ton compte</h1>
          <p className="mt-2 text-sm leading-6 text-slate-600">Un pseudo, un email et tu peux commencer a jouer.</p>
        </div>

        {error && <div className="mt-5"><ErrorMessage message={error} /></div>}

        <form onSubmit={handleSubmit} className="mt-6 space-y-5">
          <label className="block">
            <span className="form-label">Pseudo</span>
            <input className="form-input" type="text" name="pseudo" autoComplete="nickname" required value={form.pseudo} onChange={updateField} />
          </label>

          <label className="block">
            <span className="form-label">Email</span>
            <input className="form-input" type="email" name="email" autoComplete="email" required value={form.email} onChange={updateField} />
          </label>

          <label className="block">
            <span className="form-label">Mot de passe</span>
            <input className="form-input" type="password" name="password" autoComplete="new-password" required value={form.password} onChange={updateField} />
          </label>

          <button type="submit" className="btn-primary w-full" disabled={isSubmitting}>
            {isSubmitting ? "Creation..." : "Creer le compte"}
          </button>
        </form>

        <p className="mt-6 text-center text-sm text-slate-500">
          Deja inscrit ? <Link to="/login" className="font-bold text-slate-950 hover:underline">Se connecter</Link>
        </p>
      </div>
    </div>
  );
}
