import { useAuth } from "@/contexts/auth-context"
import { useBalance } from "@/hooks/use-user"
import { Button } from "./ui/button"

export function BalanceBox() {
  const { user } = useAuth()
  const { data } = useBalance(user?.id)

  return (
    <article className="space-y-4 rounded-lg border bg-background p-6">
      <div className="space-y-1">
        <div className="relative">
          <h2 className="text-lg font-medium text-zinc-500">Seu saldo:</h2>
          <span className="absolute right-0 top-0 text-xs text-zinc-500">
            {data.ts.toLocaleString("pt-br")}
          </span>
        </div>{" "}
        <p className="text-4xl font-bold text-green-600">
          R$ {data.balance.toFixed(2)}
        </p>
        <p className="text-sm font-medium text-zinc-400">
          *valor do saldo nesse banco
        </p>
        <p className="text-sm font-medium text-zinc-400">
          **valor pode estar desatualizado
        </p>
      </div>
      <Button>Fazer depósito</Button>
    </article>
  )
}
