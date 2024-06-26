import { BALANCE_REFETCH_INTERVAL } from "@/constants/query"
import { env } from "@/env"
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query"

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

interface DepositData {
  amount: number
  userId: number
}

export function useDeposit() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: async (data: DepositData) => {
      const response = await fetch(
        `${env.VITE_BANK_URL}/api/payments/deposit`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify({
            user_id: data.userId,
            amount: data.amount,
          }),
        },
      )

      if (!response.ok) {
        throw new Error(
          (await response.json()).message || "Erro ao realizar a transação.",
        )
      }

      const res = await response.json()
      return {
        message: res.message,
      } satisfies {
        message: string
      }
    },
    onSuccess: () => {
      queryClient.invalidateQueries({
        queryKey: ["balance", "transactions"],
      })
    },
  })
}
