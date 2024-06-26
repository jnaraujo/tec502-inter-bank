import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog"
import { Button } from "./ui/button"
import { Input } from "./ui/input"
import { Label } from "./ui/label"

export interface Operation {
  to: string
  amount: number
}

interface Props {
  open: boolean
  onOpenChange: (open: boolean) => void
  onOperationCreated: (operation: Operation) => void
}

export function CreateTransactionDialog(props: Props) {
  function handleSendCommand(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault()

    const formData = new FormData(event.currentTarget)
    props.onOperationCreated({
      to: formData.get("to") as string,
      amount: Number(formData.get("amount")),
    })

    props.onOpenChange(false)
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
            <Button type="submit">Adicionar operação</Button>
          </form>
        </div>
      </DialogContent>
    </Dialog>
  )
}
