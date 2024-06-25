import { BALANCE_REFETCH_INTERVAL } from "@/constants/query"
import { env } from "@/env"
import { useQuery } from "@tanstack/react-query"

export function useBalance(userId?: number) {
  return useQuery({
    queryFn: async () => {
      if (!userId) {
        return {
          balance: 0,
          ts: new Date(),
        }
      }

      const resp = await fetch(`${env.VITE_BANK_URL}/api/accounts/${userId}`)

      if (!resp.ok) {
        throw new Error("Não foi possível pegar os dados do usuário.")
      }

      return {
        balance: Number((await resp.json()).balance),
        ts: new Date(),
      }
    },
    initialData: {
      balance: 0,
      ts: new Date(),
    },
    refetchInterval: BALANCE_REFETCH_INTERVAL,
    queryKey: ["balance", userId],
  })
}
