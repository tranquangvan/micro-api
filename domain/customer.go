package domain

import "micro-api/dto"

type Customer struct {
	Id          string `db:"customer_id"`
	Name        string
	City        string
	Zipcode     string
	DateofBirth string `db:"date_of_birth"`
	Status      string
}

func (customer *Customer) statusAsText() string {
	statusAsText := "active"
	if customer.Status == "0" {
		statusAsText = "inactive"
	}
	return statusAsText
}

func (customer *Customer) ToDTO() *dto.CustomerResponse {
	return &dto.CustomerResponse{
		Id:          customer.Id,
		Name:        customer.Name,
		City:        customer.City,
		Zipcode:     customer.Zipcode,
		DateofBirth: customer.DateofBirth,
		Status:      customer.statusAsText(),
	}
}
