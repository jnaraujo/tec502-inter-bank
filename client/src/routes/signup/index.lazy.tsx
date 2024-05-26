import { SignUpForm } from "@/components/signup-form"
import { createLazyFileRoute } from "@tanstack/react-router"

export const Route = createLazyFileRoute("/signup/")({
  component: SignUpPage,
})

function SignUpPage() {
  return (
    <main className="flex flex-1 items-center justify-center">
      <SignUpForm />
    </main>
  )
}
