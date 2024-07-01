import { createLazyFileRoute, Link } from "@tanstack/react-router"

export const Route = createLazyFileRoute("/")({
  component: Index,
})

function Index() {
  return (
    <div className="p-2">
      <h3>Bem-vindo ao InterBank</h3>
      <ul>
        <li>
          <Link to="/login">Login</Link>
        </li>
        <li>
          <Link to="/signup">Criar uma conta</Link>
        </li>
        <li>
          <Link to="/dashboard">Dashboard do banco</Link>
        </li>
      </ul>
    </div>
  )
}
