import { z } from "zod"

export const IBK_REGEX = /^\d+-\d+$/
export const createOperationFormSchema = z
  .object({
    from: z.string().regex(IBK_REGEX, {
      message: "O formado do IBK do pagador é inválido.",
    }),
    to: z.string().regex(IBK_REGEX, {
      message: "O formado do IBK do beneficiário é inválido.",
    }),
    amount: z.coerce
      .number({
        message: "O valor deve ser um número.",
      })
      .min(0.1, {
        message: "O valor deve ser maior que R$ 0.10",
      }),
  })
  .refine((data) => data.from !== data.to, {
    message: "O IBK do pagador e do beneficiário não podem ser iguais.",
  })
