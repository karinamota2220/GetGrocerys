package handlers

import (
	"context"
	"net/http"
	"strconv"

	"GETALBUMS/db"
	"GETALBUMS/models"
	"GETALBUMS/utils"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx"
)

// ItemOperations is a struct that holds the database interface
type ItemOperations struct {
	DB db.DBInterface
}

// CreateItem handles creating a new item
func (ops *ItemOperations) CreateItem(c *gin.Context) {
	var item models.Todo
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !utils.ValidateItemName(item.GroceryItem) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid grocery item. Only alphabetic characters are allowed."})
		return
	}

	err := ops.DB.QueryRow(context.Background(),
		"INSERT INTO items (grocery_item, price) VALUES ($1, $2) RETURNING number_items",
		item.GroceryItem, item.Price).Scan(&item.NumberItems)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, item)
}

// GetItems handles fetching all items
func (ops *ItemOperations) GetItems(c *gin.Context) {
	rows, err := ops.DB.Query(context.Background(), "SELECT number_items, grocery_item, price FROM items")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var items []models.Todo
	for rows.Next() {
		var item models.Todo
		err := rows.Scan(&item.NumberItems, &item.GroceryItem, &item.Price)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		items = append(items, item)
	}

	c.JSON(http.StatusOK, items)
}

// GetItem handles fetching a single item by number_items
func (ops *ItemOperations) GetItem(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item ID"})
		return
	}

	var item models.Todo
	err = ops.DB.QueryRow(context.Background(),
		"SELECT number_items, grocery_item, price FROM items WHERE number_items=$1", id).
		Scan(&item.NumberItems, &item.GroceryItem, &item.Price)

	if err == pgx.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, item)
}

// UpdateItem handles updating an existing item
func (ops *ItemOperations) UpdateItem(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item ID"})
		return
	}

	var item models.Todo
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !utils.ValidateItemName(item.GroceryItem) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid grocery item. Only alphabetic characters are allowed."})
		return
	}

	_, err = ops.DB.Exec(context.Background(),
		"UPDATE items SET grocery_item=$1, price=$2 WHERE number_items=$3",
		item.GroceryItem, item.Price, id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	item.NumberItems = strconv.Itoa(id)
	c.JSON(http.StatusOK, item)
}

// DeleteItem handles deleting an item
func (ops *ItemOperations) DeleteItem(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item ID"})
		return
	}

	_, err = ops.DB.Exec(context.Background(), "DELETE FROM items WHERE number_items=$1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
