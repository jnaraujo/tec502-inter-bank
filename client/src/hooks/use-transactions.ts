import { Transaction } from "@/@types/transaction"
import { TRANSACTIONS_REFETCH_INTERVAL } from "@/constants/query"
import { useBank } from "@/stores/bank-store"
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query"

export function useTransactions(userId?: number) {
  return useQuery({
    queryFn: async ({ signal }) => {
      if (!userId) {
        return []
      }

      const address = useBank.getState().address
      const resp = await fetch(
        `${address}/api/accounts/${userId}/transactions`,
        {
          signal,
        },
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
      const address = useBank.getState().address
      const response = await fetch(`${address}/api/payments/pay`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          author: transaction.author,
          operations: transaction.operations.map((op) => ({
            from: op.from,
            to: op.to,
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
