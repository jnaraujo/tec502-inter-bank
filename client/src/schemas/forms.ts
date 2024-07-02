import { z } from "zod"
import { IBK_REGEX } from "./operation"

export const loginFormSchema = z.object({
  ibk: z.string().regex(IBK_REGEX, {
    message: "O formado do IBK do pagador é inválido.",
  }),
})

export const setupAddressFormSchema = z.object({
  address: z.string().url({
    message: "O endereço do banco precisa ser um URL.",
  }),
})

export const signUpFormSchema = z
  .object({
    type: z.enum(["individual", "legal", "joint"], {
      required_error: "Selecione um tipo de conta",
    }),
    name: z
      .string()
      .min(8, {
        message: "O nome deve ter no mínimo 6 caracteres",
      })
      .max(255, {
        message: "O nome deve ter no máximo 255 caracteres",
      }),
    document: z
      .string()
      .min(4, {
        message: "O documento deve ter no mínimo 4 caracteres",
      })
      .max(16, {
        message: "O documento deve ter no máximo 16 caracteres",
      }),
    secondDocument: z
      .string()
      .min(4, {
        message: "O documento deve ter no mínimo 4 caracteres",
      })
      .max(16, {
        message: "O documento deve ter no máximo 16 caracteres",
      })
      .optional()
      .or(z.literal("")),
  })
  .refine(
    (data) => {
      if (data.type !== "joint") return true
      return data.secondDocument !== ""
    },
    {
      message: "O documento do segundo titular é obrigatório",
      path: ["secondDocument"],
    },
  )
