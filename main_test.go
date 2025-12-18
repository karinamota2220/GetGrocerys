package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func SetUpRouter() *gin.Engine {
	router := gin.Default()
	return router
}

func TestHomepageHandler(t *testing.T) {
	mockResponse := `{"message":"Welcome to the Tech Company listing API with Golang"}`
	r := SetUpRouter()
	r.GET("/", HomepageHandler)
	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	responseData, _ := io.ReadAll(w.Body)
	assert.Equal(t, mockResponse, string(responseData))
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetGrocerys(t *testing.T) {
	r := SetUpRouter()
	r.GET("/grocerys", getGrocerys)
	req, _ := http.NewRequest("GET", "/grocerys", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var grocerys []Grocery
	json.Unmarshal(w.Body.Bytes(), &grocerys)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotEmpty(t, grocerys)

}
