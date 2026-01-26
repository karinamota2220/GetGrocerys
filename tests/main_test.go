package tests

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"GETALBUMS/db"
	"GETALBUMS/handlers"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	// Load env from project root
	if err := godotenv.Load("../.env"); err != nil {
		fmt.Println("Failed to load .env:", err)
		os.Exit(1)
	}

	// FORCE IPv4 for Windows
	connStr := fmt.Sprintf(
		"postgres://%s:%s@127.0.0.1:%s/%s?sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	pool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		panic(err)
	}

	if err := pool.Ping(context.Background()); err != nil {
		panic(err)
	}

	// Assign to global pool used by handlers
	db.Pool = pool

	code := m.Run()

	pool.Close()
	os.Exit(code)
}

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.Use(handlers.RequestTime()) // middleware

	// register routes for testing
	r.GET("/grocerys/:numberitems", handlers.GetGrocerysByNumberItems)
	r.GET("/grocerys", handlers.GetGrocerys)
	r.POST("/grocerys", handlers.PostGrocerys)
	r.PUT("/grocerys/:numberitems", handlers.UpdateGrocery)
	r.DELETE("/grocerys/:numberitems", handlers.DeleteGrocerys)

	return r
}

func TestGetGrocerysByNumberItems(t *testing.T) {
	ctx := context.Background()

	// Clean table
	_, err := db.Pool.Exec(ctx, "DELETE FROM grocerys")
	require.NoError(t, err)

	// Insert test row
	_, err = db.Pool.Exec(ctx,
		`INSERT INTO grocerys (numberitems, groceryitem, price)
		 VALUES ($1, $2, $3)`,
		"1", "Bread", 3.99,
	)
	require.NoError(t, err)

	r := setupTestRouter()

	req := httptest.NewRequest(http.MethodGet, "/grocerys/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var resp struct {
		Grocery handlers.Grocery `json:"grocery"`
	}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	require.Equal(t, "Bread", resp.Grocery.GroceryItem)
}
