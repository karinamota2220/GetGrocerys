package handlers

import (
	"GETALBUMS/db"

	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Grocery struct {
	NumberItems string  `json:"numberitems"`
	GroceryItem string  `json:"groceryitem"`
	Price       float64 `json:"price"`
}

func HomepageHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "We are loading your grocery items"})
}

// replalce http with gin
// Description: Logging middleware, request timing
// middleware for time
func RequestTime() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("RequestTime", time.Now().Format(time.RFC3339))
		c.Next() // continue to the next handler
	}
}

// /Get all groceries from
//
//	the database and return them as JSON, along with the request time.
func GetGrocerys(c *gin.Context) {
	RequestTime, _ := c.Get("RequestTime")

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
		"RequestTime": RequestTime,
		"grocerys":    grocerys,
	})
}

// postGrocerys responds with the list of all grocerys as JSON.
// postGrocerys adds an grocery item from JSON received in the request body.

func PostGrocerys(c *gin.Context) {
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

// ///////////// double check what this is doing
// looking for a number item a grocery item
func GetGrocerysByNumberItems(c *gin.Context) {
	numberitems := c.Param("numberitems")
	requestTime, _ := c.Get("RequestTime")

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
func UpdateGrocery(c *gin.Context) {
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

func DeleteGrocerys(c *gin.Context) {
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
