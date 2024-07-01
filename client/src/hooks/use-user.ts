import { BALANCE_REFETCH_INTERVAL } from "@/constants/query"
import { useBank } from "@/stores/bank-store"
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query"

export function useBalance(userId?: number) {
  return useQuery({
    queryFn: async ({ signal }) => {
      if (!userId) {
        return {
          balance: 0,
          ts: new Date(),
        }
      }

      const address = useBank.getState().address
      const resp = await fetch(`${address}/api/accounts/${userId}`, {
        signal,
      })
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
  userIBK: string
  amount: number
}

export function useDeposit() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: async (data: DepositData) => {
      const address = useBank.getState().address
      const response = await fetch(`${address}/api/payments/deposit`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          acc_ibk: data.userIBK,
          amount: data.amount,
        }),
      })

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
