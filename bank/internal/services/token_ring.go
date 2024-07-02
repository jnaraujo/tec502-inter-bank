package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/jnaraujo/tec502-inter-bank/bank/internal/config"
	"github.com/jnaraujo/tec502-inter-bank/bank/internal/interbank"
	"github.com/jnaraujo/tec502-inter-bank/bank/internal/storage"
	"github.com/jnaraujo/tec502-inter-bank/bank/internal/token"
)

func SetupTokenRing() {
	if storage.Ring.FindBankWithLowestId().Id == config.Env.BankId {
		slog.Info("I'm the bank with the lowest id")
		// verifica se o token já esta na rede.
		if !IsTokenOnRing() {
			// se não estiver, cria o token
			BroadcastToken(config.Env.BankId)
		}
	}

	bank := AskBankWithToken()
	if bank != nil && bank.Owner == config.Env.BankId {
		storage.Token.Set(*bank)
		return
	}

	// define o token como o banco com menor id
	storage.Token.Set(token.Token{
		Owner: storage.Ring.FindBankWithLowestId().Id,
		Ts:    time.Now(),
	})
}

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
	client := &http.Client{}
	for _, bank := range storage.Ring.List() {
		req, _ := http.NewRequest("PUT", "http://"+bank.Addr+"/interbank/token", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		client.Do(req)
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
		slog.Info("Não conseguiu passar o token para um banco valido. Mantém o token localmente.")
		BroadcastToken(config.Env.BankId) // faz o broadcast do token para os outros bancos - para garantir que o token não se perca
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
		slog.Error("Banco não respondeu ao token", "bank", bank.Id, "error", err)

		nextBank := storage.Ring.Next(id)
		if nextBank == nil {
			panic("Não tem banco no Token Ring!")
		}
		return findNextValidBank(nextBank.Id)
	}

	return &bank.Id
}
