package client

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/AdamWalker95/orders-api/model"
)

type SqlRepo struct {
	Client *sql.DB
}

func (s *SqlRepo) doesCustAlreadyExist(email string) bool {
	var foundUser model.LoginDetails

	query := "SELECT * FROM CUSTOMERS WHERE email = ?;"
	row := s.Client.QueryRow(query, email)

	// If no error customer record already exists
	err := row.Scan(&foundUser.CustomerID, &foundUser.Email, &foundUser.Password)
	if err != nil {
		return false
	}

	return true
}

func (s *SqlRepo) Insert(ctx context.Context, newUser model.LoginDetails) error {

	if err := s.doesCustAlreadyExist(newUser.Email); err != false {
		return fmt.Errorf("Customer's email is already on record")
	}

	query := "INSERT INTO CUSTOMERS(email, password) VALUES (?, ?)"

	stmt, err := s.Client.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("Error %s when preparing SQL statement", err)
	}
	defer stmt.Close()
	res, err := stmt.ExecContext(ctx, newUser.Email, newUser.Password)
	if err != nil {
		return fmt.Errorf("Error %s when inserting row into products table", err)
	}
	_, err = res.RowsAffected()
	if err != nil {
		return fmt.Errorf("Error %s when finding rows affected", err)
	}

	fmt.Printf("Successfully Created Customer Account for %s", newUser.Email)
	return nil
}

var ErrNotExist = errors.New("order does not exist")

func (s *SqlRepo) FindByID(user string) (model.LoginDetails, error) {

	var foundUser model.LoginDetails

	query := "SELECT * FROM CUSTOMERS WHERE email = ?;"
	row := s.Client.QueryRow(query, user)
	err := row.Scan(&foundUser.CustomerID, &foundUser.Email, &foundUser.Password)
	if err != nil {
		return model.LoginDetails{}, fmt.Errorf("Failed to find customer on system: %w", err)
	}

	return foundUser, nil
}

func (s *SqlRepo) Update(ctx context.Context, user model.LoginDetails) error {
	query := "UPDATE CUSTOMERS SET password = ? WHERE email = ?;"
	stmt, err := s.Client.Prepare(query)
	if err != nil {
		return fmt.Errorf("Error preparing update: %w", err)
	}
	defer stmt.Close()

	var res sql.Result
	res, err = stmt.Exec(user.Password, user.Email)
	_, err = res.RowsAffected()
	if err != nil {
		return fmt.Errorf("Error %s when updating row", err)
	}

	fmt.Printf("Successfully Updated Customer Account for %s", user.Email)
	return nil
}

func (s *SqlRepo) DeleteByID(ctx context.Context, user string) error {
	query := "DELETE FROM CUSTOMERS WHERE email = ?;"
	del, err := s.Client.Prepare(query)
	if err != nil {
		return fmt.Errorf("Error preparing delete: %w", err)
	}
	defer del.Close()

	var res sql.Result
	res, err = del.Exec(user)

	_, err = res.RowsAffected()

	if err != nil {
		return fmt.Errorf("Error deleting record: %w", err)
	}

	fmt.Printf("Successfully Deleted Customer Account for %s", user)
	return nil
}
