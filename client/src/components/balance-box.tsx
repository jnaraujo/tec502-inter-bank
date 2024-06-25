import { Button } from "./ui/button"

export function BalanceBox() {
  return (
    <article className="space-y-4 rounded-lg border bg-background p-6">
      <div className="space-y-1">
        <h2 className="text-lg font-medium text-zinc-500">Seu saldo:</h2>
        <p className="text-4xl font-bold text-green-600">R$ 250,00</p>
        <span className="text-sm font-medium text-zinc-400">nesse banco*</span>
      </div>
      <Button>Fazer depósito</Button>
    </article>
  )
}
