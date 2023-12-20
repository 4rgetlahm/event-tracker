import React from 'react'
import { BrowserRouter, Route, Routes } from 'react-router-dom'
import Landing from './pages/Landing'
import DependencyInjector from './components/DependencyInjector'
import { LoginCallback } from './pages/auth/LoginCallback'
import ProtectedRoute from './components/ProtectedRoutes'
import { Navbar } from './components/Navbar'

function App() {
    return (
        <BrowserRouter>
            <DependencyInjector>
                <>
                    <Navbar />
                    <Routes>
                        <Route
                            path="/oauth2/callback"
                            element={<LoginCallback />}
                        />
                        <Route
                            path="/"
                            element={
                                <ProtectedRoute>
                                    <Landing />
                                </ProtectedRoute>
                            }
                        />
                    </Routes>
                </>
            </DependencyInjector>
        </BrowserRouter>
    )
}

export default App
