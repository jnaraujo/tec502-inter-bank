import { BalanceBox } from "@/components/balance-box"
import { TransactionListBox } from "@/components/transaction-list-box"
import { TransferBox } from "@/components/transfer-box"
import { WelcomeBox } from "@/components/welcome-box"
import { useAuth } from "@/contexts/auth-context"
import { createLazyFileRoute, Navigate } from "@tanstack/react-router"

export const Route = createLazyFileRoute("/dashboard/")({
  component: DashboardPage,
})

function DashboardPage() {
  const { isAuthenticated, user } = useAuth()

  if (user === undefined) return
  if (!isAuthenticated) {
    return <Navigate to="/login" />
  }

  return (
    <main className="container grid flex-1 grid-cols-[1fr_400px] gap-6 overflow-auto p-0">
      <section className="grid grid-rows-[auto_1fr] gap-x-6 gap-y-4 overflow-auto">
        <WelcomeBox />
        <TransactionListBox />
      </section>

      <section className="grid grid-rows-[auto_1fr] space-y-4 overflow-auto">
        <BalanceBox />
        <TransferBox />
      </section>
    </main>
  )
}
