<div align="center">
  <h1>InterBank - Sistema de pagamentos bancários descentralizado</h1>
  <p>
    <strong>Projeto desenvolvido para a disciplina TEC502 - MI Concorrência e Conectividade</strong>
  </p>

  ![Most used language](https://img.shields.io/github/languages/top/jnaraujo/tec502-inter-bank?style=flat-square)
  ![GitHub](https://img.shields.io/github/license/jnaraujo/tec502-inter-bank)
</div>

Nos últimos anos, o cenário bancário no Brasil tem passado por uma revolução digital significativa. A criação do sistema de pagamento instantâneo Pix, por exemplo, foi um marco importante para a modernização do sistema financeiro nacional. Conforme relatado pelo Banco Central do Brasil, a adesão dos brasileiros a essas novas formas de movimentações financeiras tem sido expressiva, promovendo a inclusão financeira de diversas camadas da população que anteriormente dependiam de métodos tradicionais como boletos e cheques.

Como forma de criar uma solução descentralizada para pagamentos bancários, o InterBank foi desenvolvido. O sistema é composto por uma rede de bancos (nós) que se comunicam entre si para realizar transações financeiras de forma segura e eficiente. Assim, qualquer usuário usando qualquer um dos bancos participantes do consórcio, pode realizar transações entre suas suas contas, independente da instituição financeira de origem. Desse modo, um usuário pode criar transações com operações em diferentes bancos, de forma atômica e consistente.

Neste contexto, o método de Token Ring foi escolhido como a solução para resolver o problema de concorrência entre os bancos participantes. Este método, amplamente utilizado em redes de computadores, garante que cada banco tenha a oportunidade de acessar e atualizar as informações das contas de forma ordenada e sem conflitos, prevenindo o "duplo gasto" e assegurando a consistência dos dados. Além disso, para o desenvolvimento do projeto, foram utilizadas tecnologias como Docker, ReactJS e Go.

## Sumário

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

### Como utilizar?
#### Como criar de uma conta
#### Como realizar um depósito
#### Como realizar uma transferência
#### Como visualizar o extrato de uma conta
#### Como visualizar o saldo de uma conta

## Arquitetura do projeto
O sistema foi divido em duas partes principais: a interface gráfica do banco e o código do banco. A interface gráfica foi desenvolvida utilizando ReactJS e a biblioteca TanStack Query para gerenciamento de estado e requisições HTTP. Já o código do banco foi desenvolvido em Go, utilizando o framework Fiber.

### Interface gráfica
A interface gráfica do banco foi desenvolvida utilizando ReactJS e a biblioteca TanStack Query para gerenciamento de estado e requisições HTTP. A interface é composta por 3 páginas principais: a página de criação de conta, a página de login e a página do banco, onde é possível realizar operações como depósito, transferência, visualização de extrato e visualização de saldo.

```bash
client
├── src # Código fonte da aplicação
│   ├── @types # Alguns tipos que são utilizados em várias partes da aplicação 
│   ├── components # Componentes React reutilizáveis
│   ├── constants # Definem valores constantes, como o tempo de atualização do saldo
│   ├── contexts # Contextos são utilizados para compartilhar informações entre componentes de forma global
│   ├── hooks # Hooks são funções que permitem adicionar funcionalidades ao componente
│   ├── lib # Funções utilitárias
│   ├── routes # Definição das rotas da aplicação
│   ├── schemas # Schemas utilizados para validação de dados
│   ├── services # Funções que realizam requisições HTTP
```

### Código do banco
O código do banco foi desenvolvido em Go, utilizando o framework Fiber. Nele estão implementadas as rotas que são utilizadas para realizar as operações de criação de conta, login, depósito, transferência, visualização de extrato e visualização de saldo. Além disso, nele estão os códigos referentes ao sistema de Token Ring, que é utilizado para garantir a consistência dos dados, e o InterBank, que é responsável por garantir a comunicação entre os bancos.

```bash
broker
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
│   ├── transaction_processor # Serviço que roda em background para processar as transações
|   ├── validate # Funções para validação de dados
```

## Comunicação
### Protocolo de comunicação
### Token Ring
### InterBank

## Rotas da API