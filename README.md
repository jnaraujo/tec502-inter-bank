<div align="center">
  <h1>InterBank - Sistema de pagamentos bancários descentralizado</h1>
  <p>
    <strong>Projeto desenvolvido para a disciplina TEC502 - MI Concorrência e Conectividade</strong>
  </p>

  ![Most used language](https://img.shields.io/github/languages/top/jnaraujo/tec502-inter-bank?style=flat-square)
  ![GitHub](https://img.shields.io/github/license/jnaraujo/tec502-inter-bank)
</div>

<div align="center">
  <img src="/images/interface.png" alt="Client web" height="400px" width="auto" />
</div>

<br />

Nos últimos anos, o sistema bancário brasileiro tem passado por grandes revoluções. Com a criação do meio de pagamento Pix, milhões de brasileiros passaram a realizar transferência de forma simples, instantânea e sem taxas. Segundo dados do [Banco Central](https://www.bcb.gov.br/detalhenoticia/803/noticia), o Pix foi responsável por mais de 40 bilhões de transações realizadas, totalizando mais de R$ 17 trilhões movimentados.

Como forma de criar um sucessor para o Pix descentralizado e integrado para pagamentos bancários para um país sem banco central, foi desenvolvido o InterBank. O objetivo do InterBank é promover uma integração entre diferentes bancos participantes do consórcio, permitindo aos clientes realizarem transferências em suas contas a partir de qualquer banco. Além disso, cada transação passa a ser vista como um pacote, permitindo ao cliente realizar diferentes operações na mesma operação, de maneira atômica e consistente.

Como forma de solucionar esse problema, foram utilizadas tecnologias como Docker, ReactJS e Go. Além disso, para resolução do problema da concorrência e consistência das transações, o método Token Ring e uma variação do protocolo Two-Phase Commit foram utilizados. Desse modo, foi possível implementar uma solução que permita transações interbancárias ordenadas, sem conflito e consistente.

## Sumário
- [Sumário](#sumário)
- [Sobre o projeto](#sobre-o-projeto)
  - [Tecnologias utilizadas](#tecnologias-utilizadas)
  - [Como executar o projeto?](#como-executar-o-projeto)
- [Como utilizar?](#como-utilizar)
    - [Como definir o Banco?](#como-definir-o-banco)
    - [Como criar de uma conta](#como-criar-de-uma-conta)
    - [Como realizar um depósito](#como-realizar-um-depósito)
    - [Como visualizar o saldo de uma conta](#como-visualizar-o-saldo-de-uma-conta)
    - [Como realizar uma transferência](#como-realizar-uma-transferência)
    - [Como adicionar uma nova operação](#como-adicionar-uma-nova-operação)
    - [Como visualizar o extrato de uma conta](#como-visualizar-o-extrato-de-uma-conta)
- [Arquitetura do projeto](#arquitetura-do-projeto)
  - [Interface gráfica](#interface-gráfica)
  - [Banco e Interbank](#banco-e-interbank)
- [Transações interbancárias](#transações-interbancárias)
  - [Atomicidade](#atomicidade)
  - [Assincronia](#assincronia)
  - [Consistência](#consistência)
- [Token Ring](#token-ring)
  - [Concorrência distribuída](#concorrência-distribuída)
  - [Transações simultâneas](#transações-simultâneas)
  - [Estrutura do Token Ring](#estrutura-do-token-ring)
  - [Inicialização do Token Ring](#inicialização-do-token-ring)
  - [Passagem do Token](#passagem-do-token)
  - [Detecção e recuperação de falhas](#detecção-e-recuperação-de-falhas)
    - [Queda de algum dos bancos](#queda-de-algum-dos-bancos)
    - [Perda de token](#perda-de-token)
    - [Duplicação de token](#duplicação-de-token)
- [Comunicação](#comunicação)
  - [Comunicação entre a interface e o banco](#comunicação-entre-a-interface-e-o-banco)
    - [Rotas da API](#rotas-da-api)
      - [POST /api/accounts](#post-apiaccounts)
      - [POST /api/accounts/auth](#post-apiaccountsauth)
      - [DELETE /api/accounts/:id](#delete-apiaccountsid)
      - [POST /api/payments/deposit](#post-apipaymentsdeposit)
      - [GET /api/accounts/:id](#get-apiaccountsid)
      - [GET /api/accounts/:id/all](#get-apiaccountsidall)
      - [GET /api/accounts/:id/transactions](#get-apiaccountsidtransactions)
      - [POST /api/payments/pay](#post-apipaymentspay)
  - [Comunicação entre os bancos (InterBank)](#comunicação-entre-os-bancos-interbank)
    - [Rotas do InterBank](#rotas-do-interbank)
      - [POST /interbank/prepare](#post-interbankprepare)
      - [POST /interbank/rollback](#post-interbankrollback)
      - [POST /interbank/commit](#post-interbankcommit)
      - [GET /account/:document](#get-accountdocument)
      - [POST /interbank/account/ibk/:ibk](#post-interbankaccountibkibk)
      - [PUT /interbank/token](#put-interbanktoken)
      - [GET /interbank/token](#get-interbanktoken)
      - [GET /interbank/token/ok](#get-interbanktokenok)
- [Sincronização dos dados internamente no banco](#sincronização-dos-dados-internamente-no-banco)
- [Testes](#testes)
- [Conclusão](#conclusão)

## Sobre o projeto
### Tecnologias utilizadas
- Geral
  - [Docker](https://www.docker.com/): Plataforma de código aberto para criação, execução e gerenciamento de aplicações em containers.
  - [Docker Compose](https://docs.docker.com/compose/): Ferramenta para definir e executar aplicações Docker em múltiplos containers.
- Cliente
  - [React](https://reactjs.org/): Biblioteca JavaScript para a criação de interfaces de usuário.
  - [Vite](https://vitejs.dev/): Ferramenta de build para aplicações web.
  - [TypeScript](https://www.typescriptlang.org/): Superset da linguagem JavaScript que adiciona tipagem estática ao código.
  - [TanStack Query](https://tanstack.com/query/latest): Biblioteca para gerenciamento de estado e requisições HTTP. Responsável por fazer a comunicação com o Broker.
- Bancos
  - [Go](https://golang.org/): Linguagem de programação utilizada para o desenvolvimento do Broker.
  - [Fiber](https://gofiber.io/): Framework web para Go.
- Testes de integração
  - [Python](https://www.python.org/): Linguagem de programação utilizada para o teste de integração

### Como executar o projeto?
1. Clone o repositório localmente:
```bash
git clone https://github.com/jnaraujo/tec502-inter-bank
```

2. Acesse o diretório do projeto:
```bash
cd tec502-inter-bank
```

3. Execute o comando abaixo para iniciar o projeto:
```bash
docker-compose up --build
```

4. Acesse o endereço `http://localhost:3000` no seu navegador para acessar a interface do cliente.

## Como utilizar?
#### Como definir o Banco?
<div align="center">
<img src="./images/set-address.png" alt="Definir endereço do banco" height="300px" width="auto" /> <br/>
<em>Figura 1. Definir endereço do banco</em>
</div>

Para definir o endereço do banco, escreva o endereço do banco no campo de texto e clique no botão "Acessar banco". O endereço do banco é utilizado para realizar a comunicação entre os bancos do consórcio. Após definir o endereço do banco, você será redirecionado para a página de login.

#### Como criar de uma conta
<div align="center">
<img src="./images/create-account.png" alt="Criar uma nova conta" height="300px" width="auto" /> <br/>
<em>Figura 2. Criar uma nova conta</em>
</div>

Para criar uma nova conta, preencha o formulário com o nome do usuário, os documentos (CPF ou CNPJ) e o tipo da conta (pessoa física, jurídica ou conjunta) e clique no botão "Criar conta". Após criar a conta, você será redirecionado para a página do banco.

<div align="center">
<img src="./images/create-joint-account.png" alt="Criar uma nova conta conjunta" height="300px" width="auto" /> <br/>
<em>Figura 3. Criar uma nova conta conjunta</em>
</div>

Para criar uma conta conjunta, preencha o formulário com o nome da conta e os documentos dos usuários (CPF ou CNPJ) e clique no botão "Criar conta" (ambos os usuários devem ter uma conta individual previamente criada no banco). Após criar a conta, você será redirecionado para a página do banco.

#### Como realizar um depósito
<div align="center">
<img src="./images/create-deposit.png" alt="Criar um novo depósito" height="300px" width="auto" /> <br/>
<em>Figura 4. Criar um novo depósito</em>
</div>

Para realizar um depósito, preencha o formulário o valor a ser depositado e clique no botão "Fazer depósito". Após realizar o depósito, você verá o novo saldo da conta.

#### Como visualizar o saldo de uma conta
<div align="center">
<img src="./images/view-balance.png" alt="Visualizar saldo da conta" height="300px" width="auto" /> <br/>
<em>Figura 5. Visualizar saldo da conta</em>
</div>

Para visualizar o saldo de uma conta, basta verificar na caixa "Seu saldo" o valor atual da conta. O saldo é atualizado de tempos em tempos, garantindo que as informações estejam sempre atualizadas.

#### Como realizar uma transferência
<div align="center">
<img src="./images/create-transaction.png" alt="Criar transação" height="300px" width="auto" /> <br/>
<em>Figura 6. Criar transação</em>
</div>

Para realizar uma transferência, adicione uma nova operação clicando no botão "Adicionar nova operação". Em seguida, preencha o formulário. Após adicionar todas as operações, clique no botão "Realizar transferência".

#### Como adicionar uma nova operação
<div align="center">
<img src="./images/add-new-operation.png" alt="Adicionar nova operação" height="300px" width="auto" /> <br/>
<em>Figura 7. Adicionar nova operação</em>
</div>

Para criar uma nova operação, preencha o formulário com o IBK do pagador, o IBK do beneficiário e o valor a ser transferido e clique no botão "Adicionar operação".

#### Como visualizar o extrato de uma conta
<div align="center">
<img src="./images/transaction-list.png" alt="Visualizar extrato da conta" height="300px" width="auto" /> <br/>
<em>Figura 8. Visualizar extrato da conta</em>
</div>

Para visualizar o extrato de uma conta, basta verificar a lista de transações. O extrato é atualizado de tempos em tempos, garantindo que as informações estejam sempre atualizadas.

## Arquitetura do projeto
O sistema foi divido em duas partes principais: a interface gráfica do banco e o código do banco. A interface gráfica foi desenvolvida utilizando ReactJS e a biblioteca TanStack Query para gerenciamento de estado e requisições HTTP. Já o código do banco foi desenvolvido em Go, utilizando o framework Fiber.

### Interface gráfica
Para o desenvolvimento da interface gráfica, foi utilizado a biblioteca ReactJS, uma biblioteca de código aberto para interfaces intuitivas. Além disso, a biblioteca TanStack Query foi utilizado para gerenciamento de estados e requisições HTTP, permitindo atualizar as informações em tempo real (como o saldo e a lista de transações).

A interface do usuário é formada por 4 páginas principais: a página de seleção do banco do usuário, a página de login, a página de registro e a página do banco, onde é possível realizar operações como depósito, transferência, visualização de extrato e visualização de saldo.

```bash
client
├── src # Código fonte da aplicação
│   ├── @types # Alguns tipos que são utilizados em várias partes da aplicação 
│   ├── components # Componentes React reutilizáveis
│   ├── constants # Definem valores constantes, como o tempo de atualização do saldo
│   ├── contexts # Contextos são utilizados para compartilhar informações entre componentes de forma global
|   ├── stores # Stores são utilizadas para gerenciar o estado da aplicação
│   ├── hooks # Hooks são funções que permitem adicionar funcionalidades ao componente
│   ├── lib # Funções utilitárias
│   ├── routes # Definição das rotas da aplicação
│   ├── schemas # Schemas utilizados para validação de dados
│   ├── services # Funções que realizam requisições HTTP
```

### Banco e Interbank
O desenvolvimento do banco e do InterBank se derem de maneira conjunta, visto que ambos devem funcionar de maneira próxima. Ambos foram desenvolvidos em Go, uma linguagem de programação compilada criada pelo google, utilizando o Fiber, um framework web escrito em Go e com foco em simplicidade e velocidade.

Para as funções internas do banco, foram implementadas [rotas](bank/internal/routes/bank) para realizar operações de criação de conta, login, depósito, transferência, visualização de extrato e visualização de saldo. Além disso, para o sistema InterBank, foram implementados sistema como o [Token Ring](bank/internal/services/token_ring.go), o [processador de transações em segundo plano](bank/internal/transaction_processor/processor.go) e as [rotas](bank/internal/routes/interbank) que permitem que as transações ocorram.

```bash
bank
├── bank-api # Arquivos para teste da api
├── cmd # Comandos para execução da api e do token ring
├── internal
│   ├── config # Configurações do banco
│   ├── constants # Definições de constantes
│   ├── http # Configuração do servidor HTTP
│   ├── interbank # Definição dos padrões do InterBank
│   ├── models # Definição dos modelos de dados
│   ├── routes # Definição das rotas da aplicação
│   ├── services # Serviços utilizados para realizar as operações
│   ├── storage # Armazenamento dos dados
│   ├── token # Definição da estrutura de um Token
│   ├── transaction_processor # Serviço que roda em segundo plano para processar as transações
|   ├── validate # Funções para validação de dados
```

## Transações interbancárias
O principal objetivo do InterBank é promover uma integração entre os bancos do consorcio, permitindo que transações possam ser realizadas entre as contas de diferentes bancos. Desse modo, é importante criar um sistema seguro e eficiente, que permita ao usuário realizar transações atômicas, consistentes e livre de errors.

Para tal, toda transação criada no InterBank é única, contendo campos de ID, ID da transação pai (no caso de ser uma transação final), dono da transação, tipo da transação, operações que serão realizadas, data de criação, data de atualização e o status (pendente, sucesso ou falha). No que diz respeito ao tipo da transação, elas podem ser do tipo `pacote` ou do tipo `final`. Transações do tipo `pacote` são transações que podem conter várias operações, sendo o tipo definido quando uma transação é criada na interface do usuário. Por outro lado, transações do tipo `final` é a transação propriamente dita, ou seja, a transação que realmente terá efeito na conta.

Por exemplo, quando um usuário cria uma transação no banco A que envia dinheiro de uma conta de sua propriedade no banco B para uma conta de terceiros no banco C, a transação criada no banco A serão do tipo `pacote`, enquanto as transações no banco B e no banco C serão do tipo `final`. Essa separação é importante pois torna mais simples o processo de confirmação e reversão de transações, além de separar o pacote de transações da transação que realmente irá ter algum efeito na conta.

```go
// Código de bank/internal/models/transaction.go
type Transaction struct {
	Id         TransactionId     `json:"id"`
  ParentId   *TransactionId    `json:"parent_id"`
	Owner      interbank.IBK     `json:"owner"`
	Type       TransactionType   `json:"type"`
	Operations []Operation       `json:"operations"`
	CreatedAt  time.Time         `json:"created_at"`
	UpdatedAt  time.Time         `json:"updated_at"`
	Status     TransactionStatus `json:"status"`
}
```

Vale destacar que cada operação pertencente a uma transação possui estrutura própria, possuindo campos de ID único, conta de origem, conta de destino, tipo da operação (depósito ou transferência), valor, status (pendente, sucesso ou falha), data de criação e data de atualização. O campo de status se repete pois como é necessário a confirmação de todos os bancos envolvidos em determinada transação, o momento em que cada operação é realizada varia.

```go
// Código de bank/internal/models/operation.go
type Operation struct {
	Id        uuid.UUID       `json:"id"`
	From      interbank.IBK   `json:"from"`
	To        interbank.IBK   `json:"to"`
	Type      OperationType   `json:"type"`
	Amount    decimal.Decimal `json:"amount"`
	Status    OperationStatus `json:"status"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}
```

Todas as transações realizadas entre os bancos do consórcio tem como objetivo a atomicidade, assincronia e consistência. Isso significa que as transações são realizadas de forma completa e consistente, sem que ocorram falhas ou interrupções.

### Atomicidade
Atomicidade é uma das propriedades ACID (Atomicidade, Consistência, Isolamento e Durabilidade) que garante que as transação serão realizadas de forma completa e irredutível, ou seja, a transação só é executada se todas as operações forem executadas completamente, não sendo aceitos estados intermediários. Desse modo, caso uma operação não possa ser executada, nenhuma operação será executada. Além disso, caso alguma operação já confirmada venha a falha, todas são revertidas.

<div align="center">
<img src="./images/2pc-ok.png" alt="Operação de duas fases" height="300px" width="auto" /> <br/>
<em>Figura 9. Operação de duas fases</em>
</div>

Para garantir a atomicidade no Interbank, foi utilizado uma variação do protocolo [Two-Phase Commit](https://martinfowler.com/articles/patterns-of-distributed-systems/two-phase-commit.html). Nesse protocolo, as transação realizadas em duas etapas: preparação (prepare) e confirmação (commit). Na etapa de preparação, as operações são enviadas para cada banco envolvido e é aguardado a confirmação se a operação pode ser realizada. Na segundo etapa (confirmação), as operações são de fato realizadas, caso todas todas as operações na etapa anterior tenham sido confirmadas.

<div align="center">
<img src="./images/2pc-error.png" alt="Operação de duas fases com falha" height="300px" width="auto" /> <br/>
<em>Figura 10. Operação de duas fases com falha</em>
</div>

No caso de alguma das preparações terem falhado, todas as operações são desfeitas (rollback). Além disso, caso algum erro ocorra na etapa de confirmação, as operações também são desfeitas. Em ambos os casos de falha, o status é atualizado tanto nas operações quanto da própria transação.

No código abaixo, a função `ProcessTransaction` é responsável por processar uma transação. Nela, a função `Prepare` é utilizada para preparar as operações de débito e crédito. Caso ocorra algum erro durante a preparação ou confirmação das operações, a transação é marcada como falha e as operações são revertidas. A função `Rollback` é utilizada para reverter as operações e a função `Commit` é utilizada para confirmar as operações. Além disso, as operações são marcadas como sucesso ou falha, garantindo que a transação seja realizada de forma completa e consistente.

```go
// Código de bank/internal/services/inter_bank.go
func ProcessTransaction(tr models.Transaction) error {
	externalTransactions := []txProcess{} // transações externas que foram realizadas
	isSuccess := true
	for _, op := range tr.Operations {
		txDebit := Prepare(op, StepDebit) // prepara a operação de débito
		if txDebit == nil { // se ocorrer algum erro, a transação é marcada como falha
			isSuccess = false
			break
		}
		externalTransactions = append(externalTransactions, txProcess{Tx: txDebit, Step: StepDebit})

		txCredit := Prepare(op, StepCredit) // prepara a operação de crédito
		if txCredit == nil { // se ocorrer algum erro, a transação é marcada como falha
			isSuccess = false
			break
		}
		externalTransactions = append(externalTransactions, txProcess{Tx: txCredit, Step: StepCredit})
	}

	if !isSuccess { // se ocorreu algum erro, as transações ja feitas devem sofrer rollback
		for _, tx := range externalTransactions { // as transações são revertidas
			Rollback(tx.Tx.Id, tx.Tx.Operations[0], tx.Step)
		}
		for _, op := range tr.Operations { // as operações da transação são marcadas como falha
			storage.Transactions.UpdateOperationStatus(tr, op, models.OperationStatusFailed)
		}
		storage.Transactions.UpdateTransactionStatus(tr, models.TransactionStatusFailed) // a transação é marcada como falha
		return errors.New("transaction failed")
	}

	for _, tx := range externalTransactions { // se todas as operações foram preparadas
		ok := Commit(tx.Tx.Id, tx.Tx.Operations[0], tx.Step) // as operações são confirmadas
		if !ok {
			isSuccess = false // se ocorrer algum erro, a transação é marcada como falha
			break
		}
	}

	if !isSuccess { // se ocorreu algum erro, as transações ja feitas devem sofrer rollback
	 // ...
	}

	for _, op := range tr.Operations { // as operações da transação são marcadas como sucesso
		storage.Transactions.UpdateOperationStatus(tr, op, models.OperationStatusSuccess)
	}
	storage.Transactions.UpdateTransactionStatus(tr, models.TransactionStatusSuccess) // a transação é marcada como sucesso
	return nil
}
```

### Assincronia
Todas as transações realizadas através do InterBank são realizadas de maneira assíncrona. Ou seja, todas as operações criadas são processadas em segundo plano, sem que o usuário precise esperar que a transação seja concluída. Desse modo, enquanto uma operação está sendo processada, o usuário pode realizar outras operações, como visualizar o status da transação, visualizar o saldo ou criar novas operações.

Vale destacar que quando uma transação é criada em um banco, ela é adicionada na [fila interna](bank/internal/storage/transaction_queue.go) do banco, garantindo que toda transação criada no mesmo banco seja executada em ordem. Essa fila é processada em background por um [serviço](bank/internal/transaction_processor/processor.go) que é responsável por processar as transações de forma assíncrona. O serviço verifica de tempos em tempos se o banco tem a posse do [token](#token-ring) e, caso tenha, processa a transação. Caso contrário, a transação é mantida na fila até que o banco possua o token. Desse modo, nenhuma transação é realizada até que o banco possua o token. Isso garante que as transações sejam realizadas de forma ordenada e sem conflitos.

Assim que o token é adquirido, o banco começa o processamento das transações, uma de cada vez e na ordem em que foram adicionadas na fila. Como cada transações possui N operações, o banco processa cada operação de forma atômica. Assim, caso uma operação falhe, a transação é marcada como falha e nenhuma operação é realizada.

Todas as operações realizada no processamento das operações (como verificar se o usuário existe, adicionar fundos na conta, subtrair fundos da conta, etc) são realizados através do InterBank.

### Consistência
Consistência é a garantia de que nenhuma operação realizada no sistema deixará os dados inconsistentes. Ou seja, nenhuma das transações realizadas irá atualizar os dados pela metade ou alterar os dados de maneira irregular. No InterBank, a consistência do sistema é garantida através do uso do [Token Ring](#token-ring), um protocolo de acesso ao meio que garante que somente um banco irá realizar suas transação de cada vez.

Desse modo, cada transação no sistema ocorre de maneira única, garantindo que nenhuma outra transação irá interferir na transação atual. Além disso, devido a natureza [atômica](#atomicidade) das transações, é garantido que nenhuma transação será feita pela metade.

Ademais, como todas as transações são processadas de maneira [ordenada e sem conflito](#assincronia). caso duas transações sejam criadas no mesmo banco, a primeira transação é processada antes da segunda transação. Além disso, todas as operações presentes em uma transação são realizadas de forma completa e correta, garantindo que as transações sejam realizadas de forma consistente. O InterBank, porém, não garante que a ordem das transações seja a mesma em todos os bancos, sendo possível uma transação criada no tempo P+1 no banco 2 ser processada antes de uma transação criada no tempo P no banco 1. Ainda assim, a consistência é garantida, pois todas as operações serão realizadas de forma atômica e consistente.

## Token Ring
O Token Ring é um protocolo de acesso ao meio definido pelo padrão IEEE 802.5 e baseado em topologia em anel, amplamente utilizado em redes e computadores. O Token Ring utiliza um `token` para garantir a ordem no sistema, sendo este passado de nó em nó, permitindo que todos tenham a oportunidade de acessar a rede de forma ordenada e sem conflitos.

<div align="center">
<img src="./images/token-ring.png" alt="Token Ring" height="300px" width="auto" /> <br/>
<em>Figura 9. Token Ring</em>
</div>

No contexto do InterBank, o Token Ring é utilizado para garantir que cada banco no consórcio tenha a oportunidade de acessar e atualizar as informações das contas de forma ordenada e sem conflitos. O token é passado de banco em banco, seguindo a ordem dos IDs dos bancos. Quando um banco possui o token, ele pode realizar operações de leitura e escrita nos dados armazenados. Caso um banco deseje realizar uma operação e não possua o token, ele deve esperar até que o token seja passado para ele.

### Concorrência distribuída
O uso do Token Ring garante que apenas um banco terá acesso a rede por vez, impedindo que ocorram conflitos entre os bancos. Assim, mesmo com diferentes transações sendo criadas na rede ao mesmo tempo, apenas um banco poderá processar as suas transações locais por vez. Além disso, como cada banco possui sua fila local e processa apenas uma transação por vez, é garantido que as transações sejam realizadas de forma ordenada e sem conflitos.

Desse modo, as operações são realizada uma de cada vez e na ordem em que foram adicionadas na fila, garantindo que o saldo final da conta seja consistente e sem duplicação de dados. Além disso, devido a natureza atômica das transações, até as transações que falharam são processadas de forma consistente.

### Transações simultâneas
Garantir que diferentes usuários possam realizar transações simultâneas é um dos principais desafios de um sistema distribuído. Com o método de Token Ring, é garantido que as transações serão realizadas de forma ordenada e sem conflitos, mesmo que diferentes usuários estejam realizando transações ao mesmo tempo.

Desse modo, embora N transações possam ser criadas em M banco simultaneamente, somente o banco com a posse do token pode executar suas transações. Além disso, todas as transações são executadas em ordem pelo banco, de modo que somente uma transação é executada por vez.

Assim, mesmo em caso de transações simultâneas ocorrendo no InterBank, nenhuma delas é executada de maneira concorrente. Nesse sentido, transações que afetem o mesmo usuário, tanto no mesmo banco quanto em outros bancos, serão processadas sem conflitos ou duplicação de dados.

### Estrutura do Token Ring
O [Token Ring](bank/internal/storage/ring.go) é composto por um conjunto de bancos (nós) que se comunicam entre si para realizar transações financeiras de forma segura e eficiente. Com o algoritmo de Token Ring, é garantido que todos os bancos terão acesso ao token.

Todos os bancos do consórcio são definidos com antecedência e cada banco possui um ID único. O ID é utilizado para determinar a ordem em que os bancos acessam e atualizam as informações das contas. O token é passado de banco em banco, seguindo a ordem dos IDs dos bancos.

No InterBank, o Token Ring também guarda informação do endereço IP do banco, permitindo que futuras consultas e operações possam ser realizadas.

```go
// Código de bank/internal/storage/ring.go
type ringData struct {
   Id   interbank.BankId
   Addr string
}

// implementação de um token ring para
// comunicação entre os bancos
type ringStorage struct {
   mu   sync.RWMutex
   ring []ringData
}
```

### Inicialização do Token Ring
Quando o sistema é iniciado, o banco com ID mais baixo é o responsável por criar o token e passá-lo para o próximo banco. O token é passado de banco em banco, seguindo a ordem dos IDs dos bancos. Quando o token chega no último banco, ele é passado de volta para o primeiro banco, fechando o anel. Esse processo é repetido indefinidamente, garantindo que cada banco tenha a oportunidade executar suas transações. O código a seguir demonstra como o token é passado de banco em banco.

Caso o banco de ID mais baixo não esteja online no momento em que o sistema iniciar, o banco seguinte na ordem é responsável por inicializar o Token Ring.

```go
// código de bank/internal/services/token_ring.go
// Se o banco atual é o banco com menor ID
if storage.Ring.FindBankWithLowestId().Id == config.Env.BankId {
// verifica se o token já esta na rede.
if !services.IsTokenOnRing() {
   // se não estiver, cria o token
   services.BroadcastToken(config.Env.BankId)
}
}
```

### Passagem do Token
Quando um banco possui o token, ele pode realizar suas operações nos bancos do sistema. Caso um banco deseje realizar uma operação e não possua o token, ele deve esperar até que o token seja passado para ele. Assim que o banco detentor do token terminar de processar suas transações, ele passa o token para o próximo banco. Caso o próximo banco não esteja disponível, ele tentará enviar para o banco seguinte, e assim por diante. Caso nenhum banco esteja disponível, o token é mantido no banco atual, até que outros bancos fiquem ativos.

O código abaixo demonstra como a passagem do token é realizada. A primeira etapa é verificar quem é o próximo token do anel. Em seguida, a função `findNextValidBank` é responsável por encontrar o próximo banco válido, ou seja, que esteja online. Caso ele encontre um banco válido, é feito um multicast para todos os bancos avisando quem é o novo detentor do token. Caso contrário, o banco atual faz um multicast avisando que o token continua com ele.

```go
func PassToken() {
	// envia a transação para o próximo banco
	nextBank := storage.Ring.Next(config.Env.BankId)
	if nextBank == nil { ... }

	nextBankId := findNextValidBank(nextBank.Id)
	if nextBankId == nil {
		slog.Info("Não conseguiu passar o token para um banco valido. Mantém o token localmente.")
		BroadcastToken(config.Env.BankId) // faz o broadcast do token para os outros bancos - para garantir que o token não se perca
		return
	}

	BroadcastToken(*nextBankId)
}
```

### Detecção e recuperação de falhas
Todos os sistemas são suscetíveis a falhas, principalmente sistemas distribuídos. Desse modo, é imprescindível a esses sistema que as falhas seja detectadas e o corrigidas. Com o método de Token Ring, falhas como queda de algum dos bancos do sistema, perda/duplicação de token e transações não realizadas são os erros mais prováveis. Desse modo, mecanismos de detecção e recuperação de falhas foram introduzidos ao sistema.

<div align="center">
<img src="./images/token-ring-bank-down.png" alt="Detecção e recuperação de falhas" height="300px" width="auto" /> <br/>
<em>Figura 10. Detecção e recuperação de falhas</em>
</div>

#### Queda de algum dos bancos
Devido a natureza distribuída do sistema e da instabilidade da rede, o sistema precisa lidar com possíveis instabilidade e quedas de algum dos nós da rede. Assim, quando o banco detentor do token passar para o próximo banco e esse estiver indisponível, ele tentará repassar para o banco seguinte. Caso este também esteja indisponível, o processo irá se reiniciar com o próximo banco, até algum banco responder.

Caso nenhum banco esteja disponível, o token é mantido no banco atual, até que algum banco esteja disponível.

#### Perda de token
Caso o banco que possua o token venha a cair antes de repassar o token, o primeiro banco a nota a ausência do token é responsável por criar um novo e avisar a todos. Para isso, ele utiliza o horário de criação do token (estrutura `Ts` do [Token](#estrutura-do-token)) somado a um tempo X, que varia de acordo com o banco, como referência para verificar se o token foi perdido. O tempo X é definido como `2^ID do banco`, garantindo que o banco com ID menor tenha prioridade.

```go
// se o tempo de espera para o token for excedido
// o primeiro banco a perceber solicita o token
// bancos com IDs menores têm prioridade
bankTokenPriority := math.Pow(2, float64(config.Env.BankId))
maxTokenWaitDuration := time.Duration(float64(constants.MaxWaitTimeForTokenInterBank) + bankTokenPriority)
if time.Since(storage.Token.Get().Ts) > maxTokenWaitDuration {
   slog.Info("Tempo de espera para token interbancário excedido. Solicitando token...")
   services.BroadcastToken(config.Env.BankId) // faz um broadcast a todos os bancos avisando que o token agora é do banco atual
}
```

Dado que o tempo de espera `bankTokenPriority` é exponencial, a chance de dois bancos pedirem o token ao mesmo tempo é baixa. Além disso, caso ocorra, o mecanismo de detecção de [duplicação de tokens](#duplicação-de-token) invalidaria o segundo. Isso garante que o token nunca seja perdido e o sistema continue funcionando.

#### Duplicação de token
Para garantir que não ocorra duplicação de tokens, antes de iniciar o processamento das transações o banco envia um multicast para todos os bancos do consórcio perguntando quem é o dono do token. Caso ele mesmo seja o dono do token, é iniciado o processamento das transações. Caso contrário, ele atualiza as informações internas sobre quem é o dono do token, e cancela o processamento das transações.

Desse modo, quando um banco que caiu voltar a ativa, ele sempre vai ter a atualização mais recente sobre o token.

## Comunicação
Como forma de padronizar a comunicação tanto entre a interface e o banco, quanto entre os bancos do consórcio, foi utilizado o padrão de API REST. O uso de APIs REST permite que as operações sejam realizadas de forma simples e eficiente, além de garantir a interoperabilidade entre diferentes sistemas.

### Comunicação entre a interface e o banco
Como forma de padronizar a comunicação entre a interface e o banco, foi definido um conjunto de rotas que são utilizadas para realizar as operações de criação de conta, login, depósito, transferência, visualização de extrato e visualização de saldo. Para realizar essas operações, a interface faz requisições HTTP para o banco, que por sua vez processa a requisição e retorna uma resposta. Essas requisições seguem o padrão REST, onde cada operação é realizada através de um método HTTP específico.

#### Rotas da API
##### POST /api/accounts
Esta [rota](/bank/internal/routes/bank/create-account.go) é utilizada para criar uma nova conta no banco. Ela recebe um objeto JSON contendo o nome do usuário (ou razão social), os documentos (CPF ou CNPJ, podendo ser mais de um em caso de conta conjunta) e o tipo da conta (pessoa física, jurídica ou conjunta). O banco então cria a conta e retorna um objeto JSON contendo os dados da conta criada.

Exemplo de requisição para criar uma conta individual:
```http
POST /api/accounts

{
  "name": "João da Silva",
  "documents": ["000.000.000-01"],
  "type": "individual"
}
```

Exemplo de resposta:
```http
201 Created

{
  "id": 1,
  "name": "José da Silva",
  "documents": [
    "000.000.000-01"
  ],
  "type": "individual",
  "ibk": "1-1",
  "created_at": "2024-06-28T19:54:46.16347775-03:00",
  "balance": "0"
}
```

Exemplo de requisição para criar uma conta conjunta:
```http
POST /api/accounts

{
  "name": "João e Maria da Silva",
  "documents": ["000.000.000-01", "000.000.000-01"],
  "type": "joint"
}
```

##### POST /api/accounts/auth
A [rota](/bank/internal/routes/bank/auth.go) de autenticação é utilizada para realizar o login de um usuário. Ela recebe um objeto JSON contendo o InterBank Key (IBK) do usuário. O banco então verifica se o IBK é válido e retorna um objeto JSON contendo os dados da conta.

Exemplo de requisição:
```http
POST /api/accounts/auth

{
  "ibk": "1-1"
}
```

Exemplo de resposta:
```http
200 OK

{
  "id": 1,
  "name": "José da Silva",
  "documents": [
    "000.000.000-01"
  ],
  "type": "individual",
  "ibk": "1-1",
  "created_at": "2024-06-28T19:54:46.16347775-03:00",
  "balance": "0"
}
```

Exemplo de resposta para IBK inválido:
```http
401 Unauthorized

{
  "message": "Conta não encontrada"
}
```

##### DELETE /api/accounts/:id
Esta [rota](/bank/internal/routes/bank/delete-account.go) é utilizada para deletar uma conta no banco. Ela recebe o ID da conta a ser deletada e o banco então deleta a conta e retorna um objeto JSON contendo os dados da conta deletada.

Exemplo de requisição:
```http
DELETE /api/accounts/1
```

Exemplo de resposta:
```http
200 OK

{
  "message": "Conta deletada com sucesso"
}
```

##### POST /api/payments/deposit
Esta [rota](/bank/internal/routes/bank/deposit-route.go) é utilizada para realizar um depósito em uma conta. Ela recebe um objeto JSON contendo o IBK da conta e o valor a ser depositado. O banco então realiza o depósito e retorna um objeto JSON contendo os dados da transação. Como o deposito é uma operação atômica, a transação é realizada de forma instantânea, não dependendo do InterBank para ser concluída.

Exemplo de requisição:
```http
POST /api/payments/deposit

{
  "acc_ibk": "1-1",
  "amount": "100"
}
```

Exemplo de resposta:
```http
200 OK

{
  "id": "f64417e9-683c-4de8-a74b-31133002a808",
  "owner": "1-1",
  "parent_id": null,
  "type": "final",
  "operations": [
    {
      "id": "69600a35-2083-43ac-ba5d-a8a62467eab1",
      "from": "1-1",
      "to": "1-1",
      "type": "deposit",
      "amount": "100",
      "status": "success",
      "created_at": "2024-06-28T20:03:24.38772134-03:00",
      "updated_at": "2024-06-28T20:03:24.387726119-03:00"
    }
  ],
  "created_at": "2024-06-28T20:03:24.387722803-03:00",
  "updated_at": "2024-06-28T20:03:24.387726249-03:00",
  "status": "success"
}
```

##### GET /api/accounts/:id
Esta [rota](/bank/internal/routes/bank/find-account.go) é utilizada para visualizar os dados de uma conta. Ela recebe o ID da conta e o banco então retorna um objeto JSON contendo os dados da conta.

Exemplo de requisição:
```http
GET /api/accounts/1
```

Exemplo de resposta:
```http
200 OK

{
  "id": 1,
  "name": "José da Silva",
  "documents": [
  "000.000.000-01"
  ],
  "type": "individual",
  "ibk": "1-1",
  "created_at": "2024-06-28T20:03:19.09762633-03:00",
  "balance": "100"
}
```

##### GET /api/accounts/:id/all
Esta [rota](/bank/internal/routes/bank/find-all-accounts.go) é responsável por retornar todas as contas do usuário (tanto as contas no banco atual, quanto as contas em outros bancos). Ela recebe o ID da conta e o banco é responsável por enviar um multicast para todos os bancos do consórcio, solicitando as contas do usuário. Cada banco então retorna um objeto JSON contendo os dados da conta.

Exemplo de requisição:
```http
GET /api/accounts/1/all
```

Exemplo de resposta:
```http
200 OK

[
  {
    "id": 1,
    "name": "José da Silva",
    "documents": [
      "000.000.000-01"
    ],
    "type": "individual",
    "ibk": "1-1",
    "created_at": "2024-06-28T20:03:19.09762633-03:00",
    "balance": "100"
  },
  {
    "id": 1,
    "name": "José da Silva",
    "documents": [
      "000.000.000-01"
    ],
    "type": "individual",
    "ibk": "2-1",
    "created_at": "2024-06-28T20:08:00.643185361-03:00",
    "balance": "0"
  }
]
```

##### GET /api/accounts/:id/transactions
Esta [rota](/bank/internal/routes/bank/list-account-transactions.go) é utilizada para visualizar as transações de uma conta. Ela recebe o ID da conta e retorna um objeto JSON contendo as transações realizadas pela conta.

Exemplo de requisição:
```http
GET /api/accounts/1/transactions
```

Exemplo de resposta:
```http
200 OK

[
  {
    "id": "5f86cc57-a57b-4bcc-b707-3df28eaa679c",
    "owner": "1-1",
    "type": "package",
    "parent_id": null,
    "operations": [
      {
        "id": "153e7f58-0764-4767-bde6-cd5ac00ae73b",
        "from": "1-1",
        "to": "2-1",
        "type": "transfer",
        "amount": "50",
        "status": "failed",
        "created_at": "2024-06-28T20:08:52.768665384-03:00",
        "updated_at": "2024-06-28T20:08:53.043557766-03:00"
      },
      {
        "id": "ad5e4166-734a-4669-87b3-337c09def0e6",
        "from": "2-1",
        "to": "2-2",
        "type": "transfer",
        "amount": "100",
        "status": "failed",
        "created_at": "2024-06-28T20:08:52.7687795-03:00",
        "updated_at": "2024-06-28T20:08:53.043558096-03:00"
      }
    ],
    "created_at": "2024-06-28T20:08:52.768782115-03:00",
    "updated_at": "2024-06-28T20:08:53.043558226-03:00",
    "status": "failed"
  },
  {
    "id": "f64417e9-683c-4de8-a74b-31133002a808",
    "owner": "1-1",
    "type": "final",
    "parent_id": null,
    "operations": [
      {
        "id": "69600a35-2083-43ac-ba5d-a8a62467eab1",
        "from": "1-1",
        "to": "1-1",
        "type": "deposit",
        "amount": "100",
        "status": "success",
        "created_at": "2024-06-28T20:03:24.38772134-03:00",
        "updated_at": "2024-06-28T20:03:24.387726119-03:00"
      }
    ],
    "created_at": "2024-06-28T20:03:24.387722803-03:00",
    "updated_at": "2024-06-28T20:03:24.387726249-03:00",
    "status": "success"
  }
]
```

##### POST /api/payments/pay
Esta [rota](/bank/internal/routes/bank/pay-route.go) é utilizada para realizar uma transferência entre contas. Ela recebe um objeto JSON contendo o IBK do autor da transação e uma lista de operações, onde cada operação contém o IBK da conta de origem, o IBK da conta de destino e o valor a ser transferido. O banco então realiza a transferência e retorna um objeto JSON com os detalhes da transação. Como a transferência depende do InterBank, a transação é criada de forma assíncrona e processada em background.

Exemplo de requisição:
```http
POST /api/payments/pay

{
  "author": "1-1",
  "operations": [
    {
      "from": "1-1",
      "to": "2-1",
      "amount": 50
    },
    {
      "from": "2-1",
      "to": "2-2",
      "amount": 100
    }
  ]
}
```

Exemplo de resposta:
```http
200 OK

{
  "id": "374837eb-9f72-49d3-ae27-aa5ea68c2fd9",
  "owner": "1-1",
  "type": "package",
  "parent_id": null,
  "operations": [
    {
      "id": "b04dc295-2de8-41dc-9cda-116fe07baeb0",
      "from": "1-1",
      "to": "2-1",
      "type": "transfer",
      "amount": "50",
      "status": "pending",
      "created_at": "2024-06-29T18:20:50.076733513Z",
      "updated_at": "2024-06-29T18:20:50.076733563Z"
    },
    {
      "id": "01e98458-6235-4a71-b369-af06b605ccff",
      "from": "2-1",
      "to": "2-2",
      "type": "transfer",
      "amount": "100",
      "status": "pending",
      "created_at": "2024-06-29T18:20:50.076876201Z",
      "updated_at": "2024-06-29T18:20:50.076876251Z"
    }
  ],
  "created_at": "2024-06-29T18:20:50.076883013Z",
  "updated_at": "2024-06-29T18:20:50.076883064Z",
  "status": "pending"
}
```

### Comunicação entre os bancos (InterBank)
Como forma de padronizar a comunicação entre os bancos do consórcio, foi definido um conjunto de padrões para a comunicação entre os bancos. O InterBank é responsável por garantir que as mensagens sejam enviadas e recebidas de forma correta, além de garantir a consistência dos dados.

#### Rotas do InterBank
##### POST /interbank/prepare
Esta [rota](/bank/internal/routes/interbank/prepare.go) é utilizada para preparar uma transação. Ela recebe um objeto JSON contendo a operação a ser realizada e qual o passo a ser realizado (débito ou crédito). O banco então prepara a transação e retorna um objeto JSON indicando que a transação foi preparada com sucesso.

Exemplo de requisição:
```http
POST /interbank/add-credit

{
  "parent_id": "f635b354-025b-4ab5-bf5c-a46a36998ebc",
  "operation": {
    "from": "1-1",
    "to": "2-1",
    "amount": 100
  },
  "step": "credit"
}
```

Exemplo de resposta:
```http
200 OK
{
 "id": "5f86cc57-a57b-4bcc-b707-3df28eaa679c",
 "owner": "1-1",
 "type": "final",
 "parent_id": "f635b354-025b-4ab5-bf5c-a46a36998ebc",
 "operations": [
   {
     "id": "153e7f58-0764-4767-bde6-cd5ac00ae73b",
     "from": "1-1",
     "to": "2-1",
     "type": "transfer",
     "amount": "100",
     "status": "pending",
     "created_at": "2024-06-28T20:08:52.768665384-03:00",
     "updated_at": "2024-06-28T20:08:53.043557766-03:00"
   },
 ],
 "created_at": "2024-06-28T20:08:52.768782115-03:00",
 "updated_at": "2024-06-28T20:08:53.043558226-03:00",
 "status": "pending"
},
```

##### POST /interbank/rollback
Esta [rota](/bank/internal/routes/interbank/rollback.go) é utilizada para reverter uma transação. Ela recebe um objeto JSON contendo o ID da transação a ser revertida e qual o passo a ser revertido (débito ou crédito). O banco então reverte a transação e retorna um objeto JSON indicando que a transação foi revertida com sucesso.

Exemplo de requisição:
```http
POST /interbank/rollback

{
  "tx_id": "5f86cc57-a57b-4bcc-b707-3df28eaa679c",
  "step": "credit"
}
```

Exemplo de resposta:
```http
200 OK
{
 "id": "5f86cc57-a57b-4bcc-b707-3df28eaa679c",
 "owner": "1-1",
 "type": "final",
 "parent_id": "f635b354-025b-4ab5-bf5c-a46a36998ebc",
 "operations": [
   {
     "id": "153e7f58-0764-4767-bde6-cd5ac00ae73b",
     "from": "1-1",
     "to": "2-1",
     "type": "transfer",
     "amount": "100",
     "status": "failed",
     "created_at": "2024-06-28T20:08:52.768665384-03:00",
     "updated_at": "2024-06-28T20:08:53.043557766-03:00"
   },
 ],
 "created_at": "2024-06-28T20:08:52.768782115-03:00",
 "updated_at": "2024-06-28T20:08:53.043558226-03:00",
 "status": "failed"
},
```

##### POST /interbank/commit
Esta [rota](/bank/internal/routes/interbank/commit.go) é utilizada para confirmar uma transação. Ela recebe um objeto JSON contendo o ID da transação a ser confirmada e qual o passo a ser confirmado (débito ou crédito). O banco então confirma a transação e retorna um objeto JSON indicando que a transação foi confirmada com sucesso.

Exemplo de requisição:
```http
POST /interbank/add-credit

{
  "tx_id": "5f86cc57-a57b-4bcc-b707-3df28eaa679c",
  "step": "credit"
}
```

Exemplo de resposta:
```http
200 OK
{
 "id": "5f86cc57-a57b-4bcc-b707-3df28eaa679c",
 "owner": "1-1",
 "type": "final",
 "parent_id": "f635b354-025b-4ab5-bf5c-a46a36998ebc",
 "operations": [
   {
     "id": "153e7f58-0764-4767-bde6-cd5ac00ae73b",
     "from": "1-1",
     "to": "2-1",
     "type": "transfer",
     "amount": "100",
     "status": "success",
     "created_at": "2024-06-28T20:08:52.768665384-03:00",
     "updated_at": "2024-06-28T20:08:53.043557766-03:00"
   },
 ],
 "created_at": "2024-06-28T20:08:52.768782115-03:00",
 "updated_at": "2024-06-28T20:08:53.043558226-03:00",
 "status": "success"
},
```

##### GET /account/:document
Esta rota é utilizada para buscar todas as contas que um banco possui associadas a um documento. Ela recebe o documento (CPF ou CNPJ) e o banco então retorna um objeto JSON contendo as contas associadas ao documento.

Exemplo de requisição:
```http
GET /account/000.000.000-01
```

Exemplo de resposta:
```http
200 OK

[
  {
    "id": 1,
    "name": "José da Silva",
    "documents": [
      "000.000.000-01"
    ],
    "type": "individual",
    "ibk": "1-1",
    "created_at": "2024-06-28T20:03:19.09762633-03:00",
    "balance": "100"
  }
]
```

##### POST /interbank/account/ibk/:ibk
Esta [rota](/bank/internal/routes/interbank/account-ibk.go) é utilizada para buscar uma conta em um banco específico. Ela recebe o IBK da conta e o banco então retorna um objeto JSON contendo os dados da conta.

Exemplo de requisição:
```http
POST /interbank/account/ibk/1-1
```

Exemplo de resposta:
```http
200 OK

{
  "id": 1,
  "name": "José da Silva",
  "documents": [
    "000.000.000-01"
  ],
  "type": "individual",
  "ibk": "1-1",
  "created_at": "2024-06-28T20:03:19.09762633-03:00",
  "balance": "100"
}
```

##### PUT /interbank/token
Esta [rota](/bank/internal/routes/interbank/set-token.go) é utilizada para enviar um token para um banco específico. Ela recebe um objeto JSON contendo o ID do banco que irá receber o token e a data de quando o token foi criado.

Exemplo de requisição:
```http
PUT /interbank/token

{
  "to": 1,
  "ts": "2024-06-28T20:03:19.09762633-03:00"
}
```

Exemplo de resposta:
```http
200 OK

{
  "message": "Token setado com sucesso"
}
```

##### GET /interbank/token
Esta [rota](/bank/internal/routes/interbank/get-token.go) é utilizada para retornar quem o banco acha que é o dono do token.

Exemplo de requisição:
```http
GET /interbank/token
```

Exemplo de resposta:
```http
200 OK

{
  "to": 1,
  "ts": "2024-06-28T20:03:19.09762633-03:00"
}
```

##### GET /interbank/token/ok
Esta [rota](/bank/internal/routes/interbank/can-receive-token.go) verifica se o banco pode receber o token.

Exemplo de requisição:
```http
GET /interbank/token/ok
```

Exemplo de resposta:
```http
200 OK
```

## Sincronização dos dados internamente no banco
Devido a natureza distribuída do sistema, leituras e escritas podem ocorrer de forma concorrente no banco. Por exemplo, dois usuário podem tentar realizar um depósito na mesma conta ao mesmo tempo, o que pode causar inconsistências nos dados.

Para resolver o problema de sincronia interna, foram utilizados [mecanismos de lock](https://github.com/jnaraujo/tec502-inter-bank/blob/main/bank/internal/storage/transactions.go#L15) (mutexes) para garantir que apenas uma operações de escrita seja realizada por vez. Assim, antes que qualquer operação de leitura ou escrita no dados armazenados seja realizada, um lock é adquirido. Isso garante que as operações sejam realizadas de forma ordenada e sem conflitos.

Por exemplo, na operação de depósito, o lock é adquirido antes de adicionar o valor na conta e liberado após a operação ser concluída. Isso garante, no caso abaixo, que apenas uma transação seja salva por vez.

```go
// Código de bank/internal/storage/transactions.go
func (ts *transactionsStorage) Save(tr models.Transaction) {
   ts.mu.Lock()
   ts.data[tr.Id] = tr
   ts.mu.Unlock()
}
```

## Testes
Para garantir que o sistema de consórcio bancário funcione corretamente, foram implementados testes unitários e de integração. Os testes unitários são utilizados para testar funções específicas do código, enquanto os testes de integração são utilizados para testar a integração entre diferentes componentes do sistema.

Por exemplo, para testar transações simultâneas entre os bancos, foram criados testes de integração que simulam a criação de transações em diferentes bancos ao mesmo tempo. Os testes verificam se as transações são realizadas de forma ordenada e sem conflitos, mesmo que diferentes usuários estejam realizando transações ao mesmo tempo. Para isso, foram utilizadas threads para simular a criação de transações em diferentes bancos ao mesmo tempo. Esses testes foram implementados em Python e estão disponíveis no diretório `tests`.

## Conclusão
O sistema de consórcio bancário desenvolvido é uma solução eficiente e segura para a realização de transações financeiras entre diferentes bancos. A utilização de APIs REST, transações atômicas, Token Ring e transações assíncronas garantem que as operações sejam realizadas de forma ordenada e sem conflitos, mesmo que diferentes usuários estejam realizando transações ao mesmo tempo. Além disso, foram utilizadas utilizando tecnologias modernas e atuais, como React, Go e Docker.

Desse modo, o sistema desenvolvido é capaz de lidar com os requisitos propostos, além de ser capaz de lidar com as possíveis falhas do sistema de maneira satisfatória. Além disso, o sistema passou em todos os testes de transações concorrentes e atômicas.

