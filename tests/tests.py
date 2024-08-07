import utils
import api
import time
from threading import Thread

WAIT_TIME_FOR_TRANSACTIONS = 5

def multipleTransactions1(addrs):
  cpf_1_1 = utils.randomCpf()
  cpf_2_2 = utils.randomCpf()
  cpf_3_3 = utils.randomCpf()
  
  acc_1_1 = api.createAccount("José da Silva", [cpf_1_1], "individual", addrs[0])
  acc_2_2 = api.createAccount("Frederico Machado", [cpf_2_2], "individual", addrs[1])
  acc_3_3 = api.createAccount("Maria de Souza", [cpf_3_3], "individual", addrs[2])
  
  api.createDeposit(acc_1_1["ibk"], 100, addrs[0])
  api.createDeposit(acc_2_2["ibk"], 100, addrs[1])
  api.createDeposit(acc_3_3["ibk"], 100, addrs[2])
  
  t1 = Thread(target=api.pay, args=(acc_1_1["ibk"], [
      {
        "from": acc_1_1["ibk"],
        "to": acc_2_2["ibk"],
        "amount": 50,
      },
      {
        "from": acc_1_1["ibk"],
        "to": acc_2_2["ibk"],
        "amount": 50,
      }
    ], addrs[0]))
  
  t2 = Thread(target=api.pay, args=(acc_2_2["ibk"], [
      {
        "from": acc_2_2["ibk"],
        "to": acc_1_1["ibk"],
        "amount": 20,
      },
      {
        "from": acc_2_2["ibk"],
        "to": acc_1_1["ibk"],
        "amount": 25,
      }
    ], addrs[1]))
  
  t3 = Thread(target=api.pay, args=(acc_3_3["ibk"], [
      {
        "from": acc_3_3["ibk"],
        "to": acc_1_1["ibk"],
        "amount": 30,
      },
      {
        "from": acc_3_3["ibk"],
        "to": acc_2_2["ibk"],
        "amount": 30,
      },
    ], addrs[2]))
  
  # Inicia as threads
  t1.start()
  t2.start()
  t3.start()
  
  # Espera as threads terminarem
  t1.join()
  t2.join()
  t3.join()
  
  time.sleep(WAIT_TIME_FOR_TRANSACTIONS) # Espera um pouco para as transações serem processadas
  
  acc_1_1 = api.findAccount(acc_1_1["id"], addrs[0])
  acc_2_2 = api.findAccount(acc_2_2["id"], addrs[1])
  acc_3_3 = api.findAccount(acc_3_3["id"], addrs[2])

  if int(acc_1_1["balance"]) != 75:
    print("Erro: Saldo de José da Silva incorreto")
  if int(acc_2_2["balance"]) != 185:
    print("Erro: Saldo de Frederico Machado incorreto")
  if int(acc_3_3["balance"]) != 40:
    print("Erro: Saldo de Maria de Souza incorreto")
    
  # Tear down
  api.deleteUser(acc_1_1["id"], addrs[0])
  api.deleteUser(acc_2_2["id"], addrs[1])
  api.deleteUser(acc_3_3["id"], addrs[2])

def multipleTransactions2(addrs):
  cpf_1_1 = utils.randomCpf()
  cpf_2_2 = utils.randomCpf()
  cpf_3_3 = utils.randomCpf()
  
  acc_1_1 = api.createAccount("José da Silva", [cpf_1_1], "individual", addrs[0])
  acc_2_2 = api.createAccount("Frederico Machado", [cpf_2_2], "individual", addrs[1])
  acc_3_3 = api.createAccount("Maria de Souza", [cpf_3_3], "individual", addrs[2])
  
  api.createDeposit(acc_1_1["ibk"], 100, addrs[0])
  api.createDeposit(acc_2_2["ibk"], 50, addrs[1])
  
  threads = []
  for i in range(100):
    t1 = Thread(target=api.pay, args=(acc_1_1["ibk"], [
        {
          "from": acc_1_1["ibk"],
          "to": acc_2_2["ibk"],
          "amount": 0.5,
        },
        {
          "from": acc_1_1["ibk"],
          "to": acc_3_3["ibk"],
          "amount": 0.5,
        }
      ], addrs[0]))
    
    t2 = Thread(target=api.pay, args=(acc_2_2["ibk"], [
        {
          "from": acc_2_2["ibk"],
          "to": acc_1_1["ibk"],
          "amount": 0.4,
        },
        {
          "from": acc_2_2["ibk"],
          "to": acc_3_3["ibk"],
          "amount": 0.1,
        }
      ], addrs[1]))
    
    t1.start()
    t2.start()
    
    threads.append(t1)
    threads.append(t2)
  
  for t in threads:
    t.join()
  
  time.sleep(WAIT_TIME_FOR_TRANSACTIONS) # Espera um pouco para as transações serem processadas
  
  acc_1_1 = api.findAccount(acc_1_1["id"], addrs[0])
  acc_2_2 = api.findAccount(acc_2_2["id"], addrs[1])
  acc_3_3 = api.findAccount(acc_3_3["id"], addrs[2])
  
  if int(acc_1_1["balance"]) != 40:
    print("Erro: Saldo de José da Silva incorreto")
  if int(acc_2_2["balance"]) != 50:
    print("Erro: Saldo de Frederico Machado incorreto")
  if int(acc_3_3["balance"]) != 60:
    print("Erro: Saldo de Maria de Souza incorreto")
    
  # Tear down
  api.deleteUser(acc_1_1["id"], addrs[0])
  api.deleteUser(acc_2_2["id"], addrs[1])
  api.deleteUser(acc_3_3["id"], addrs[2])

