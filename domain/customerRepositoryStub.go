package domain

type CustomerRepositoryStub struct {
	customers []Customer
}

func (s CustomerRepositoryStub) FindAll() ([]Customer, error) {
	return s.customers, nil
}

func NewCustomerRepositoryStub() CustomerRepositoryStub {
	customers := []Customer{
		{"1001", "ABC", "Da Nang", "5555", "1999-9-9", "1"},
		{"1000", "XYZ", "Da Nang", "444", "2000-9-9", "1"},
	}
	return CustomerRepositoryStub{customers}
}
