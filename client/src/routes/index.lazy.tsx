import { Button } from "@/components/ui/button"
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card"
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form"
import { Input } from "@/components/ui/input"
import { toast } from "@/components/ui/use-toast"
import { setupAddressFormSchema } from "@/schemas/forms"
import { useBank } from "@/stores/bank-store"
import { zodResolver } from "@hookform/resolvers/zod"
import { createLazyFileRoute, useNavigate } from "@tanstack/react-router"
import { useForm } from "react-hook-form"
import { z } from "zod"

export const Route = createLazyFileRoute("/")({
  component: Index,
})

function Index() {
  const navigate = useNavigate()
  const { setAddress } = useBank()
  const form = useForm<z.infer<typeof setupAddressFormSchema>>({
    resolver: zodResolver(setupAddressFormSchema),
    defaultValues: {
      address: "",
    },
  })

  async function onSubmit(data: z.infer<typeof setupAddressFormSchema>) {
    try {
      await fetch(`${data.address}/api`, {
        signal: AbortSignal.timeout(200),
      })
    } catch (error) {
      toast({
        title: "Erro ao acessar o banco",
        description: "O banco não respondeu.",
        variant: "destructive",
      })
      return
    }
    setAddress(data.address)
    navigate({
      to: "/login",
    })
  }

  return (
    <main className="flex flex-1 items-center justify-center">
      <Card className="mx-auto w-96 max-w-sm">
        <CardHeader>
          <CardTitle className="text-xl">Bem-vindo ao InterBank</CardTitle>
          <CardDescription>
            Digite o endereço ip do banco que deseja acessar:
          </CardDescription>
        </CardHeader>
        <CardContent>
          <Form {...form}>
            <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-6">
              <FormField
                control={form.control}
                name="address"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Endereço do Banco:</FormLabel>
                    <FormControl>
                      <Input
                        placeholder="Ex: yellowbank.example.com"
                        // type="url"
                        {...field}
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <Button type="submit">Acessar banco</Button>
            </form>
          </Form>
        </CardContent>
      </Card>
    </main>
  )
}
