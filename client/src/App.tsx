import React from "react";
import { BrowserRouter, Route, Routes } from "react-router-dom";
import Landing from "./pages/Landing";
import DependencyInjector from "./components/DependencyInjector";
import { LoginCallback } from "./pages/auth/LoginCallback";
import ProtectedRoute from "./components/ProtectedRoutes";

function App() {
  return (
    <DependencyInjector>
        <BrowserRouter>
            <Routes>
                <Route path="/oauth2/callback" element={<LoginCallback />} />
                  <Route path="/" element={<ProtectedRoute><Landing /></ProtectedRoute>} />
            </Routes>
        </BrowserRouter>
    </DependencyInjector>
  );
}

export default App;
