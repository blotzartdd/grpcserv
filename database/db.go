package database

import (
	"database/sql"
	"fmt"
	"grpcserv/models"
	"log"

	_ "github.com/lib/pq"
)

type DB struct {
	conn *sql.DB
}

func NewDB(host, port, user, password, dbname string) (*DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	conn, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	if err = conn.Ping(); err != nil {
		return nil, err
	}

	db := &DB{conn: conn}
	if err = db.createTables(); err != nil {
		return nil, err
	}

	log.Println("Successfully connected to PostgreSQL database")
	return db, nil
}

func (db *DB) createTables() error {
	query := `
	CREATE TABLE IF NOT EXISTS accounts (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) UNIQUE NOT NULL,
		amount BIGINT NOT NULL DEFAULT 0,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`

	_, err := db.conn.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create accounts table: %v", err)
	}

	log.Println("Database tables created successfully")
	return nil
}

func (db *DB) CreateAccount(name string, amount int64) error {
	query := `INSERT INTO accounts (name, amount) VALUES ($1, $2)`
	_, err := db.conn.Exec(query, name, amount)
	if err != nil {
		return fmt.Errorf("failed to create account: %v", err)
	}
	return nil
}

func (db *DB) DeleteAccount(name string) error {
	query := `DELETE FROM accounts WHERE name = $1`
	result, err := db.conn.Exec(query, name)
	if err != nil {
		return fmt.Errorf("failed to delete account: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("account not found")
	}

	return nil
}

func (db *DB) ChangeAccountName(oldName, newName string) error {
	query := `UPDATE accounts SET name = $1, updated_at = CURRENT_TIMESTAMP WHERE name = $2`
	result, err := db.conn.Exec(query, newName, oldName)
	if err != nil {
		return fmt.Errorf("failed to change account name: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("account not found")
	}

	return nil
}

func (db *DB) ChangeAccountAmount(name string, newAmount int64) error {
	query := `UPDATE accounts SET amount = $1, updated_at = CURRENT_TIMESTAMP WHERE name = $2`
	result, err := db.conn.Exec(query, newAmount, name)
	if err != nil {
		return fmt.Errorf("failed to change account amount: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("account not found")
	}

	return nil
}

func (db *DB) GetAccount(name string) (*models.Account, error) {
	query := `SELECT name, amount FROM accounts WHERE name = $1`
	row := db.conn.QueryRow(query, name)

	var account models.Account
	err := row.Scan(&account.Name, &account.Amount)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("account not found")
		}
		return nil, fmt.Errorf("failed to get account: %v", err)
	}

	return &account, nil
}

func (db *DB) Close() error {
	return db.conn.Close()
}
