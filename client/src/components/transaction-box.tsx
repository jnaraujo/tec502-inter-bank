import { cn } from "@/lib/utils"
import {
  ArrowRight,
  CircleCheckBig,
  CircleDashed,
  CircleX,
  DollarSign,
} from "lucide-react"

type TransactionStatus = "success" | "pending" | "failed"

interface Operation {
  amount: number
  from: string
  to: string
  type: string
}

interface Props {
  createdAt: string
  status: TransactionStatus
  type: "final" | "package"
  operations: Array<Operation>
}

export function TransactionBox(props: Props) {
  function getLabel() {
    switch (props.status) {
      case "success":
        return "Transação realizada"
      case "failed":
        return "Transação falhou"
      default:
        return "Transação pendente"
    }
  }

  return (
    <div className="flex items-center space-x-3 px-2">
      {props.status === "success" ? (
        <CircleCheckBig className="size-8 text-green-500" />
      ) : null}
      {props.status === "failed" ? (
        <CircleX className="size-8 text-red-500" />
      ) : null}
      {props.status === "pending" ? (
        <CircleDashed className="size-8 text-zinc-400" />
      ) : null}

      <div className="h-full w-full space-y-0.5">
        <h3
          className={cn("text-sm font-medium", {
            "text-green-600": props.status === "success",
            "text-red-600": props.status === "failed",
            "text-zinc-500": props.status === "pending",
          })}
        >
          {getLabel()} -{" "}
          {props.type == "package"
            ? "(Pacote de transações)"
            : "(Transação final)"}
        </h3>

        <div>
          {props.operations.map((op, idx) => (
            <div key={idx} className="flex items-center space-x-1">
              <div className="flex items-center gap-1">
                <span className="text-sm text-zinc-500">De: {op.from}</span>
                <ArrowRight className="size-4 text-zinc-500" />
                <span className="text-sm text-zinc-500">Para: {op.to}</span>
                <div className="ml-2 flex items-center gap-1">
                  <DollarSign className="size-4 text-zinc-500" />
                  <span className="text-sm text-zinc-500">
                    R$ {op.amount.toFixed(2)}
                  </span>
                </div>
                <span className="text-sm text-zinc-400">({op.type})</span>
              </div>
            </div>
          ))}
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
