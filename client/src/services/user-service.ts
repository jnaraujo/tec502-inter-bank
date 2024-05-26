import { env } from "@/env"

export async function auth(id: string) {
  const response = await fetch(`${env}/api/accounts/${id}`)
  if (!response.ok) {
    throw new Error("Failed to fetch user")
  }
  return response.json()
}
