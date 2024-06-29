import utils
import api
import time

def singleTransactionWithMultipleOperations(addrs=["localhost:3001", "localhost:3002"]):
  cpf_1_1 = utils.randomCpf()
  cnpj_1_2 = utils.randomCpf()
  cpf_2_3 = utils.randomCpf()
  
  acc_1_1 = api.createAccount("José da Silva", [cpf_1_1], "individual", addrs[0])
  acc_1_2 = api.createAccount("Pedro Souza", [cnpj_1_2], "legal", addrs[0])
  acc_1_joint = api.createAccount("José e Pedro Conta Conjunta", [cpf_1_1, cnpj_1_2], "joint", addrs[0])
  acc_2_1 = api.createAccount("José da Silva", [cpf_1_1], "individual", addrs[1])
  acc_2_3 = api.createAccount("Frederico Machado", [cpf_2_3], "individual", addrs[1])
  
  api.createDeposit(acc_1_joint["ibk"], 100, addrs[0])
  api.createDeposit(acc_1_1["ibk"], 100, addrs[0])
  
  api.pay(acc_1_1["ibk"], [
    {
      "from_acc_ibk": acc_1_joint["ibk"],
      "to_acc_ibk": acc_2_1["ibk"],
      "amount": 50,
    },
    {
      "from_acc_ibk": acc_1_1["ibk"],
      "to_acc_ibk": acc_2_1["ibk"],
      "amount": 100,
    },
    {
      "from_acc_ibk": acc_2_1["ibk"],
      "to_acc_ibk": acc_1_2["ibk"],
      "amount": 150,
    },
    {
      "from_acc_ibk": acc_1_joint["ibk"],
      "to_acc_ibk": acc_2_3["ibk"],
      "amount": 15,
    }
  ], addrs[0])
  
  time.sleep(1.5) # Espera um pouco para as transações serem processadas
  # 1.5 segundos por que tem um delay de 1 segundo nos bancos para processar as transações!!
  
  # trs = api.findAllTransactions(acc_1_1["id"], addrs[0])
  # print(trs)
  
  # Verifica se as contas estão corretas
  acc_1_1 = api.findAccount(acc_1_1["id"], addrs[0])
  acc_1_2 = api.findAccount(acc_1_2["id"], addrs[0])
  acc_1_joint = api.findAccount(acc_1_joint["id"], addrs[0])
  acc_2_1 = api.findAccount(acc_2_1["id"], addrs[1])
  acc_2_3 = api.findAccount(acc_2_3["id"], addrs[1])
  
  if int(acc_1_1["balance"]) != 0:
    print("Erro: Saldo de José da Silva no banco 1 incorreto")
  if int(acc_1_2["balance"]) != 150:
    print("Erro: Saldo de Pedro Souza no banco 1 incorreto")
  if int(acc_1_joint["balance"]) != 35:
    print("Erro: Saldo da conta conjunta no banco 1 incorreto")
  if int(acc_2_1["balance"]) != 0:
    print("Erro: Saldo de José da Silva no banco 2 incorreto")
  if int(acc_2_3["balance"]) != 15:
    print("Erro: Saldo de Frederico Machado no banco 2 incorreto")
  
  # Tear down
  a = api.deleteUser(acc_1_1["id"], addrs[0])
  api.deleteUser(acc_1_2["id"], addrs[0])
  api.deleteUser(acc_1_joint["id"], addrs[0])
  api.deleteUser(acc_2_1["id"], addrs[1])
  api.deleteUser(acc_2_3["id"], addrs[1])