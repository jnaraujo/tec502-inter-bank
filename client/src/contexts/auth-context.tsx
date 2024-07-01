/* eslint-disable react-refresh/only-export-components */
import { User } from "@/@types/user"
import { auth, createAccount } from "@/services/user-service"
import Cookies from "js-cookie"
import React, { createContext, useContext, useEffect, useState } from "react"

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

function getUserFromCookies(): User | null | undefined {
  const cookieData = Cookies.get("user")
  if (!cookieData) return undefined
  try {
    return JSON.parse(cookieData) as User
  } catch (error) {
    Cookies.remove("user")
    return null
  }
}

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [user, setUser] = useState<User | null | undefined>(undefined)
  const isAuthenticated = !!user

  async function login(data: ILogin) {
    const res = await auth(data.ibk)
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
      if (!userCookie) {
        setUser(null)
        return
      }
      try {
        await auth(userCookie.ibk)
        setUser(userCookie)
      } catch (error) {
        Cookies.remove("user")
        setUser(null)
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
