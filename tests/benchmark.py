import api
import threading
import time
from datetime import datetime
import utils

addrs = ["localhost:3001", "localhost:3002", "localhost:3003"]

N = 500

def main():
  print("="*5, "Running Benchmark", "="*5)
  cpf_1_1 = utils.randomCpf()
  cpf_3_3 = utils.randomCpf()
  
  acc_1_1 = api.createAccount("José da Silva", [cpf_1_1], "individual", addrs[0])
  acc_2_1 = api.createAccount("José da Silva", [cpf_1_1], "individual", addrs[1])
  acc_3_3 = api.createAccount("Frederico Machado", [cpf_3_3], "individual", addrs[2])
  
  api.createDeposit(acc_1_1["ibk"], N, addrs[0])
  api.createDeposit(acc_3_3["ibk"], N, addrs[2])
  
  # cria N threads com as transações
  threads = []
  for i in range(N):
    t1 = threading.Thread(target=api.pay, args=(acc_1_1["ibk"], [
        {
          "from": acc_1_1["ibk"],
          "to": acc_2_1["ibk"],
          "amount": 1,
        },
        {
          "from": acc_2_1["ibk"],
          "to": acc_3_3["ibk"],
          "amount": 1,
        }
      ], addrs[0]))
    
    t2 = threading.Thread(target=api.pay, args=(acc_3_3["ibk"], [
        {
          "from": acc_3_3["ibk"],
          "to": acc_1_1["ibk"],
          "amount": 1,
        },
      ], addrs[2]))
    
    threads.append(t1)
    threads.append(t2)
    
  # inicia as threads
  for t in threads:
    t.start()
    
  # espera todas as threads terminarem
  for t in threads:
    t.join()
    
  # espera um pouco para as transações serem processadas
  time.sleep(10)
  
  txsFromAcc1 = api.findAllTransactions(acc_1_1["id"], addrs[0])
  txsFromAcc2 = api.findAllTransactions(acc_2_1["id"], addrs[1])
  txsFromAcc3 = api.findAllTransactions(acc_3_3["id"], addrs[2])
  txs = txsFromAcc1 + txsFromAcc2 + txsFromAcc3
  
  timeToProcessInMs = 0
  txsProcessed = 0
  for tx in txs:
    if tx["status"] != "success":
      print("Transaction not processed: ", tx)
      return
    if tx["type"] == "final":
      continue
    
    txsProcessed += 1
    
    created_at = datetime.strptime(tx["created_at"][:-4] + "Z", "%Y-%m-%dT%H:%M:%S.%fZ").timestamp()
    updated_at = datetime.strptime(tx["updated_at"][:-4] + "Z", "%Y-%m-%dT%H:%M:%S.%fZ").timestamp()
    timeToProcessInMs += (updated_at - created_at) * 1000
    
  print("Tempo médio de processamento: ", round(timeToProcessInMs / txsProcessed, 2), "ms")
  print("Transações processadas: ", txsProcessed)
  
  print("="*5, "Benchmark Finished", "="*5)
  
  # Tear down
  api.deleteUser(acc_1_1["id"], addrs[0])
  api.deleteUser(acc_2_1["id"], addrs[1])
  api.deleteUser(acc_3_3["id"], addrs[2])
main()