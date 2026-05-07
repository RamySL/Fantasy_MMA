import { useState } from "react";
import { Link, NavLink } from "react-router-dom";
import { useAuth } from "../context/AuthContext.jsx";

const navItems = [
  { to: "/cards", label: "Cartes" },
  { to: "/my-predictions", label: "Mes predictions" },
  { to: "/leaderboard", label: "Classement" },
];

function navClass({ isActive }) {
  return `rounded-xl px-3 py-2 text-sm font-semibold transition ${
    isActive ? "bg-slate-950 text-white" : "text-slate-600 hover:bg-slate-100 hover:text-slate-950"
  }`;
}

export default function Navbar() {
  const { user, isAuthenticated, isLoading, logout } = useAuth();
  const [isOpen, setIsOpen] = useState(false);

  async function handleLogout() {
    await logout();
    setIsOpen(false);
  }

  return (
    <header className="sticky top-0 z-30 border-b border-slate-200 bg-white/85 backdrop-blur">
      <nav className="app-container flex h-16 items-center justify-between gap-4">
        <Link to="/" className="flex items-center gap-3" onClick={() => setIsOpen(false)}>
          <span className="flex h-10 w-10 items-center justify-center rounded-2xl bg-slate-950 text-sm font-black text-white">FM</span>
          <span className="text-base font-black tracking-tight text-slate-950">Fantasy MAA</span>
        </Link>

        <div className="hidden items-center gap-1 md:flex">
          {navItems.map((item) => (
            <NavLink key={item.to} to={item.to} className={navClass}>
              {item.label}
            </NavLink>
          ))}
        </div>

        <div className="hidden items-center gap-3 md:flex">
          {isLoading ? (
            <span className="text-sm text-slate-500">Session...</span>
          ) : isAuthenticated ? (
            <>
              <span className="text-sm font-semibold text-slate-700">{user.pseudo}</span>
              <button type="button" onClick={handleLogout} className="btn-secondary">
                Deconnexion
              </button>
            </>
          ) : (
            <>
              <Link to="/login" className="btn-ghost">
                Connexion
              </Link>
              <Link to="/register" className="btn-primary">
                Inscription
              </Link>
            </>
          )}
        </div>

        <button type="button" className="btn-secondary md:hidden" onClick={() => setIsOpen((value) => !value)} aria-expanded={isOpen}>
          Menu
        </button>
      </nav>

      {isOpen && (
        <div className="border-t border-slate-200 bg-white md:hidden">
          <div className="app-container flex flex-col gap-2 py-4">
            {navItems.map((item) => (
              <NavLink key={item.to} to={item.to} className={navClass} onClick={() => setIsOpen(false)}>
                {item.label}
              </NavLink>
            ))}

            <div className="mt-3 border-t border-slate-200 pt-3">
              {isAuthenticated ? (
                <div className="flex flex-col gap-3">
                  <span className="text-sm font-semibold text-slate-700">Connecte: {user.pseudo}</span>
                  <button type="button" onClick={handleLogout} className="btn-secondary">
                    Deconnexion
                  </button>
                </div>
              ) : (
                <div className="grid gap-2">
                  <Link to="/login" className="btn-secondary" onClick={() => setIsOpen(false)}>
                    Connexion
                  </Link>
                  <Link to="/register" className="btn-primary" onClick={() => setIsOpen(false)}>
                    Inscription
                  </Link>
                </div>
              )}
            </div>
          </div>
        </div>
      )}
    </header>
  );
}
