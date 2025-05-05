import { Routes, Route } from "react-router";

import HomePage from "./pages/home";
import AuthPage from "./pages/auth";
import ProfilePage from "./pages/profilePage";

function Router() {
  return (
    <Routes>
      <Route path="/" element={<HomePage />} />
      <Route path="auth/callback/" element={<AuthPage />} />
      <Route path="profile/:username/" element={<ProfilePage />} />
    </Routes>
  );
}

export default Router;
