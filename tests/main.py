from threading import Thread
import tests
import utils
  
def main():
  print("="*5, "Running Tests", "="*5)
  
  threads = []
  for i in range(10):
    addrs = utils.randomAddrs()
    t = Thread(target=tests.singleTransactionWithMultipleOperations, args=(addrs,))
    t.start()
    threads.append(t)
    
  # Espera todas as threads terminarem
  for t in threads:
    t.join()
    
  print("="*5, "Tests Finished", "="*5)
main()