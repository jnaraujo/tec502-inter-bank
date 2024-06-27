import { User } from "@/@types/user"
import { env } from "@/env"

export async function auth(ibk: string) {
  const response = await fetch(`${env.VITE_BANK_URL}/api/accounts/auth`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ acc_ibk: ibk }),
  })
  if (response.status == 404) {
    throw new Error("Usuário não encontrado")
  }

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
    document: res.document,
    type: res.type,
  } satisfies User
}

interface CreateAccountUser {
  name: string
}

export async function createAccount(account: CreateAccountUser) {
  const response = await fetch(`${env.VITE_BANK_URL}/api/accounts`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(account),
  })

  if (response.status == 409) {
    throw new Error("Usuário já existe")
  }

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
    document: res.document,
    type: res.type,
  } satisfies User
}
