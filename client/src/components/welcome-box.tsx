import { useAuth } from "@/contexts/auth-context"

export function WelcomeBox() {
  const { user } = useAuth()

  function userTypeToHuman(type?: string) {
    switch (type) {
      case "individual":
        return "Pessoa Física"
      case "legal":
        return "Pessoa Jurídica"
      case "joint":
        return "Conta Conjunta"
      default:
        return "Desconhecido"
    }
  }

  return (
    <article className="grid h-fit grid-cols-2 gap-6 rounded-lg border bg-background p-6">
      <div className="space-y-2">
        <h1 className="text-xl font-medium text-zinc-900">
          Bem-vindo, {user?.name}!
        </h1>
        <p className="text-zinc-500">
          Você está usando o nosso Sistema Inter Bancário (InterBank). Com ele,
          você pode realizar transações entre suas contas em diferentes bancos
          no consorcio.
        </p>
      </div>
      <div>
        <h2 className="text-lg font-medium text-zinc-900">Dados da Conta:</h2>
        <ul className="space-y-1">
          <li className="text-sm text-zinc-500">IBK: {user?.ibk}</li>
          <li className="text-sm text-zinc-500">
            Tipo: {userTypeToHuman(user?.type)}
          </li>
          <li className="text-sm text-zinc-500">
            Documento(s): {user?.documents.join(", ")}
          </li>
        </ul>
      </div>
    </article>
  )
}
