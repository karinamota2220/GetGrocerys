package models

type Todo struct {
	NumberItems string  `json:"numberitems"`
	GroceryItem string  `json:"groceryitem"`
	Price       float64 `json:"price"`
}

func (Todo) TableName() string {
	return "grocerys"
}
