import path from "path"
import { TanStackRouterVite } from "@tanstack/router-vite-plugin"
import legacy from "@vitejs/plugin-legacy"
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
    legacy({
      targets: ["ie >= 11"],
      renderLegacyChunks: true,
      modernPolyfills: true,
    }),
    TanStackRouterVite(),
  ],
  resolve: {
    alias: {
      "@": path.resolve(__dirname, "./src"),
    },
  },
})
