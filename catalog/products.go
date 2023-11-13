package main

type Product struct {
	ID          int64   `json:"id"`
	Title       string  `json:"title" form:"title" binding:"required"`
	Description string  `json:"description" form:"description" binding:"required"`
	Category    string  `json:"category" form:"category" binding:"required"`
	Price       float64 `json:"price" form:"price" binding:"required"`
	Stock       int64   `json:"stock" form:"stock" binding:"required"`
	image       []byte
}

var Products = []Product{
	{ID: 1, Title: "Playstation 5", Category: "Gaming", Description: "The latest gaming console from Sony", Price: 499.99, Stock: 15},
}
