import { Toaster } from "@/components/ui/toaster"
import { AuthProvider } from "@/contexts/auth-context"
import { createRootRoute, Outlet } from "@tanstack/react-router"
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
