import { useAuth } from "@/contexts/auth-context"
import { useBalance, useDeposit } from "@/hooks/use-user"
import { useState } from "react"
import { DepositDialog } from "./deposit-dialog"
import { Button } from "./ui/button"
import { toast } from "./ui/use-toast"

export function BalanceBox() {
  const { user } = useAuth()
  const { data } = useBalance(user?.id)
  const { mutate: deposit } = useDeposit()
  const [open, setOpen] = useState(false)

  function handleSubmit(amount: number) {
    if (!user) return

    deposit(
      {
        amount,
        userIBK: user.ibk,
      },
      {
        onSuccess: () => {
          toast({
            title: "Depósito criado com sucesso!",
            description:
              'Você pode acompanhar o status da transação pela caixa "Suas transações"..',
          })
        },
        onError: (error) => {
          toast({
            title: "Erro ao criar o depósito.",
            description: (error as any).message,
            variant: "destructive",
          })
        },
      },
    )
  }

  return (
    <>
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
        <Button
          onClick={() => {
            setOpen(true)
          }}
        >
          Fazer depósito
        </Button>
      </article>

      <DepositDialog
        open={open}
        onOpenChange={setOpen}
        onSubmit={handleSubmit}
      />
    </>
  )
}
