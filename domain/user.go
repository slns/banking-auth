package domain

import "database/sql"

type User struct {
	Username   string         `db:"username"`
	Password   string 		  `db:"password"`
	Role       string         `db:"role"`
	CustomerId sql.NullString `db:"customer_id"`
	CreatedOn  sql.NullString `db:"created_on"`
}