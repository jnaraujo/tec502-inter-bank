import { useAuth } from "@/contexts/auth-context"
import { createLazyFileRoute, Navigate } from "@tanstack/react-router"

export const Route = createLazyFileRoute("/dashboard/")({
  component: DashboardPage,
})

function DashboardPage() {
  const auth = useAuth()

  if (!auth.isAuthenticated) {
    return <Navigate to="/login" />
  }

  return (
    <main className="grid flex-1">
      <h1>
        Bem-vindo, {auth.user?.name}! (Id: {auth.user?.id})
      </h1>
    </main>
  )
}
