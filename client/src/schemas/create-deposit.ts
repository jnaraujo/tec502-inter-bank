import { z } from "zod"

export const createDepositFormSchema = z.object({
  amount: z.coerce
    .number({
      message: "O valor deve ser um número.",
    })
    .min(0.1, {
      message: "O valor deve ser maior que R$ 0.10",
    }),
})
