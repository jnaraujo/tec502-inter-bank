import { env } from "@/env"

interface Operation {
  from: string
  to: string
  amount: number
}

interface Transaction {
  author: string
  operations: Array<Operation>
}

export async function sendTransaction(transaction: Transaction) {
  const response = await fetch(`${env.VITE_BANK_URL}/api/payments/pay`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({
      author: transaction.author,
      operations: transaction.operations.map((op) => ({
        from_user_ibk: op.from,
        to_user_ibk: op.to,
        amount: op.amount,
      })),
    }),
  })

  if (!response.ok) {
    throw new Error((await response.json()).error[0].value)
  }

  const res = await response.json()
  return {
    message: res.message,
  } satisfies {
    message: string
  }
}
