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
import { toast } from "./ui/use-toast"

const formSchema = z.object({
  name: z
    .string()
    .min(8, {
      message: "O nome deve ter no mínimo 6 caracteres",
    })
    .max(255, {
      message: "O nome deve ter no máximo 255 caracteres",
    }),
  document: z.string().min(4).max(16),
})

export function SignUpForm() {
  const router = useRouter()
  const auth = useAuth()
  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      name: "",
      document: "",
    },
  })

  async function onSubmit(data: z.infer<typeof formSchema>) {
    try {
      await auth.signUp(data)
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
              name="name"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Nome</FormLabel>
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
