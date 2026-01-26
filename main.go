package main

import (
	"GETALBUMS/db"
	"GETALBUMS/handlers"

	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	// connect to database
	pool, err := db.InitDB()
	if err != nil {
		fmt.Println("Database connection failed:", err)
		return
	}
	defer pool.Close()
	fmt.Println("Database connection established")

	// Optional: show grocerys in console (for testing)
	db.ShowGrocerys()

	//  In-memory todo list API with POST/GET endpoints
	router := gin.Default()
	router.GET("/", handlers.HomepageHandler)
	router.GET("/grocerys", handlers.RequestTime(), handlers.GetGrocerys)
	router.GET("/grocerys/:numberitems", handlers.RequestTime(), handlers.GetGrocerysByNumberItems)
	router.POST("/grocerys", handlers.RequestTime(), handlers.PostGrocerys)
	router.PUT("/grocerys", handlers.RequestTime(), handlers.UpdateGrocery)
	router.DELETE("/grocerys/:numberitems", handlers.RequestTime(), handlers.DeleteGrocerys)

	router.Run(":8081")
}
