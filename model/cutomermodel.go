package model

type CustomerModel struct {
	CustomerID  int64
	FirstName   string
	LastName    string
	Email       string
	PhoneNumber string
	Address     string
	Password    string `json:"-"`
	CreatedAt   string
	UpdatedAt   string
}
