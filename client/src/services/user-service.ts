import { User } from "@/@types/user"
import { env } from "@/env"

export async function auth(id: string) {
  const response = await fetch(`${env.VITE_BANK_URL}/api/accounts/${id}`)
  if (!response.ok) {
    throw new Error("Failed to fetch user")
  }

  const res = await response.json()
  return {
    id: res.id,
    name: res.name,
    balance: res.balance,
    createdAt: new Date(res.created_at),
    ibk: res.ibk,
  } as User
}

interface CreateAccountUser {
  name: string
}

export async function createAccount(user: CreateAccountUser) {
  const response = await fetch(`${env.VITE_BANK_URL}/api/accounts`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(user),
  })
  if (!response.ok) {
    throw new Error("Failed to signup user")
  }
  const res = await response.json()
  return {
    id: res.id,
    name: res.name,
    balance: res.balance,
    createdAt: new Date(res.created_at),
    ibk: res.ibk,
  } as User
}