def multipleTransactions3(addrs):
  threads = []
  for i in range(10):
    # addresses = utils.shuffleAddrs(addrs)
    t = Thread(target=singleTransactionWithMultipleOperations, args=(addrs,))
    t.start()
    threads.append(t)
    
  # Espera todas as threads terminarem
  for t in threads:
    t.join()

def singleTransactionWithMultipleOperations(addrs):
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
      "from": acc_1_joint["ibk"],
      "to": acc_2_1["ibk"],
      "amount": 50,
    },
    {
      "from": acc_1_1["ibk"],
      "to": acc_2_1["ibk"],
      "amount": 100,
    },
    {
      "from": acc_2_1["ibk"],
      "to": acc_1_2["ibk"],
      "amount": 150,
    },
    {
      "from": acc_1_joint["ibk"],
      "to": acc_2_3["ibk"],
      "amount": 15,
    }
  ], addrs[0])
  
  time.sleep(WAIT_TIME_FOR_TRANSACTIONS) # Espera um pouco para as transações serem processadas
  
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

def testFailureTransactions(addrs):
  cpf_1_1 = utils.randomCpf()
  cpf_2_2 = utils.randomCpf()
  
  acc_1_1 = api.createAccount("José da Silva", [cpf_1_1], "individual", addrs[0])
  acc_2_2 = api.createAccount("Frederico Machado", [cpf_2_2], "individual", addrs[1])
  
  api.createDeposit(acc_1_1["ibk"], 100, addrs[0])
  api.createDeposit(acc_2_2["ibk"], 100, addrs[1])
  
  # testa uma transação que tenta enviar mais dinheiro do que o saldo
  tx1 = api.pay(acc_1_1["ibk"], [{
    "from": acc_1_1["ibk"],
    "to": acc_2_2["ibk"],
    "amount": 150,
  }], addrs[0])
    
  # testa uma transação que tenta enviar dinheiro para uma conta inexistente
  tx2 = api.pay(acc_1_1["ibk"], [{
    "from": acc_1_1["ibk"],
    "to": "2-20",
    "amount": 50,
  }], addrs[0])
    
  # testa uma transação que tenta enviar dinheiro para uma conta de terceiros
  tx3 = api.pay(acc_1_1["ibk"], [{
    "from": acc_2_2["ibk"],
    "to": acc_1_1["ibk"],
    "amount": 50,
  }], addrs[0])
  if tx3["message"] != "Usuário não pode fazer transferências com contas de terceiros":
    print("Erro: Mensagem de erro incorreta")
  
  # testa uma transacao que tenta enviar dinheiro para a mesma conta
  tx4 = api.pay(acc_1_1["ibk"], [{
    "from": acc_1_1["ibk"],
    "to": acc_1_1["ibk"],
    "amount": 50,
  }], addrs[0])
  if tx4["message"] != "Conta de origem e destino não podem ser iguais":
    print("Erro: Mensagem de erro incorreta")
  
  time.sleep(WAIT_TIME_FOR_TRANSACTIONS) # Espera um pouco para as transações serem processadas
  
  # Verifica se as transações falharam
  transactions = api.findAllTransactions(acc_1_1["id"], addrs[0])
  
  for tx in transactions:
    if tx["id"] == tx1["id"] and tx["status"] != "failed":
      print("Erro: Transação 1 deveria ter falhado")
    if tx["id"] == tx2["id"] and tx["status"] != "failed":
      print("Erro: Transação 2 deveria ter falhado")
      
  # Tear down
  api.deleteUser(acc_1_1["id"], addrs[0])
  api.deleteUser(acc_2_2["id"], addrs[1])