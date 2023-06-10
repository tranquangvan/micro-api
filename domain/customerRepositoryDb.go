package domain

import (
	"database/sql"
	"log"
	"micro-api/errs"
	"time"
)

func NewCustomerRepositoryDb() CustomerRepositoryDb {
	client, err := sql.Open("mysql", "user:password@/dbname")
	if err != nil {
		panic(err)
	}
	// See "Important settings" section.
	client.SetConnMaxLifetime(time.Minute * 3)
	client.SetMaxOpenConns(10)
	client.SetMaxIdleConns(10)
	return CustomerRepositoryDb{client}
}

type CustomerRepositoryDb struct {
	client *sql.DB
}

func (d *CustomerRepositoryDb) FindAll() ([]Customer, error) {
	findAllSql := "select customer_id, name, city, zipcode, date_of_birth, status from customers"

	rows, err := d.client.Query(findAllSql)
	if err != nil {
		log.Println("Error while querying customer table " + err.Error())
		return nil, err
	}
	customers := make([]Customer, 0)

	for rows.Next() {
		var customerRow Customer
		err := rows.Scan(&customerRow)
		if err != nil {
			log.Println("Error while scanning customer table " + err.Error())
			return nil, err
		}
		customers = append(customers, customerRow)
	}
	return customers, err
}
func (d *CustomerRepositoryDb) FindById(id string) (*Customer, *errs.AppError) {
	findCustomerSql := "select customer_id, name, city, zipcode, date_of_birth, status from customers where customer_id = ?"
	row := d.client.QueryRow(findCustomerSql)
	var customer Customer
	err := row.Scan(&customer)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errs.NewNotFoundError("Customer not found")
		} else {
			log.Println("Error while scanning customer " + err.Error())
			return nil, errs.NewUnexprectedError("Unexpected database error")
		}
	}
	return &customer, nil
}
