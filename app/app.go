package app

import (
	"fmt"
	"log"
	"micro-api/domain"
	"micro-api/service"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

func Start() {
	router := mux.NewRouter()
	clientDb := getDbClient()
	customerRepository := domain.NewCustomerRepositoryDb(clientDb)
	accountRepository := domain.NewAccountRepositoryDb(clientDb)
	ch := CustomerHandlers{service.NewCustomerService(customerRepository)}
	ah := AccountHandlers{service.NewAccountService(accountRepository)}

	router.HandleFunc("/customers", ch.getAllCustomers).Methods(http.MethodGet)
	router.HandleFunc("/customers/{customer_id:[0-9]+}", ch.getCustomer).Methods(http.MethodGet)
	router.HandleFunc("/customers/{customer_id:[0-9]+}/account", ah.NewAccount).Methods(http.MethodPost)
	router.HandleFunc("/transactions/{customer_id:[0-9]+}/account/{account_id:[0-9]+}", ah.MakeTransaction).Methods(http.MethodPost)
	fmt.Println("Start server on localhost:8000")
	log.Fatal(http.ListenAndServe(":8000", router))
}

func getDbClient() *sqlx.DB {
	client, err := sqlx.Open("mysql", "mysql:mysqlpassword@(182.20.0.1:3305)/micro-api")
	if err != nil {
		panic(err)
	}
	// See "Important settings" section.
	client.SetConnMaxLifetime(time.Minute * 3)
	client.SetMaxOpenConns(10)
	client.SetMaxIdleConns(10)
	return client
}
