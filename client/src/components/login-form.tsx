import { useAuth } from "@/contexts/auth-context"
import { loginFormSchema } from "@/schemas/forms"
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

export function LoginForm() {
  const auth = useAuth()
  const router = useRouter()
  const form = useForm<z.infer<typeof loginFormSchema>>({
    resolver: zodResolver(loginFormSchema),
    defaultValues: {
      ibk: "",
    },
  })

  async function onSubmit(data: z.infer<typeof loginFormSchema>) {
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
              name="ibk"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Seu IBK:</FormLabel>
                  <FormControl>
                    <Input placeholder="Ex: 1-23" {...field} />
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
