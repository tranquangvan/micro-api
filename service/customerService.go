package service

import (
	"micro-api/domain"
	"micro-api/dto"
	"micro-api/errs"
)

type CustomerService interface {
	GetAllCustomer(status string) ([]dto.CustomerResponse, *errs.AppError)
	GetCustomer(string) (*dto.CustomerResponse, *errs.AppError)
}

type DefaultCustomerService struct {
	repo domain.CustomerRepositoryDb
}

func (s DefaultCustomerService) GetAllCustomer(status string) ([]dto.CustomerResponse, *errs.AppError) {
	customers, err := s.repo.FindAll(status)
	if err != nil {
		return nil, err
	}
	customerResponse := []dto.CustomerResponse{}
	for i := range customers {
		response := customers[i].ToDTO()
		customerResponse = append(customerResponse, *response)
	}
	return customerResponse, nil
}

func (s DefaultCustomerService) GetCustomer(id string) (*dto.CustomerResponse, *errs.AppError) {
	customer, err := s.repo.FindById(id)
	if err != nil {
		return nil, err
	}
	response := customer.ToDTO()
	return response, nil
}

func NewCustomerService(repository domain.CustomerRepositoryDb) DefaultCustomerService {
	return DefaultCustomerService{repository}
}
