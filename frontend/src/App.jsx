import { Navigate, Route, Routes } from "react-router-dom";
import Layout from "./components/Layout.jsx";
import ProtectedRoute from "./components/ProtectedRoute.jsx";
import CardDetailsPage from "./pages/CardDetailsPage.jsx";
import CardsPage from "./pages/CardsPage.jsx";
import HomePage from "./pages/HomePage.jsx";
import LeaderboardPage from "./pages/LeaderboardPage.jsx";
import LoginPage from "./pages/LoginPage.jsx";
import MyPredictionsPage from "./pages/MyPredictionsPage.jsx";
import NotFoundPage from "./pages/NotFoundPage.jsx";
import RegisterPage from "./pages/RegisterPage.jsx";

export default function App() {
  return (
    <Routes>
      <Route element={<Layout />}>
        <Route index element={<HomePage />} />
        <Route path="cards" element={<CardsPage />} />
        <Route path="cards/:id" element={<CardDetailsPage />} />
        <Route path="leaderboard" element={<LeaderboardPage />} />

        <Route element={<ProtectedRoute />}>
          <Route path="my-predictions" element={<MyPredictionsPage />} />
        </Route>

        <Route path="login" element={<LoginPage />} />
        <Route path="register" element={<RegisterPage />} />
        <Route path="home" element={<Navigate to="/" replace />} />
        <Route path="*" element={<NotFoundPage />} />
      </Route>
    </Routes>
  );
}
