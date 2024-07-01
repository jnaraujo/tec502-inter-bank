/* eslint-disable react-refresh/only-export-components */
import { User } from "@/@types/user"
import { auth, createAccount } from "@/services/user-service"
import React, { createContext, useContext, useState } from "react"

interface ILogin {
  ibk: string
}

interface ISignUp {
  name: string
  documents: string[]
  type: "individual" | "legal" | "joint"
}

interface IAuthContext {
  user: User | null | undefined
  isAuthenticated: boolean
  login: (data: ILogin) => Promise<void>
  signUp: (data: ISignUp) => Promise<void>
}

export const AuthContext = createContext<IAuthContext>({} as IAuthContext)
export const useAuth = () => useContext(AuthContext)

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [user, setUser] = useState<User | null | undefined>(undefined)
  const isAuthenticated = !!user

  async function login(data: ILogin) {
    const res = await auth(data.ibk)
    setUser(res)
  }

  async function signUp(data: ISignUp) {
    const res = await createAccount(data)
    setUser(res)
  }

  return (
    <AuthContext.Provider value={{ user, signUp, login, isAuthenticated }}>
      {children}
    </AuthContext.Provider>
  )
}
