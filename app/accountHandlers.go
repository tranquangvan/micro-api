package app

import (
	"encoding/json"
	"micro-api/dto"
	"micro-api/service"
	"net/http"

	"github.com/gorilla/mux"
)

type AccountHandlers struct {
	service service.AccountService
}

func (h AccountHandlers) NewAccount(w http.ResponseWriter, r *http.Request) {
	dataParam := mux.Vars(r)
	customerId := dataParam["customer_id"]
	var accountRequest dto.NewAccountRequest
	err := json.NewDecoder(r.Body).Decode(&accountRequest)
	if err != nil {
		writeResponse(w, http.StatusBadRequest, err.Error())
	} else {
		accountRequest.CustomerId = customerId
		account, appErr := h.service.NewAccount(accountRequest)
		if appErr != nil {
			writeResponse(w, appErr.Code, appErr.Message)
		} else {
			writeResponse(w, http.StatusCreated, account)
		}
	}
}

func (h AccountHandlers) MakeTransaction(w http.ResponseWriter, r *http.Request) {
	dataParam := mux.Vars(r)
	customerId := dataParam["customer_id"]
	accountId := dataParam["account_id"]

	var rqTransaction dto.TransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&rqTransaction); err != nil {
		writeResponse(w, http.StatusBadRequest, err.Error())
	} else {
		rqTransaction.CustomerId = customerId
		rqTransaction.AccountId = accountId
		// Make Transaction
		transaction, appError := h.service.MakeTransaction(rqTransaction)
		if appError != nil {
			writeResponse(w, appError.Code, appError.Message)
		} else {
			writeResponse(w, http.StatusCreated, transaction)
		}
	}
}
