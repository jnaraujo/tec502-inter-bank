import { useAuth } from "@/contexts/auth-context"
import { useBank } from "@/stores/bank-store"

export function Footer() {
  const { address } = useBank()
  const { user } = useAuth()
  return (
    <footer className="container mx-auto mb-2 rounded-lg border bg-background p-3 text-center">
      <p className="text-zinc-700">InterBank @ {new Date().getFullYear()}</p>
      <p className="text-sm text-zinc-500">
        Bank Address: {address || "address not set"} | Logged in as:{" "}
        {user ? user.name : "not logged in"}
      </p>
    </footer>
  )
}
