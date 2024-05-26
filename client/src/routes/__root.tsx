import { Toaster } from "@/components/ui/toaster"
import { AuthProvider } from "@/contexts/auth-context"
import { createRootRoute, Outlet } from "@tanstack/react-router"
import { TanStackRouterDevtools } from "@tanstack/router-devtools"

export const Route = createRootRoute({
  component: () => (
    <>
      <div className="flex min-h-[100svh] flex-col bg-muted font-sans">
        <AuthProvider>
          <Outlet />
        </AuthProvider>
      </div>
      <TanStackRouterDevtools />
      <Toaster />
    </>
  ),
})
