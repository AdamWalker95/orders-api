package model

type CustomerId struct {
	CustomerID int `json:"customer_id"`
}

type LoginDetails struct {
	CustomerID int    `json:"customer_id"`
	Email      string `json:"email"`
	Password   string `json:"password"`
}
