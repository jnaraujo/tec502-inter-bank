
def randomCpf():
  cpf = ""
  for i in range(11):
    cpf += str(random.randint(0, 9))
  return cpf

def randomAddrs():
  addrs = ["localhost:3001", "localhost:3002"]
  random.shuffle(addrs)
  return addrs