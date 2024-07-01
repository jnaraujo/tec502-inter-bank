import { useAuth } from "@/contexts/auth-context"
import { zodResolver } from "@hookform/resolvers/zod"
import { Link, useRouter } from "@tanstack/react-router"
import { useForm } from "react-hook-form"
import { z } from "zod"
import { Button } from "./ui/button"
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "./ui/card"
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "./ui/form"
import { Input } from "./ui/input"
import { RadioGroup, RadioGroupItem } from "./ui/radio-group"
import { toast } from "./ui/use-toast"

const formSchema = z
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

export function SignUpForm() {
  const router = useRouter()
  const auth = useAuth()
  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      type: "individual",
      name: "",
      document: "",
      secondDocument: "",
    },
  })

  async function onSubmit(data: z.infer<typeof formSchema>) {
    try {
      await auth.signUp({
        name: data.name,
        type: data.type,
        documents: [data.document, data.secondDocument].filter(
          Boolean,
        ) as string[],
      })
      router.navigate({
        to: "/dashboard",
      })
    } catch (error) {
      console.log(error)
      toast({
        title: "Erro ao criar conta",
        description: (error as any).message,
        variant: "destructive",
      })
    }
  }

  return (
    <Card className="mx-auto w-96 max-w-sm">
      <CardHeader>
        <CardTitle className="text-xl">Criar conta</CardTitle>
        <CardDescription>
          Entre com suas credenciais para criar uma nova conta
        </CardDescription>
      </CardHeader>
      <CardContent>
        <Form {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-6">
            <FormField
              control={form.control}
              name="type"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Tipo de conta</FormLabel>
                  <FormControl>
                    <RadioGroup
                      onValueChange={field.onChange}
                      defaultValue={field.value}
                      className="flex justify-between"
                    >
                      <FormItem className="flex items-center space-x-3 space-y-0">
                        <FormControl>
                          <RadioGroupItem value="individual" />
                        </FormControl>
                        <FormLabel className="font-normal">
                          Individual
                        </FormLabel>
                      </FormItem>

                      <FormItem className="flex items-center space-x-3 space-y-0">
                        <FormControl>
                          <RadioGroupItem value="legal" />
                        </FormControl>
                        <FormLabel className="font-normal">Legal</FormLabel>
                      </FormItem>

                      <FormItem className="flex items-center space-x-3 space-y-0">
                        <FormControl>
                          <RadioGroupItem value="joint" />
                        </FormControl>
                        <FormLabel className="font-normal">Conjunta</FormLabel>
                      </FormItem>
                    </RadioGroup>
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="name"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>
                    {form.getValues("type") === "joint"
                      ? "Nome da conta conjunta"
                      : "Nome do titular"}
                  </FormLabel>
                  <FormControl>
                    <Input placeholder="Ex: John Doe" {...field} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
            <FormField
              control={form.control}
              name="document"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Documento</FormLabel>
                  <FormControl>
                    <Input placeholder="Ex: 123.456.789-99" {...field} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
            {form.getValues("type") === "joint" && (
              <FormField
                control={form.control}
                name="secondDocument"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Documento do segundo titular</FormLabel>
                    <FormControl>
                      <Input placeholder="Ex: 123.456.789-99" {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
            )}
            <Button type="submit">Cadastrar</Button>
          </form>
        </Form>
      </CardContent>
      <CardFooter className="space-x-1 text-sm">
        <span>Já tem conta?</span>
        <Link className="text-purple-700" to="/login">
          Entre agora!
        </Link>
      </CardFooter>
    </Card>
  )
}
