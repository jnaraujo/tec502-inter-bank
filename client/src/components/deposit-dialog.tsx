import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog"
import { handleZodError } from "@/lib/zod"
import { createDepositFormSchema } from "@/schemas/deposit"
import { useState } from "react"
import { Button } from "./ui/button"
import { Input } from "./ui/input"
import { Label } from "./ui/label"

interface Props {
  open: boolean
  onOpenChange: (open: boolean) => void
  onSubmit: (amount: number) => void
}

export function DepositDialog(props: Props) {
  const [errors, setErrors] = useState<string[]>([])

  function handleSendCommand(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault()
    const formData = new FormData(event.currentTarget)

    try {
      const data = createDepositFormSchema.parse({
        amount: formData.get("amount"),
      })

      props.onSubmit(data.amount)
      props.onOpenChange(false)
    } catch (error: any) {
      setErrors(handleZodError(error))
    }
  }

  return (
    <Dialog open={props.open} onOpenChange={props.onOpenChange}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Fazer um novo depósito</DialogTitle>
          <DialogDescription>
            Faça um novo depósito para a sua conta.
          </DialogDescription>
        </DialogHeader>
        <div>
          <form className="space-y-4" onSubmit={handleSendCommand}>
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

            {errors.length > 0 && (
              <div className="space-y-1">
                {errors.map((error, idx) => (
                  <p key={idx} className="text-sm text-red-500">
                    {error}
                  </p>
                ))}
              </div>
            )}

            <Button type="submit">Fazer depósito</Button>
          </form>
        </div>
      </DialogContent>
    </Dialog>
  )
}
