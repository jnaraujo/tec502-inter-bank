from threading import Thread
import tests
import utils
import random

addrs = ["localhost:3001", "localhost:3002", "localhost:3003"]

def shuffleAddrs():
  a = addrs.copy()
  random.shuffle(a)
  return a

def main():
  print("="*5, "Running Tests", "="*5)
  
  tests.multipleTransactions1(addrs)
  tests.multipleTransactions2(addrs)
  
  threads = []
  for i in range(10):
    t = Thread(target=tests.singleTransactionWithMultipleOperations, args=(shuffleAddrs(),))
    t.start()
    threads.append(t)
    
  # Espera todas as threads terminarem
  for t in threads:
    t.join()
    
  print("="*5, "Tests Finished", "="*5)
main()