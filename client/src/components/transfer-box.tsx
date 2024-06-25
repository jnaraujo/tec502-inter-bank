import { useAuth } from "@/contexts/auth-context"
import { sendTransaction } from "@/services/payment-service"
import { ArrowLeftRight, DollarSign, Plus, Trash } from "lucide-react"
import { useState } from "react"
import { CreateTransactionDialog, Operation } from "./create-transaction-dialog"
import { Button } from "./ui/button"
import { toast } from "./ui/use-toast"

export function TransferBox() {
  const { user } = useAuth()
  const [operations, setOperations] = useState<Array<Operation>>([])
  const [openDialog, setOpenDialog] = useState(false)

  function handleNewOperation(operation: Operation) {
    setOperations((prev) => [operation, ...prev])
  }

  function removeOperation(idx: number) {
    setOperations((prev) => prev.filter((_, i) => i !== idx))
  }

  async function handleSubmit(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault()
    if (!user) return

    if (operations.length === 0) {
      toast({
        title: "Adicione uma operação.",
        description:
          "Você precisa de pelo meno uma operação para realizar a transferência.",
        variant: "destructive",
      })
      return
    }

    try {
      await sendTransaction({
        author: user.ibk,
        operations: operations.map((op) => ({
          from: user.ibk,
          to: op.to,
          amount: op.amount,
        })),
      })

      setOperations([])

      toast({
        title: "Transação criada com sucesso!",
        description:
          'Você pode acompanhar o status da transação pela caixa "Suas transações".',
      })
    } catch (error) {
      toast({
        title: "Erro ao criar a transação.",
        description: (error as any).message,
        variant: "destructive",
      })
    }
  }

  return (
    <>
      <article className="space-y-4 rounded-lg border bg-background p-6">
        <h2 className="text-lg font-medium text-zinc-800">
          Realizar transferência
        </h2>
        <p className="text-zinc-500">
          Sua transferência pode ter várias operações! Se uma delas falhar,
          todas são desfeitas.
        </p>

        <form className="space-y-4" onSubmit={handleSubmit}>
          <h3 className="font-medium text-zinc-600">
            Operações da transferência
          </h3>
          <div>
            <div className="h-64 space-y-2 overflow-y-auto overflow-x-hidden px-4">
              {operations.length === 0 ? (
                <p className="p-6 text-center text-sm text-zinc-500">
                  Sem operações ainda. Adicione uma nova operação para aparecer
                  aqui.
                </p>
              ) : null}
              {operations.map((op, idx) => (
                <div className="flex items-center gap-2 text-zinc-400">
                  <div className="flex w-full items-center gap-2">
                    <div className="flex items-center gap-1">
                      <ArrowLeftRight className="size-4 text-zinc-600" />
                      <span className="min-w-20 text-zinc-500">
                        Para: {op.to}
                      </span>
                    </div>
                    <div className="flex items-center gap-1">
                      <DollarSign className="size-4 text-green-600" />
                      <span className="text-green-500">
                        R$ {op.amount.toFixed(2)}
                      </span>
                    </div>
                  </div>
                  <button
                    onClick={() => {
                      removeOperation(idx)
                    }}
                  >
                    <Trash className="size-5 text-red-400 transition-colors duration-200 hover:text-red-500" />
                  </button>
                </div>
              ))}
            </div>
            <Button
              className="mt-4 flex items-center gap-1"
              variant="secondary"
              onClick={() => {
                setOpenDialog(true)
              }}
            >
              <Plus /> Adicionar nova Operação
            </Button>
          </div>

          <Button type="submit">Realizar transferência</Button>
        </form>
      </article>

      <CreateTransactionDialog
        open={openDialog}
        onOpenChange={setOpenDialog}
        onOperationCreated={handleNewOperation}
      />
    </>
  )
}
