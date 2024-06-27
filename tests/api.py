import http.client
import json
import random

def createAccount(name: str, documents: list, type: str, addr: str) -> dict:
  conn = http.client.HTTPConnection(addr)
  data = {
    "name": name,
    "documents": documents,
    "type": type
  }
  conn.request("POST", "/api/accounts", json.dumps(data))
  res = conn.getresponse()
  data = res.read()
  return json.loads(data.decode("utf-8"))


def deleteUser(accountId: int, addr: str) -> dict:
  conn = http.client.HTTPConnection(addr)
  conn.request("DELETE", "/api/accounts/" + str(accountId))
  res = conn.getresponse()
  data = res.read()
  return json.loads(data.decode("utf-8"))

def createDeposit(accIBK: str, amount: int, addr: str) -> dict:
  conn = http.client.HTTPConnection(addr)
  data = {
    "acc_ibk": accIBK,
    "amount": amount
  }
  conn.request("POST", "/api/payments/deposit", json.dumps(data))
  res = conn.getresponse()
  data = res.read()
  return json.loads(data.decode("utf-8"))

def pay(authorIBK: str, operations: list, addr: str) -> dict:
  conn = http.client.HTTPConnection(addr)
  data = {
    "author": authorIBK,
    "operations": operations
  }
  conn.request("POST", "/api/payments/pay", json.dumps(data))
  res = conn.getresponse()
  data = res.read()
  return json.loads(data.decode("utf-8"))

def findAccount(account_id: int, addr: str) -> dict:
  conn = http.client.HTTPConnection(addr)
  conn.request("GET", "/api/accounts/" + str(account_id))
  res = conn.getresponse()
  data = res.read()
  return json.loads(data.decode("utf-8"))

def findAllTransactions(account_id: int, addr: str) -> dict:
  conn = http.client.HTTPConnection(addr)
  conn.request("GET", "/api/accounts/" + str(account_id) + "/transactions")
  res = conn.getresponse()
  data = res.read()
  return json.loads(data.decode("utf-8"))