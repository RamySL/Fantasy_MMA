import { useState } from "react";
import { Link, useLocation, useNavigate } from "react-router-dom";
import ErrorMessage from "../components/ErrorMessage.jsx";
import { useAuth } from "../context/AuthContext.jsx";

function getRedirectPath(locationState) {
  if (!locationState?.from) {
    return "/";
  }

  if (typeof locationState.from === "string") {
    return locationState.from;
  }

  return locationState.from.pathname || "/";
}

export default function LoginPage() {
  const navigate = useNavigate();
  const location = useLocation();
  const { login } = useAuth();
  const [form, setForm] = useState({ email: "", password: "" });
  const [error, setError] = useState("");
  const [isSubmitting, setIsSubmitting] = useState(false);

  function updateField(event) {
    const { name, value } = event.target;
    setForm((current) => ({ ...current, [name]: value }));
  }

  async function handleSubmit(event) {
    event.preventDefault();
    setError("");
    setIsSubmitting(true);

    try {
      await login(form);
      navigate(getRedirectPath(location.state), { replace: true });
    } catch (err) {
      setError(err.message || "Connexion impossible.");
    } finally {
      setIsSubmitting(false);
    }
  }

  return (
    <div className="mx-auto max-w-md">
      <div className="card-panel">
        <div>
          <p className="text-sm font-semibold uppercase tracking-[0.25em] text-slate-500">Connexion</p>
          <h1 className="mt-2 text-3xl font-black tracking-tight text-slate-950">Bon retour</h1>
          <p className="mt-2 text-sm leading-6 text-slate-600">Connecte-toi pour sauvegarder tes predictions.</p>
        </div>

        {error && <div className="mt-5"><ErrorMessage message={error} /></div>}

        <form onSubmit={handleSubmit} className="mt-6 space-y-5">
          <label className="block">
            <span className="form-label">Email</span>
            <input className="form-input" type="email" name="email" autoComplete="email" required value={form.email} onChange={updateField} />
          </label>

          <label className="block">
            <span className="form-label">Mot de passe</span>
            <input className="form-input" type="password" name="password" autoComplete="current-password" required value={form.password} onChange={updateField} />
          </label>

          <button type="submit" className="btn-primary w-full" disabled={isSubmitting}>
            {isSubmitting ? "Connexion..." : "Se connecter"}
          </button>
        </form>

        <p className="mt-6 text-center text-sm text-slate-500">
          Pas encore de compte ? <Link to="/register" className="font-bold text-slate-950 hover:underline">Creer un compte</Link>
        </p>
      </div>
    </div>
  );
}
