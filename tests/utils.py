import random

def randomCpf():
  cpf = ""
  for i in range(11):
    cpf += str(random.randint(0, 9))
  return cpf

def shuffleAddrs(addrs):
  a = addrs.copy()
  random.shuffle(a)
  return a