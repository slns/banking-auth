package domain

import (
	"database/sql"

	"github.com/ashishjuyal/banking-lib/errs"
	"github.com/ashishjuyal/banking-lib/logger"
	"github.com/jmoiron/sqlx"
)

type AuthRepository interface {
	FindBy(username string, password string) (*Login, *errs.AppError)
	SaveUser(user User, customer Customer) (*User, *errs.AppError)
	GenerateAndSaveRefreshTokenToStore(authToken AuthToken) (string, *errs.AppError)
	RefreshTokenExists(refreshToken string) *errs.AppError
}

type AuthRepositoryDb struct {
	client *sqlx.DB
}

func (d AuthRepositoryDb) RefreshTokenExists(refreshToken string) *errs.AppError {
	sqlSelect := "select refresh_token from refresh_token_store where refresh_token = ?"
	var token string
	err := d.client.Get(&token, sqlSelect, refreshToken)
	if err != nil {
		if err == sql.ErrNoRows {
			return errs.NewAuthenticationError("refresh token not registered in the store")
		} else {
			logger.Error("Unexpected database error: " + err.Error())
			return errs.NewUnexpectedError("unexpected database error")
		}
	}
	return nil
}

func (d AuthRepositoryDb) GenerateAndSaveRefreshTokenToStore(authToken AuthToken) (string, *errs.AppError) {
	// generate the refresh token
	var appErr *errs.AppError
	var refreshToken string
	if refreshToken, appErr = authToken.newRefreshToken(); appErr != nil {
		return "", appErr
	}

	// store it in the store
	sqlInsert := "insert into refresh_token_store (refresh_token) values (?)"
	_, err := d.client.Exec(sqlInsert, refreshToken)
	if err != nil {
		logger.Error("unexpected database error: " + err.Error())
		return "", errs.NewUnexpectedError("unexpected database error")
	}
	return refreshToken, nil
}

func (d AuthRepositoryDb) FindBy(username, password string) (*Login, *errs.AppError) {
	var login Login
	sqlVerify := `SELECT username, u.customer_id, role, group_concat(a.account_id) as account_numbers FROM users u
                  LEFT JOIN accounts a ON a.customer_id = u.customer_id
                WHERE username = ? and password = ?
                GROUP BY a.customer_id ORDER BY account_numbers DESC`
	err := d.client.Get(&login, sqlVerify, username, password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errs.NewAuthenticationError("invalid credentials")
		} else {
			logger.Error("Error while verifying login request from database: " + err.Error())
			return nil, errs.NewUnexpectedError("Unexpected database error")
		}
	}
	return &login, nil
}

// func (d AuthRepositoryDb) FindUser(username, password string) (*User, *errs.AppError) {
// 	var user User
// 	sqlVerify := `SELECT * FROM users WHERE username = ? and password = ?`
// 	err := d.client.Get(&user, sqlVerify, username, password)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			return nil, nil
// 			} else {
// 				logger.Error("Error while verifying user request from database: " + err.Error())
// 				return nil, errs.NewUnexpectedError("Unexpected database error")
// 			}
// 		}
// 	return nil, errs.NewAuthorizationError("User already exist")
	
// }

func (d AuthRepositoryDb) SaveUser(u User, c Customer) (*User, *errs.AppError) {
	// starting the database transaction block
	tx, err := d.client.Begin()
	if err != nil {
		logger.Error("Error while starting a new transaction for bank account transaction: " + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected database error")
	}
	// inserting bank account transaction
	result, _ := tx.Exec(`INSERT INTO customers (name, date_of_birth, city, zipcode) 
											values (?, ?, ?, ?)`, c.Name, c.DateofBirth, c.City, c.Zipcode)
	
	// getting the last transaction ID from the transaction table
	customerId, err := result.LastInsertId()
	if err != nil {
		logger.Error("Error while getting the last customer id: " + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected database error")
	}

	_, err = tx.Exec(`INSERT INTO users (username, password, role, customer_id, createde_on) 
											values (?, ?, ?, ?, ?)`, u.Username, u.Password, u.Role, customerId, u.CreatedOn)

	// in case of error Rollback, and changes from both the tables will be reverted
	if err != nil {
		tx.Rollback()
		logger.Error("Error while saving User: " + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected database error")
	}
	
	// commit the transaction when all is good
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		logger.Error("Error while commiting User Customer: " + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected database error")
	}
	
	 return &u, nil
}

func NewAuthRepository(client *sqlx.DB) AuthRepositoryDb {
	return AuthRepositoryDb{client}
}
