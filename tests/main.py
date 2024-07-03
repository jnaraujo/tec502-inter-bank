from threading import Thread
import tests
import utils
import random

addrs = ["localhost:3001", "localhost:3002", "localhost:3003"]

def main():
  print("="*5, "Running Tests", "="*5)
  
  print("Testing multiple transactions 1")
  tests.multipleTransactions1(addrs)
  print("Testing multiple transactions 2")
  tests.multipleTransactions2(addrs)
  print("Testing multiple transactions 3")
  tests.multipleTransactions3(addrs)
  
  print("Testing failure transactions")
  tests.testFailureTransactions(addrs)
    
  print("="*5, "Tests Finished", "="*5)
  
def stressTest(): 
  while True:
    print("Testing")
    tests.multipleTransactions1(addrs)
    tests.multipleTransactions2(addrs)
    tests.multipleTransactions3(addrs)
    tests.testFailureTransactions(addrs)
    print("Testing")
    

main()