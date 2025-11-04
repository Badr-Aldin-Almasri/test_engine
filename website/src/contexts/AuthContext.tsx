import React, { createContext, useContext, useState, useEffect, ReactNode } from 'react'
import { User } from '@/types/auth'
import { authService } from '@/services/authService'

interface AuthContextType {
  user: User | null
  isAuthenticated: boolean
  isLoading: boolean
  login: (email: string, password: string) => Promise<void>
  register: (email: string, password: string, name: string) => Promise<void>
  logout: () => void
}

const AuthContext = createContext<AuthContextType | undefined>(undefined)

export const AuthProvider: React.FC<{ children: ReactNode }> = ({ children }) => {
  const [user, setUser] = useState<User | null>(null)
  const [isLoading, setIsLoading] = useState(true)

  useEffect(() => {
    const initAuth = () => {
      const savedUser = authService.getUser()
      const hasToken = authService.isAuthenticated()
      if (savedUser && hasToken) {
        setUser(savedUser)
      } else if (!hasToken) {
        // Clear invalid user data
        setUser(null)
      }
      setIsLoading(false)
    }
    initAuth()
  }, [])

  const login = async (email: string, password: string) => {
    const response = await authService.login({ email, password })
    setUser(response.user)
    // Ensure state is updated
    return response
  }

  const register = async (email: string, password: string, name: string) => {
    const response = await authService.register({ email, password, name })
    setUser(response.user)
    // Ensure state is updated
    return response
  }

  const logout = () => {
    authService.logout()
    setUser(null)
  }

  return (
    <AuthContext.Provider
      value={{
        user,
        isAuthenticated: !!user,
        isLoading,
        login,
        register,
        logout,
      }}
    >
      {children}
    </AuthContext.Provider>
  )
}

export const useAuth = () => {
  const context = useContext(AuthContext)
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider')
  }
  return context
}

