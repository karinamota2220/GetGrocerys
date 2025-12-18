package main

import (
	"GETALBUMS/db"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func HomepageHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "We are loading your grocery items"})
}

// replalce http with gin
// Description: Logging middleware, request timing
// middleware for time
func requestTime() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("requestTime", time.Now().Format(time.RFC3339))
		c.Next() // continue to the next handler
	}
}

// helloHandler using Gin
// At the momement not being used
func helloHandler(name string) gin.HandlerFunc {
	return func(c *gin.Context) {
		responseText := "<h1>Hello " + name + "</h1>"

		if requestTime, exists := c.Get("requestTime"); exists {
			if str, ok := requestTime.(string); ok {
				responseText += "\n<small>Generated at: " + str + "</small>"
			}
		}

		c.Data(200, "text/html; charset=utf-8", []byte(responseText))
	}
}

// todo list this is in memory
// type Todo struct {
// 	NumberItems string  `json:"numberitems"`
// 	GroceryItem string  `json:"groceryitem"`
// 	Price       float64 `json:"price"`
// }

// var grocerys = []Todo{
// 	{NumberItems: "1", GroceryItem: "Eggs", Price: 7.99},
// 	{NumberItems: "2", GroceryItem: "Fair Life Milk", Price: 10.99},
// 	{NumberItems: "3", GroceryItem: "Snapple", Price: 6.99},
// }

type Grocery struct {
	NumberItems string  `json:"numberitems"`
	GroceryItem string  `json:"groceryitem"`
	Price       float64 `json:"price"`
}

// album represents data about a record album.
// type album struct {
// 	ID     string  `json:"id"`
// 	Title  string  `json:"title"`
// 	Artist string  `json:"artist"`
// 	Price  float64 `json:"price"`
// }

// albums slice to seed record album data.
// var albums = []album{
// 	{ID: "1", Title: "All My Demons Greeting Me As A Friend", Artist: "AURORA", Price: 12.99},
// 	{ID: "2", Title: "Beautiful Chaos", Artist: "Katseye", Price: 24.99},
// 	{ID: "3", Title: "The Happy Star", Artist: "Lexie Liu", Price: 39.99},
// }

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
	// router.GET("/albums", requestTime(), getAlbums)
	// router.GET("/albums/:id", requestTime(), getAlbumByID)
	// router.POST("/albums", requestTime(), postAlbums)
	router.GET("/", HomepageHandler)
	router.GET("/karina", requestTime(), helloHandler("Karina"))
	router.GET("/grocerys", requestTime(), getGrocerys)
	router.GET("/grocerys/:numberitems", requestTime(), getGrocerysByNumberItems)
	router.POST("/grocerys", requestTime(), postGrocerys)
	router.PUT("/grocerys", requestTime(), updateGrocery)
	router.DELETE("/grocerys/:numberitems", requestTime(), deleteGrocerys)

	router.Run(":8081")
}

// getAlbums responds with the list of all albums as JSON.
// add requestTime
// func getAlbums(c *gin.Context) {
// 	requestTime, _ := c.Get("requestTime")
// 	c.IndentedJSON(http.StatusOK, gin.H{
// 		"requestTime": requestTime,
// 		"albums":      albums,
// 	})
// }

//////////////

// for test

// /Get all groceries from
//
//	the database and return them as JSON, along with the request time.
func getGrocerys(c *gin.Context) {
	requestTime, _ := c.Get("requestTime")

	rows, err := db.Pool.Query(context.Background(),
		"SELECT numberitems, groceryitem, price FROM grocerys",
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	grocerys := []Grocery{}
	for rows.Next() {
		var g Grocery
		if err := rows.Scan(&g.NumberItems, &g.GroceryItem, &g.Price); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		grocerys = append(grocerys, g)
	}

	c.JSON(http.StatusOK, gin.H{
		"requestTime": requestTime,
		"grocerys":    grocerys,
	})
}

// postAlbums adds an album from JSON received in the request body.
// func postAlbums(c *gin.Context) {
// 	var newAlbum album

// Call BindJSON to bind the received JSON to
// newAlbum.
// if err := c.BindJSON(&newAlbum); err != nil {
// 	return
// }

// Add the new album to the slice.
// 	albums = append(albums, newAlbum)
// 	c.IndentedJSON(http.StatusCreated, newAlbum)
// }

////////////

func postGrocerys(c *gin.Context) {
	var newG Grocery
	if err := c.BindJSON(&newG); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Postgres table
	_, err := db.Pool.Exec(context.Background(),
		"INSERT INTO grocerys (numberitems, groceryitem, price) VALUES ($1, $2, $3)",
		newG.NumberItems, newG.GroceryItem, newG.Price,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, newG)
}

/////////////// double check what this is doing
// looking for a number item a grocery item

func getGrocerysByNumberItems(c *gin.Context) {
	numberitems := c.Param("numberitems")
	requestTime, _ := c.Get("requestTime")

	row := db.Pool.QueryRow(context.Background(),
		"SELECT numberitems, groceryitem, price FROM grocerys WHERE numberitems=$1",
		numberitems,
	)

	var g Grocery
	if err := row.Scan(&g.NumberItems, &g.GroceryItem, &g.Price); err != nil {
		if err.Error() == "no rows in result set" {
			c.JSON(http.StatusNotFound, gin.H{"message": "Grocery item not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"requestTime": requestTime,
		"grocery":     g,
	})
}

// /////////////
func updateGrocery(c *gin.Context) {
	numberItems := c.Param("numberitems")

	var updated Grocery
	if err := c.BindJSON(&updated); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Update logic here…
	c.JSON(200, gin.H{
		"message":      "Grocery updated",
		"numberitems":  numberItems,
		"updated_data": updated,
	})
}

///////////////

func deleteGrocerys(c *gin.Context) {
	numberitems := c.Param("numberitems")

	result, err := db.Pool.Exec(
		context.Background(),
		"DELETE FROM grocerys WHERE numberitems=$1",
		numberitems,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "Grocery item not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "Grocery item deleted",
		"numberitems": numberitems,
	})
}
