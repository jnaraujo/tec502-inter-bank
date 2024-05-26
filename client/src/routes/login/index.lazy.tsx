import { LoginForm } from "@/components/login-form"
import { createLazyFileRoute } from "@tanstack/react-router"

export const Route = createLazyFileRoute("/login/")({
  component: LoginPage,
})

function LoginPage() {
  return (
    <main className="flex flex-1 items-center justify-center">
      <LoginForm />
    </main>
  )
}
