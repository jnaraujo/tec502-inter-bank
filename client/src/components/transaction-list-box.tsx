import { useAuth } from "@/contexts/auth-context"
import { useTransactions } from "@/hooks/use-transactions"
import { TransactionBox } from "./transaction-box"

export function TransactionListBox() {
  const { user } = useAuth()
  const { data: transactions } = useTransactions(user?.id)

  return (
    <article className="flex flex-col justify-between space-y-4 overflow-auto rounded-lg border bg-background p-6">
      <h2 className="text-xl font-medium text-zinc-800">Suas transações</h2>

      <div className="flex flex-col gap-6 overflow-auto py-2">
        {transactions.map((tr) => (
          <TransactionBox
            key={tr.id}
            createdAt={tr.created_at}
            status={tr.status as any}
            operations={tr.operations.map((op) => ({
              amount: Number(op.amount),
              from: op.from,
              to: op.to,
            }))}
          />
        ))}
      </div>
    </article>
  )
}
