import { createLazyFileRoute } from "@tanstack/react-router"

export const Route = createLazyFileRoute("/dashboard/")({
  component: DashboardPage,
})

function DashboardPage() {
  return (
    <main className="grid flex-1">
      <h1>Bem-vindo, Usuário!</h1>
    </main>
  )
}
