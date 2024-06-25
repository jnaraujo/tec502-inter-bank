/* eslint-disable react-refresh/only-export-components */
import { User } from "@/@types/user"
import { auth, createAccount } from "@/services/user-service"
import Cookies from "js-cookie"
import React, { createContext, useContext, useEffect, useState } from "react"

interface ILogin {
  document: string
}

interface ISignUp {
  name: string
  document: string
}

interface IAuthContext {
  user: User | null
  isAuthenticated: boolean
  login: (data: ILogin) => Promise<void>
  signUp: (data: ISignUp) => Promise<void>
}

export const AuthContext = createContext<IAuthContext>({} as IAuthContext)
export const useAuth = () => useContext(AuthContext)

function getUserFromCookies(): User | null {
  const cookieData = Cookies.get("user")
  if (!cookieData) return null
  try {
    return JSON.parse(cookieData) as User
  } catch (error) {
    Cookies.remove("user")
    return null
  }
}

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [user, setUser] = useState<User | null>(getUserFromCookies())
  const isAuthenticated = !!user

  async function login(data: ILogin) {
    const res = await auth(data.document)
    Cookies.set("user", JSON.stringify(res))
    setUser(res)
  }

  async function signUp(data: ISignUp) {
    const res = await createAccount(data)
    Cookies.set("user", JSON.stringify(res))
    setUser(res)
  }

  useEffect(() => {
    async function tryAuth() {
      const userCookie = getUserFromCookies()
      if (!userCookie) return
      try {
        await auth(userCookie.document)
      } catch (error) {
        Cookies.remove("user")
      }
    }
    tryAuth()
  }, [])

  return (
    <AuthContext.Provider value={{ user, signUp, login, isAuthenticated }}>
      {children}
    </AuthContext.Provider>
  )
}
