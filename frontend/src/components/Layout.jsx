import { Outlet } from "react-router-dom";
import Navbar from "./Navbar.jsx";

export default function Layout() {
  return (
    <div className="flex min-h-screen flex-col">
      <Navbar />
      <main className="app-container flex-1 py-8 sm:py-10">
        <Outlet />
      </main>
      <footer className="border-t border-slate-200 bg-white/70 py-6">
        <div className="app-container flex flex-col gap-2 text-sm text-slate-500 sm:flex-row sm:items-center sm:justify-between">
          <span>Fantasy MAA</span>
          <span>Predictions MMA, classement et cartes de combats.</span>
        </div>
      </footer>
    </div>
  );
}
