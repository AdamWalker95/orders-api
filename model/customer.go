package model

type Customer struct {
	Username string `json:"username"`
}

type LoginDetails struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CustomerDetails struct {
	CustomerID Customer `json:"customer_id"`
	Orders     []Order  `json:"Orders"`
}
