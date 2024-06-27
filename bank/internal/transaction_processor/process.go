package transaction_processor

import (
	"github.com/jnaraujo/tec502-inter-bank/bank/internal/models"
	"github.com/jnaraujo/tec502-inter-bank/bank/internal/services"
	"github.com/jnaraujo/tec502-inter-bank/bank/internal/storage"
)

func processTransaction(tr models.Transaction) error {
	for _, op := range tr.Operations {
		err := services.SubCreditFromAccount(op.From, op.Amount)
		if err != nil {
			rollbackOperations(tr)
			return err
		}

		err = services.AddCreditToAccount(op.To, op.Amount)
		if err != nil {
			// como falou na segunda parte, reverte a primeira parte
			services.AddCreditToAccount(op.From, op.Amount)
			rollbackOperations(tr)
			return err
		}

		storage.Transactions.UpdateOperationStatus(tr, op, models.OperationStatusSuccess)
	}

	storage.Transactions.UpdateTransactionStatus(tr, models.TransactionStatusSuccess)

	return nil
}

func rollbackOperations(tr models.Transaction) {
	for _, op := range tr.Operations {
		// so precisa reverter as que tiveram sucesso
		if op.Status == models.OperationStatusSuccess {
			services.SubCreditFromAccount(op.To, op.Amount)
			services.AddCreditToAccount(op.From, op.Amount)
		}

		storage.Transactions.UpdateOperationStatus(tr, op, models.OperationStatusFailed)
	}

	storage.Transactions.UpdateTransactionStatus(tr, models.TransactionStatusFailed)
}
