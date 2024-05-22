## Rotas internas do banco
### Autenticação
- POST /auth/login - Autentica o usuário

### Contas
- POST /accounts - Cria uma conta
- GET /accounts/:id - Retorna os dados de uma conta
- GET /accounts/:id/balance - Retorna o saldo de uma conta
- GET /accounts/:id/transactions - Retorna as transações de uma conta

### Transferências
- GET /payments/:id - Retorna os dados de uma transferência
- POST /payments/transfer - Realiza uma transferência entre contas
- POST /payments/pay - Realiza um pagamento que não seja transferência (boleto, fatura, etc)
- POST /payments/deposit - Realiza um depósito em uma conta

## Rotas do InterBank
### Transferências
- POST /interbank/transfer - Realiza uma transferência entre contas de bancos diferentes


# Protocolo de comunicação InterBank
O InterBank é um protocolo de comunicação que permite a transferência de valores entre contas de bancos diferentes. Ele é baseado em requisições HTTP e JSON.

## Chave do usuário
Cada usuário possui uma chave única de 4 bytes (32 bits) que é utilizada para identificá-lo. Essa chave é gerada pelo banco e é utilizada para identificar o usuário em todas as transações. A chave é formada pela concatenação de um identificador do banco (1 byte) e um identificador do usuário (3 bytes).

### Identificador do banco
O identificador do banco é um número inteiro não negativo que identifica o banco. Ele é utilizado para garantir que a chave do usuário seja única entre todos os bancos. O identificador do banco é um número de 1 byte (8 bits).

### Identificador do usuário
O identificador do usuário é um número inteiro não negativo que identifica o usuário dentro do banco. Ele é utilizado para garantir que a chave do usuário seja única dentro do banco. O identificador do usuário é um número de 3 bytes (24 bits).

### Exemplo
Suponha que o banco tenha o identificador 1 e o usuário tenha o identificador 1234. A chave do usuário seria a concatenação de 1 e 1234, ou seja, 1-1234.