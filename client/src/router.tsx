import { Routes, Route } from "react-router";

import HomePage from "./pages/home";
import AuthPage from "./pages/auth";

function Router() {
  return (
    <Routes>
      <Route path="/" element={<HomePage />} />
      <Route path="auth/callback/" element={<AuthPage />} />
    </Routes>
  );
}

export default Router;
