import { Footer } from "@/components/footer"
import { Toaster } from "@/components/ui/toaster"
import { AuthProvider } from "@/contexts/auth-context"
import { useBank } from "@/stores/bank-store"
import {
  createRootRoute,
  Navigate,
  Outlet,
  useLocation,
} from "@tanstack/react-router"
import React from "react"

const TanStackRouterDevtools =
  process.env.NODE_ENV === "production"
    ? () => null // Render nothing in production
    : React.lazy(() =>
        // Lazy load in development
        import("@tanstack/router-devtools").then((res) => ({
          default: res.TanStackRouterDevtools,
          // For Embedded Mode
          // default: res.TanStackRouterDevtoolsPanel
        })),
      )

export const Route = createRootRoute({
  component: Root,
})

function Root() {
  const { address } = useBank()
  const location = useLocation()

  if (location.pathname != "/" && !address) {
    return <Navigate to="/" />
  }

  return (
    <>
      <div className=" flex h-[100svh] flex-col justify-between gap-6 overflow-auto bg-muted p-6 font-sans">
        <AuthProvider>
          <Outlet />
          <Footer />
        </AuthProvider>
      </div>
      <TanStackRouterDevtools />
      <Toaster />
    </>
  )
}
