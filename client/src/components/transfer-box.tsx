import { Button } from "./ui/button"
import { Input } from "./ui/input"
import { Label } from "./ui/label"

export function TransferBox() {
  return (
    <article className="space-y-4 rounded-lg border bg-background p-6">
      <h2 className="text-lg font-medium text-zinc-500">
        Realizar transferência
      </h2>

      <form className="space-y-4">
        <div className="space-y-1">
          <Label htmlFor="to">IBK do beneficiário:</Label>
          <Input placeholder="Ex: 1-203" required />
        </div>
        <div className="space-y-1">
          <Label htmlFor="to">Valor:</Label>
          <Input placeholder="Ex: 500" type="number" required />
        </div>
        <Button type="submit">Realizar transferência</Button>
      </form>
    </article>
  )
}
