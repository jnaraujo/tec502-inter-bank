import { cn } from "@/lib/utils"
import { ArrowDownLeft, ArrowRight, ArrowUpRight } from "lucide-react"

type TransactionType = "received" | "sent"
type TransactionStatus = "success" | "pending" | "failed"

interface Props {
  amount: number
  author: string
  type: TransactionType
  createdAt: string
  status: TransactionStatus
}

export function TransactionBox(props: Props) {
  function getLabel() {
    if (props.type === "received") {
      switch (props.status) {
        case "success":
          return "Transação recebida"
        case "failed":
          return "Recebimento falhou"
        default:
          return "Recebimento pendente"
      }
    } else if (props.type === "sent") {
      switch (props.status) {
        case "success":
          return "Transação enviada"
        case "failed":
          return "Envio falhou"
        default:
          return "Envio pendente"
      }
    }
  }

  return (
    <div className="flex items-center space-x-3 px-2">
      {props.type === "received" ? (
        <ArrowDownLeft
          className={cn("size-8", {
            "text-green-500": props.status === "success",
            "text-zinc-400": props.status === "pending",
            "text-orange-600": props.status === "failed",
          })}
        />
      ) : (
        <ArrowUpRight
          className={cn("size-8", {
            "text-red-500": props.status === "success",
            "text-zinc-400": props.status === "pending",
            "text-orange-600": props.status === "failed",
          })}
        />
      )}

      <div className="h-full w-full space-y-0.5">
        <span className="text-sm font-medium text-zinc-950">{getLabel()}</span>
        <h3
          className={cn("text-2xl font-medium text-red-600", {
            "text-red-600": props.type === "sent",
            "text-green-600": props.type === "received",

            "text-zinc-400": props.status === "pending",
            "text-orange-600": props.status === "failed",
          })}
        >
          R$ {props.amount}
        </h3>
        <div className="flex items-center gap-1">
          <span className="text-sm text-zinc-500">José da Silva (1-1)</span>
          <ArrowRight className="size-4 text-zinc-500" />
          <span className="text-sm text-zinc-500">Mario da Silva (2-100)</span>
        </div>
      </div>

      <div className="flex w-[200px] justify-end self-start">
        <span className="text-left text-sm text-zinc-500">
          {new Date(props.createdAt).toLocaleString("pt-br")}
        </span>
      </div>
    </div>
  )
}
