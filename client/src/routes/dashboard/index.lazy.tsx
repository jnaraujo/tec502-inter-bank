import { TransactionBox } from "@/components/transaction-box"
import { useAuth } from "@/contexts/auth-context"
import { createLazyFileRoute, Navigate } from "@tanstack/react-router"

export const Route = createLazyFileRoute("/dashboard/")({
  component: DashboardPage,
})

function DashboardPage() {
  const { isAuthenticated, user } = useAuth()

  if (!isAuthenticated || !user) {
    return <Navigate to="/login" />
  }

  return (
    <main className="container grid h-screen grid-cols-[1fr_400px] gap-6 py-6">
      <section className="grid grid-rows-[auto_1fr] gap-x-6 gap-y-4 overflow-auto">
        <article className="flex h-fit flex-col justify-between space-y-2 rounded-lg border bg-background p-6">
          <h1 className="text-xl font-medium text-zinc-900">
            Bem-vindo, {user.name}! (IBK: {user.ibk})
          </h1>
          <p>
            Você está usando o nosso Sistema Inter Bancário (InterBank). Com
            ele, você pode realizar transações entre suas contas em diferentes
            bancos no consorcio.
          </p>
        </article>

        <article className="flex h-fit flex-col justify-between space-y-4 rounded-lg border bg-background p-6">
          <h2 className="text-xl font-medium text-zinc-900">Suas transações</h2>

          <TransactionBox
            amount={200}
            author="José da Silva"
            createdAt="2024-06-25T16:35:57.794293059-03:00"
            type="sent"
            status="pending"
          />

          <TransactionBox
            amount={200}
            author="José da Silva"
            createdAt="2024-06-25T16:35:57.794293059-03:00"
            type="received"
            status="success"
          />

          <TransactionBox
            amount={200}
            author="José da Silva"
            createdAt="2024-06-25T16:35:57.794293059-03:00"
            type="received"
            status="failed"
          />
        </article>
      </section>

      <section className="rounded-lg border bg-background p-6">
        <h2>oi</h2>
      </section>

      {/* <article className="flex flex-col justify-center gap-2 rounded-lg border bg-background p-6">
        <h2>Suas transaçoes</h2>
      </article>

      <section className="grid grid-rows-[350px_1fr] gap-4 overflow-auto"></section> */}
    </main>
  )
}
