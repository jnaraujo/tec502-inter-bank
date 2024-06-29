from threading import Thread
import random
import tests
import utils
  
def main():
  for i in range(10):
    addrs = utils.randomAddrs()
    t = Thread(target=executeTransactions, args=(addrs,))
    t.start()
main()