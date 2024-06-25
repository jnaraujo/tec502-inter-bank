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
  document: z.string().min(4, {
    message: "O documento deve ter no mínimo 4 caracteres",
  }),
})

export function LoginForm() {
  const auth = useAuth()
  const router = useRouter()
  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      document: "",
    },
  })

  async function onSubmit(data: z.infer<typeof formSchema>) {
    try {
      await auth.login(data)
      router.navigate({
        to: "/dashboard",
      })
    } catch (error) {
      toast({
        title: "Erro ao entrar na conta",
        description: (error as any).message,
        variant: "destructive",
      })
    }
  }

  return (
    <Card className="mx-auto w-96 max-w-sm">
      <CardHeader>
        <CardTitle className="text-xl">Login</CardTitle>
        <CardDescription>
          Entre com suas credenciais para entrar na sua conta
        </CardDescription>
      </CardHeader>
      <CardContent>
        <Form {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-6">
            <FormField
              control={form.control}
              name="document"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Seu Documento:</FormLabel>
                  <FormControl>
                    <Input placeholder="Ex: 123.456.789-99" {...field} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
            <Button type="submit">Entrar</Button>
          </form>
        </Form>
      </CardContent>
      <CardFooter className="space-x-1 text-sm">
        <span>Não tem conta?</span>
        <Link className="text-purple-700" to="/signup">
          Crie uma agora!
        </Link>
      </CardFooter>
    </Card>
  )
}
