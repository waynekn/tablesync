import { Routes, Route } from "react-router";

import HomePage from "./pages/home";
import AuthPage from "./pages/auth";
import ProfilePage from "./pages/profilePage";
import CreateSheet from "./pages/createSpreadSheet";
import EditSheet from "./pages/sheetEdit";

function Router() {
  return (
    <Routes>
      <Route path="/" element={<HomePage />} />
      <Route path="auth/callback/" element={<AuthPage />} />
      <Route path="profile/:username/" element={<ProfilePage />} />
      <Route path="spreadsheet/create/" element={<CreateSheet />} />
      <Route path="sheet/:sheetID/edit/" element={<EditSheet />} />
    </Routes>
  );
}

export default Router;
