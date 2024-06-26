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
    <article className="flex h-fit flex-col justify-between space-y-2 rounded-lg border bg-background p-6">
      <h1 className="text-xl font-medium text-zinc-900">
        Bem-vindo, {user?.name}! (IBK: {user?.ibk} -{" "}
        {userTypeToHuman(user?.type)})
      </h1>
      <p className="text-zinc-500">
        Você está usando o nosso Sistema Inter Bancário (InterBank). Com ele,
        você pode realizar transações entre suas contas em diferentes bancos no
        consorcio.
      </p>
    </article>
  )
}
