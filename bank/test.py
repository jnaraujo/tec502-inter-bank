import http.client
import json
import random

def run():  
  # Create users
  print("Creating users...")
  user1 = create_user("Frederico dos Santos", randomCpf(), "3001")
  user2 = create_user("José da Silva", randomCpf(), "3002")
  
  # Create deposits
  print("Creating deposits...")
  create_deposit(user1["id"], 100, "3001")
  create_deposit(user2["id"], 100, "3002")
  
  # Find accounts
  print("Finding accounts...")
  fUser1 = find_account(user1["id"], "3001")
  fUser2 = find_account(user2["id"], "3002")
  print("User1 balance: ", fUser1["balance"])
  print("User2 balance: ", fUser2["balance"])
  
  # Pay
  print("Paying...")
  p0 = pay(user1["ibk"], user2["ibk"], 50, "3001")
  p1 = pay(user1["ibk"], user2["ibk"], 50, "3001")
  p2 = pay(user2["ibk"], user1["ibk"], 30, "3002")
  print(p0, p1, p2)
  
  # Find accounts
  print("Finding accounts...")
  fUser1 = find_account(user1["id"], "3001")
  fUser2 = find_account(user2["id"], "3002")
  
  print("User1 balance: ", fUser1["balance"])
  print("User2 balance: ", fUser2["balance"])

def create_user(name, document, port):
  conn = http.client.HTTPConnection("localhost:" + port)
  payload = "{\n  \"name\": \"" + name + "\",\n  \"document\": \"" + document + "\"\n}"
  conn.request("POST", "/api/accounts", payload)
  res = conn.getresponse()
  data = res.read()
  return json.loads(data.decode("utf-8"))

def delete_user(account_id, port):
  conn = http.client.HTTPConnection("localhost:3001")
  conn.request("DELETE", "/api/accounts/" + str(account_id))
  res = conn.getresponse()
  data = res.read()
  return json.loads(data.decode("utf-8"))

def create_deposit(account_id, value, port):
  conn = http.client.HTTPConnection("localhost:"+port)
  payload = "{\n  \"user_id\": " + str(account_id) + ",\n  \"amount\": " + str(value) + "\n}"
  conn.request("POST", "/api/payments/deposit", payload)
  res = conn.getresponse()
  data = res.read()
  return json.loads(data.decode("utf-8"))

def find_account(account_id, port):
  conn = http.client.HTTPConnection("localhost:" + port)
  conn.request("GET", "/api/accounts/" + str(account_id))
  res = conn.getresponse()
  data = res.read()
  return json.loads(data.decode("utf-8"))

def pay(fromUser, toUser, value, port):
  conn = http.client.HTTPConnection("localhost:" + port)
  payload = "{\n  \"from_user_ibk\": \"" + fromUser + "\",\n  \"to_user_ibk\": \"" + toUser + "\",\n  \"amount\": " + str(value) + "\n}"

  conn.request("POST", "/api/payments/pay", payload)

  res = conn.getresponse()
  data = res.read()
  return json.loads(data.decode("utf-8"))

def randomCpf():
  cpf = ""
  for i in range(11):
    cpf += str(random.randint(0, 9))
  return cpf


if __name__ == "__main__":
  run()
  