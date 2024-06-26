import { ZodError } from "zod"

export function handleZodError({ issues }: ZodError<unknown>): string[] {
  return issues.map((issue) => issue.message)
}
