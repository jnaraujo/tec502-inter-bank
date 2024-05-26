/* eslint-disable react-refresh/only-export-components */
import React, { createContext, useContext } from "react"

interface User {
  name: string
  key: string
}

interface IAuthContext {
  user: User | null
}

export const AuthContext = createContext<IAuthContext>({} as IAuthContext)
export const useAuth = () => useContext(AuthContext)

export function AuthProvider({ children }: { children: React.ReactNode }) {
  return (
    <AuthContext.Provider value={{ user: null }}>
      {children}
    </AuthContext.Provider>
  )
}
