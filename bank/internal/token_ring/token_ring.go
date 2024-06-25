package token_ring

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/jnaraujo/tec502-inter-bank/bank/internal/config"
	"github.com/jnaraujo/tec502-inter-bank/bank/internal/interbank"
	"github.com/jnaraujo/tec502-inter-bank/bank/internal/storage"
	"github.com/jnaraujo/tec502-inter-bank/bank/internal/token_ring/token"
)

// verifica se o token já esta na rede.
func IsTokenOnRing() bool {
	for _, node := range storage.Ring.List() {
		res, err := http.Get(fmt.Sprintf("http://%s/interbank/token", node.Addr))
		if err != nil {
			continue
		}
		if res.StatusCode == http.StatusOK {
			return true
		}
	}

	return false
}

func BroadcastToken(id interbank.BankId) {
	body, _ := json.Marshal(map[string]any{
		"to": id,
		"ts": time.Now(),
	})
	for _, bank := range storage.Ring.List() {
		http.Post("http://"+bank.Addr+"/interbank/token", "application/json", bytes.NewBuffer(body))
	}
}

func AskBankWithToken() *token.Token {
	for _, bank := range storage.Ring.List() {
		res, err := http.Get("http://" + bank.Addr + "/interbank/token")
		if err != nil {
			continue
		}
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			continue
		}

		data, _ := io.ReadAll(res.Body)
		var tk token.Token
		err = json.Unmarshal(data, &tk)
		if err != nil {
			continue
		}

		return &tk
	}

	return nil
}

func PassToken() {
	// envia a transação para o próximo banco
	nextBank := storage.Ring.Next(config.Env.BankId)
	if nextBank == nil {
		panic("Não tem banco no Token Ring!")
	}

	nextBankId := findNextValidBank(nextBank.Id)
	if nextBankId == nil {
		fmt.Println("Não conseguiu passar o token para um banco valido.")
		return
	}

	BroadcastToken(*nextBankId)
}

func findNextValidBank(id interbank.BankId) *interbank.BankId {
	if id == config.Env.BankId {
		return nil
	}

	bank := storage.Ring.Find(id)
	if bank == nil {
		panic("Banco não encontrado no Token Ring!")
	}

	res, err := http.Get("http://" + bank.Addr + "/interbank/token/ok")
	if err != nil || res.StatusCode != http.StatusOK {
		nextBank := storage.Ring.Next(id)
		if nextBank == nil {
			panic("Não tem banco no Token Ring!")
		}
		return findNextValidBank(nextBank.Id)
	}

	return &bank.Id
}