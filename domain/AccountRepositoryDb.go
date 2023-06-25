package domain

import (
	"micro-api/errs"
	"micro-api/logger"
	"strconv"

	"github.com/jmoiron/sqlx"
)

type AccountRepositoryDb struct {
	client *sqlx.DB
}

func NewAccountRepositoryDb(dbClient *sqlx.DB) AccountRepositoryDb {
	return AccountRepositoryDb{
		client: dbClient,
	}
}

func (d AccountRepositoryDb) Save(a Account) (*Account, *errs.AppError) {
	sqlInsert := "INSERT INTO accounts (customer_id, opening_date, account_type, amount, status) values (?, ?, ?, ?, ?)"

	result, err := d.client.Exec(sqlInsert, a.CustomerId, a.OpeningDate, a.AccountType, a.Amount, a.Status)
	if err != nil {
		logger.Error("Error while creating new account: " + err.Error())
		return nil, errs.NewUnexprectedError("Unexpected error from database")
	}

	accountId, err := result.LastInsertId()
	if err != nil {
		logger.Error("Error while getting last insert id of new account: " + err.Error())
		return nil, errs.NewUnexprectedError("Unexpected error from database")
	}

	a.AccountId = strconv.FormatInt(accountId, 10)
	return &a, nil
}

func (d *AccountRepositoryDb) FindById(id string) (*Account, *errs.AppError) {
	findAccountSql := "select account_id, amount, status, customer_id, account_type from accounts where account_id = ?"
	var account Account
	err := d.client.Get(&account, findAccountSql, id)
	if err != nil {
		logger.Error("Error while scanning account " + err.Error())
		return nil, errs.NewUnexprectedError("Unexpected database error")
	}
	return &account, nil
}

func (d *AccountRepositoryDb) SaveTransaction(transaction Transaction) (*Transaction, *errs.AppError) {
	// starting the database transaction block
	createTransactionSql := "INSERT INTO transactions (account_id , amount, transaction_type, transaction_date) values (?, ?, ?, ?)"
	tx, err := d.client.Begin()
	if err != nil {
		logger.Error("Error while starting a new transaction for bank account transaction" + err.Error())
		return nil, errs.NewUnexprectedError("Unexpected database error")
	}

	// inserting bank account transaction
	result, err := tx.Exec(createTransactionSql, transaction.AccountId, transaction.Amount, transaction.TransactionType, transaction.TransactionDate)
	if err != nil {
		logger.Error("Error while insert a new transaction for bank account transaction" + err.Error())
		return nil, errs.NewUnexprectedError("Unexpected database error")
	}
	if transaction.IsWithDrawal() {
		_, err = tx.Exec(`Update accounts Set amount = amount - ? where account_id =?`, transaction.Amount, transaction.AccountId)
	} else {
		_, err = tx.Exec(`Update accounts Set amount = amount + ? where account_id =?`, transaction.Amount, transaction.AccountId)
	}
	// in case of error RollBack and changes from both the tabels will be reverted
	if err != nil {
		tx.Rollback()
		logger.Error("Error while saving transaction" + err.Error())
		return nil, errs.NewUnexprectedError("Unexpected database error")
	}
	// commit the transaction when all is good
	tx.Commit()
	if err != nil {
		logger.Error("Error while commit transaction" + err.Error())
		return nil, errs.NewUnexprectedError("Unexpected database error")
	}
	// getting the last transactionId from the transaction table
	transactionId, err := result.LastInsertId()
	if err != nil {
		logger.Error("Error while getting transactionId" + err.Error())
		return nil, errs.NewUnexprectedError("Unexpected database error")
	}
	// Getting the latest account information when exec transaction
	account, appErr := d.FindById(transaction.AccountId)
	if appErr != nil {
		logger.Error("Error while getting the latest account information" + appErr.Message)
		return nil, appErr
	}
	transaction.Amount = account.Amount
	transaction.TransactionId = strconv.FormatInt(transactionId, 10)
	return &transaction, nil
}
