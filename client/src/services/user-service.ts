import { User } from "@/@types/user"
import { useBank } from "@/stores/bank-store"

export async function auth(ibk: string) {
  const address = useBank.getState().address
  const response = await fetch(`${address}/api/accounts/auth`, {
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
  const address = useBank.getState().address
  const response = await fetch(`${address}/api/accounts`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(account),
  })

  if (!response.ok) {
    throw new Error((await response.json()).message)
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
