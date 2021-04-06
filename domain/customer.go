package domain

import (
	"github.com/ashishjuyal/banking-lib/errs"
)

type CustomerRepository interface {
	// status == 1  status == 0 status == ""
	FindAll(status string) ([]Customer, *errs.AppError)
	ById(string) (*Customer, *errs.AppError)
}

type Customer struct {
	Id string `db:"customer_id"`
	Name string `db:"name"`
	City string `db:"city"`
	Zipcode string `db:"zipcode"`
	DateofBirth string `db:"date_of_birth"`
	Status string `db:"status"`
}

// func (c Customer) statusAstext() string {
// 	statusAsText := "active"
// 	if c.Status == "0" {
// 		statusAsText = "inactive"
// 	}
// 	return statusAsText
// }

// func (c Customer) ToDto() dto.CustomerResponse {
	
// 	return dto.CustomerResponse{
// 		Id: 			c.Id,
// 		Name: 			c.Name,
// 		City: 			c.City,
// 		Zipcode: 		c.Zipcode,
// 		DateofBirth: 	c.DateofBirth,
// 		Status: 		c.statusAstext(),
// 	}
// }

