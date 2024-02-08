package main

import "time"

type BaseOrder struct {
	ID        int64     `json:"id"`
	Status    string    `json:"status,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	Address   string    `json:"address"`
	Email     string    `json:"email"`
}

type OrderRequest struct {
	BaseOrder
	Products []struct {
		ID       int64 `json:"id"`
		Quantity int64 `json:"quantity"`
	} `json:"products"`
	User *User `json:"user,omitempty"`
}

type Product struct {
	ID          string  `json:"id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Quantity    int64   `json:"quantity"`
	Price       float64 `json:"price"`
}

type Order struct {
	BaseOrder
	Products []Product `json:"products,omitempty"`
}
