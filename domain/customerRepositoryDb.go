package domain

import (
	"micro-api/errs"
	"micro-api/logger"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type CustomerRepository interface {
	FindAll(status string) ([]Customer, *errs.AppError)
	FindById(string) (*Customer, *errs.AppError)
}

func NewCustomerRepositoryDb(dbClient *sqlx.DB) CustomerRepositoryDb {
	return CustomerRepositoryDb{
		client: dbClient,
	}
}

type CustomerRepositoryDb struct {
	client *sqlx.DB
}

func (d *CustomerRepositoryDb) FindAll(status string) ([]Customer, *errs.AppError) {
	customers := make([]Customer, 0)
	var err error
	findAllSql := "select customer_id, name, city, zipcode, date_of_birth, status from customers"
	if status == "" {
		err = d.client.Select(&customers, findAllSql)
	} else {
		findAllSql += " where status = ? "
		err = d.client.Select(&customers, findAllSql, status)
	}

	if err != nil {
		logger.Error("Error while scanning customer " + err.Error())
		return nil, errs.NewUnexprectedError("Unexpected database error")
	}
	return customers, nil
}
func (d *CustomerRepositoryDb) FindById(id string) (*Customer, *errs.AppError) {
	findCustomerSql := "select customer_id, name, city, zipcode, date_of_birth, status from customers where customer_id = ?"
	var customer Customer
	err := d.client.Get(&customer, findCustomerSql, id)
	if err != nil {
		logger.Error("Error while scanning customer " + err.Error())
		return nil, errs.NewUnexprectedError("Unexpected database error")
	}
	return &customer, nil
}
