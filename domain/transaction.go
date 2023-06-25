package domain

import "micro-api/dto"

const WITHDRAWAL = "withdrawal"
const DEPOSIT = "deposit"

type Transaction struct {
	TransactionId   string  `db:"transaction_id"`
	AccountId       string  `db:"account_id"`
	Amount          float64 `db:"amount"`
	TransactionType string  `db:"transaction_type"`
	TransactionDate string  `db:"transaction_date"`
}

func (t Transaction) IsWithDrawal() bool {
	if t.TransactionType == WITHDRAWAL {
		return true
	}
	return false
}

func (transaction *Transaction) ToDTO() *dto.TransactionResponse {
	return &dto.TransactionResponse{
		TransactionId:   transaction.TransactionId,
		AccountId:       transaction.AccountId,
		NewBalance:      transaction.Amount,
		TransactionType: transaction.TransactionType,
		TransactionDate: transaction.TransactionDate,
	}
}
