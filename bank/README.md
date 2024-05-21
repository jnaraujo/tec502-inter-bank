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