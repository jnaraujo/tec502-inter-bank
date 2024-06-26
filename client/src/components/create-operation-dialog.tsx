import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog"
import { useAuth } from "@/contexts/auth-context"
import { handleZodError } from "@/lib/zod"
import { createOperationFormSchema } from "@/schemas/create-operation"
import { useState } from "react"
import { Button } from "./ui/button"
import { Input } from "./ui/input"
import { Label } from "./ui/label"

export interface Operation {
  from: string
  to: string
  amount: number
}

interface Props {
  open: boolean
  onOpenChange: (open: boolean) => void
  onOperationCreated: (operation: Operation) => void
}

export function CreateOperationDialog(props: Props) {
  const { user } = useAuth()
  const [errors, setErrors] = useState<string[]>([])

  function handleSendCommand(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault()

    const formData = new FormData(event.currentTarget)

    try {
      const data = createOperationFormSchema.parse({
        from: formData.get("from"),
        to: formData.get("to"),
        amount: formData.get("amount"),
      })

      props.onOperationCreated(data)
      props.onOpenChange(false)
    } catch (error: any) {
      setErrors(handleZodError(error))
    }
  }

  return (
    <Dialog open={props.open} onOpenChange={props.onOpenChange}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Adicionar nova operação</DialogTitle>
          <DialogDescription>
            Todas as operação atômicas. Ou seja, se uma falhar, todas serão
            desfeitas (até as que tiveram sucesso).
          </DialogDescription>
        </DialogHeader>
        <div>
          <form className="space-y-4" onSubmit={handleSendCommand}>
            <div className="space-y-1">
              <Label htmlFor="from">IBK do pagador:</Label>
              <Input
                id="from"
                name="from"
                placeholder="Ex: 1-203"
                required
                defaultValue={user?.ibk}
              />
            </div>

            <div className="space-y-1">
              <Label htmlFor="to">IBK do beneficiário:</Label>
              <Input id="to" name="to" placeholder="Ex: 1-203" required />
            </div>
            <div className="space-y-1">
              <Label htmlFor="amount">Valor:</Label>
              <Input
                id="amount"
                name="amount"
                placeholder="Ex: 500"
                type="number"
                step={0.01}
                min={0}
                required
              />
            </div>

            {errors.length !== 0 && (
              <div>
                {errors.map((error) => (
                  <div key={error} className="text-sm text-red-500">
                    {error}
                  </div>
                ))}
              </div>
            )}
            <Button type="submit">Adicionar operação</Button>
          </form>
        </div>
      </DialogContent>
    </Dialog>
  )
}
