package application

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/AdamWalker95/orders-api/model"
)

// Gets Database name from config
func (a *App) DatabaseLookup(database string) string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s", a.config.User, a.config.Password, a.config.MySqlAddress, database)
}

// Starts and Configures SQL Database
func (a *App) StartSqlDatabase() (*sql.DB, error) {

	// Starts up MySQL
	db, err := sql.Open("mysql", a.DatabaseLookup(""))
	if err != nil {
		fmt.Printf("Error %s when opening DB\n", err)
		return nil, fmt.Errorf("Failed to open MySQL Database: %w", err)
	}
	//defer db.Close()

	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()

	// Creates database if database doesn't already
	res, err := db.ExecContext(ctx, "CREATE DATABASE IF NOT EXISTS "+a.config.MySqlDatabaseName)
	if err != nil {
		fmt.Printf("Error %s when creating DB\n", err)
		return nil, fmt.Errorf("Failed to create MySQL Database %s: %w", a.config.MySqlDatabaseName, err)
	}

	// Checks the number of rows affected by running call to create database
	_, err = res.RowsAffected()
	if err != nil {
		fmt.Printf("Error %s when fetching rows", err)
		return nil, fmt.Errorf("Failed to fetch MySQL Database Rows: %w", err)
	}

	// Opens database
	db, err = sql.Open("mysql", a.DatabaseLookup(a.config.MySqlDatabaseName))
	if err != nil {
		fmt.Printf("Error %s when opening DB\n", err)
		return nil, fmt.Errorf("Failed to open MySQL Database %s: %w", a.config.MySqlDatabaseName, err)
	}

	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(20)
	db.SetConnMaxLifetime(time.Minute * 5)

	ctx, cancelfunc = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()

	// Confirms database can be contacted
	err = db.PingContext(ctx)
	if err != nil {
		fmt.Printf("Errors %s pinging DB", err)
		return nil, fmt.Errorf("Failed to ping MySQL Database %s: %w", a.config.MySqlDatabaseName, err)
	}
	fmt.Printf("Connected to DB %s successfully\n", a.config.MySqlDatabaseName)

	return db, nil
}

// Creates tables for database
func createTables(db *sql.DB) error {

	err := createCustomersTable(db)
	if err != nil {
		return err
	}

	return createOrdersTable(db)
}

func createCustomersTable(db *sql.DB) error {

	query := `CREATE TABLE IF NOT EXISTS CUSTOMERS(customer_id int primary key auto_increment, email text, password text)`
	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	res, err := db.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("Error %s when creating product table", err)
	}
	_, err = res.RowsAffected()
	if err != nil {
		return fmt.Errorf("Error %s when getting rows affected", err)
	}

	return nil
}

func createOrdersTable(db *sql.DB) error {

	query := `CREATE TABLE IF NOT EXISTS ORDERS(order_id int primary key auto_increment, customer_id int, created_at datetime, shipped_at datetime, completed_at datetime)`
	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	res, err := db.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("Error %s when creating product table", err)
	}
	_, err = res.RowsAffected()
	if err != nil {
		return fmt.Errorf("Error %s when getting rows affected", err)
	}

	if err := Insert(db); err != nil {
		return fmt.Errorf("Error %s when creating/finding records", err)
	}

	return nil
}

func Insert(db *sql.DB) error {

	/*ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()

	for i := 1; i < 5; i++ {
		query := "INSERT INTO ORDERS(customer_id) VALUES (?)"

		stmt, err := db.PrepareContext(ctx, query)
		if err != nil {
			return fmt.Errorf("Error %s when preparing SQL statement", err)
		}
		defer stmt.Close()
		res, err := stmt.ExecContext(ctx, 13)
		if err != nil {
			return fmt.Errorf("Error %s when inserting row into products table", err)
		}
		_, err = res.RowsAffected()
		if err != nil {
			return fmt.Errorf("Error %s when finding rows affected", err)
		}

		fmt.Printf("Successfully Created Order %d", i)
	}*/

	if err := FindByID(db); err != nil {
		return fmt.Errorf("Error %s when finding records", err)
	}
	return nil
}

func FindByID(db *sql.DB) error {

	for i := 1; i < 5; i++ {
		var foundOrder model.Order

		query := "SELECT * FROM ORDERS WHERE order_id = ?;"
		row := db.QueryRow(query, i)
		err := row.Scan(&foundOrder.OrderID, &foundOrder.CustomerID, &foundOrder.CreatedAt, &foundOrder.ShippedAt, &foundOrder.CompletedAt)
		if err != nil {
			return fmt.Errorf("Failed to find order on system: %w", err)
		}

		fmt.Println(foundOrder)
	}

	return nil
}
