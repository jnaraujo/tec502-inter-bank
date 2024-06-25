import { TRANSACTIONS_REFETCH_INTERVAL } from "@/constants/query"
import { env } from "@/env"
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query"

export interface Operation {
  id: string
  from: string
  to: string
  type: string
  amount: string
  status: string
  created_at: string
  updated_at: string
}

export interface Transaction {
  id: string
  author: string
  operations: Operation[]
  created_at: string
  updated_at: string
  status: string
}

export function useTransactions(userId?: number) {
  return useQuery({
    queryFn: async () => {
      if (!userId) {
        return []
      }

      const resp = await fetch(
        `${env.VITE_BANK_URL}/api/accounts/${userId}/transactions`,
      )

      if (!resp.ok) {
        throw new Error("Não foi possível pegar os dados do usuário.")
      }

      return (await resp.json()) as Array<Transaction>
    },
    initialData: [],
    refetchInterval: TRANSACTIONS_REFETCH_INTERVAL,
    queryKey: ["transactions", userId],
  })
}

interface NewTransaction {
  author: string
  operations: Array<{
    from: string
    to: string
    amount: number
  }>
}

export function useSendTransaction() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: async (transaction: NewTransaction) => {
      const response = await fetch(`${env.VITE_BANK_URL}/api/payments/pay`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          author: transaction.author,
          operations: transaction.operations.map((op) => ({
            from_user_ibk: op.from,
            to_user_ibk: op.to,
            amount: op.amount,
          })),
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
        queryKey: ["transactions"],
      })
    },
  })
}
