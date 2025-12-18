package db

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

var Pool *pgxpool.Pool

// Initializes PostgreSQL connection
// put all info in one func insead of two
// Connect to PostgreSQL
func InitDB() (DBInterface, error) {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env file")
	}
	// Get connection details from environment variables
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	dbname := os.Getenv("DB_NAME")

	// Create the connection string
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", user, password, host, port, dbname)

	// Initialize connection pool
	pool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Test connection
	if err := pool.Ping(context.Background()); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	fmt.Println("Connected to PostgreSQL!")
	Pool = pool
	return pool, nil
}

// Show all grocerys in the database
func ShowGrocerys() {
	rows, err := Pool.Query(context.Background(), "SELECT * FROM grocerys")
	if err != nil {
		fmt.Println("Error reading grocerys:", err)
		return
	}
	defer rows.Close()

	fmt.Println("Grocerys Table:")
	for rows.Next() {
		var numberitems, groceryitem string
		var price float64

		err := rows.Scan(&numberitems, &groceryitem, &price)
		if err != nil {
			fmt.Println("Scan error:", err)
			return
		}

		fmt.Printf("NumberItems: %s | Item: %s | Price: %.2f\n", numberitems, groceryitem, price)
	}
}
