import path from "path"
import { TanStackRouterVite } from "@tanstack/router-vite-plugin"
import react from "@vitejs/plugin-react-swc"
import Unfonts from "unplugin-fonts/vite"
import { defineConfig } from "vite"

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [
    react(),
    Unfonts({
      fontsource: {
        families: [
          {
            name: "Inter",
            weights: [400, 500, 600, 700],
            styles: ["normal"],
            subset: "latin",
          },
        ],
      },
    }),
    TanStackRouterVite(),
  ],
  resolve: {
    alias: {
      "@": path.resolve(__dirname, "./src"),
    },
  },
})
