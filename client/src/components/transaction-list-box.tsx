import { Transaction as ITransaction } from "@/@types/transaction"
import { useAuth } from "@/contexts/auth-context"
import { useTransactions } from "@/hooks/use-transactions"
import { Transaction } from "./transaction"

export function TransactionListBox() {
  const { user } = useAuth()
  const { data: transactions } = useTransactions(user?.id)

  function formatId(tr: ITransaction) {
    if (tr.type === "final" && tr.parent_id) return tr.parent_id
    return tr.id
  }

  return (
    <article className="flex flex-col space-y-4 overflow-auto rounded-lg border bg-background p-6">
      <h2 className="text-xl font-medium text-zinc-800">Suas transações</h2>

      <div className="flex flex-col gap-6 overflow-auto py-2">
        {transactions.map((tr) => (
          <Transaction
            key={tr.id}
            txId={formatId(tr)}
            createdAt={tr.created_at}
            status={tr.status as any}
            type={tr.type}
            operations={tr.operations.map((op) => ({
              amount: Number(op.amount),
              from: op.from,
              to: op.to,
              type: op.type,
            }))}
          />
        ))}
      </div>
    </article>
  )
}
