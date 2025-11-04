import React from 'react'
import { BrowserRouter as Router, Routes, Route, Navigate, useParams } from 'react-router-dom'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { AuthProvider, useAuth } from './contexts/AuthContext'
import { FlowEditor } from './components/FlowEditor'
import { Dashboard } from './components/Dashboard'
import { Login } from './components/Login'

const queryClient = new QueryClient()

const ProtectedRoute: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const { isAuthenticated, isLoading } = useAuth()

  // Check localStorage directly for immediate authentication check
  const token = localStorage.getItem('auth_token')
  const hasToken = !!token

  if (isLoading && !hasToken) {
    return (
      <div className="flex items-center justify-center h-screen">
        <div>Loading...</div>
      </div>
    )
  }

  // Allow access if either React state says authenticated OR token exists
  if (!isAuthenticated && !hasToken) {
    return <Navigate to="/login" replace />
  }

  return <>{children}</>
}

// Wrapper component for FlowEditor to get flowId from URL
const FlowEditorRoute: React.FC = () => {
  const { flowId } = useParams<{ flowId: string }>()
  
  if (!flowId) {
    return (
      <div className="flex items-center justify-center h-full">
        <div className="text-center">
          <p className="text-muted-foreground">Invalid flow ID</p>
        </div>
      </div>
    )
  }

  return <FlowEditor flowId={flowId} />
}

function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <AuthProvider>
        <Router>
          <Routes>
            <Route path="/login" element={<Login />} />
            <Route
              path="/"
              element={
                <ProtectedRoute>
                  <Dashboard />
                </ProtectedRoute>
              }
            />
            <Route
              path="/flow/:flowId/editor"
              element={
                <ProtectedRoute>
                  <FlowEditorRoute />
                </ProtectedRoute>
              }
            />
            <Route path="*" element={<Navigate to="/" replace />} />
          </Routes>
        </Router>
      </AuthProvider>
    </QueryClientProvider>
  )
}

export default App
