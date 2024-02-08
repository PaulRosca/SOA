package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/IBM/sarama"
	"github.com/gin-gonic/gin"
)

type Notification struct {
	ProductID int64 `json:"productID"`
	Quantity  int64 `json:"quantity"`
}

type User struct {
	ID    int64  `json:"id"`
	Email string `json:"email"`
	Type  string `json:"type"`
}

type Body struct {
	User *User `json:"user,omitempty"`
}

func getOrders(c *gin.Context) {
	var body Body
	data, err := c.GetRawData()
	if err == nil {
		if err := json.Unmarshal(data, &body); err != nil {
			c.IndentedJSON(http.StatusBadRequest, err.Error())
			return

		}
	}
	ids := []any{}
	var ordersQuery strings.Builder
	ordersQuery.WriteString("SELECT id, status, timestamp, address, email FROM orders")
	if body.User == nil {
		c.IndentedJSON(http.StatusUnauthorized, "User not logged in!")
		return
	}
	if body.User.Type != "admin" {
		ordersQuery.WriteString(" WHERE user_id = ?")
		ids = append(ids, body.User.ID)
	}
	orders := make([]Order, 0)
	rows, err := DB.Query(ordersQuery.String(), ids...)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, err.Error())
		return
	}
	defer rows.Close()
	for rows.Next() {
		var order Order
		if err := rows.Scan(&order.ID, &order.Status, &order.Timestamp, &order.Address, &order.Email); err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err.Error())
			return
		}
		itemRows, err := DB.Query("SELECT order_item.product_id, quantity, title, description, price FROM order_item INNER JOIN catalog.product ON catalog.product.id = order_item.product_id  WHERE order_id = ?", order.ID)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err.Error())
			return
		}
		products := make([]Product, 0)
		for itemRows.Next() {
			var product Product
			if err := itemRows.Scan(&product.ID, &product.Quantity, &product.Title, &product.Description, &product.Price); err != nil {
				c.IndentedJSON(http.StatusInternalServerError, err.Error())
				return
			}
			products = append(products, product)
		}
		order.Products = products
		orders = append(orders, order)
	}
	c.IndentedJSON(http.StatusOK, orders)
}

func addOrder(producer sarama.SyncProducer) gin.HandlerFunc {
	return func(c *gin.Context) {
		var orderReq OrderRequest
		if err := c.ShouldBind(&orderReq); err != nil {
			c.IndentedJSON(http.StatusBadRequest, err.Error())
			return
		}

		fmt.Println(orderReq)
		result, err := DB.Exec("INSERT INTO orders (status, address, email) VALUES (?, ?, ?)", "processing", orderReq.Address, orderReq.Email)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err.Error())
			return
		}
		id, err := result.LastInsertId()
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err.Error())
			return
		}

		if orderReq.User != nil {
			result, err = DB.Exec("UPDATE orders SET user_id = ? WHERE id = ?", orderReq.User.ID, id)
			if err != nil {
				c.IndentedJSON(http.StatusInternalServerError, err.Error())
				return
			}
		}

		orderReq.ID = id
		for _, prod := range orderReq.Products {
			_, err := DB.Exec("INSERT INTO order_item (order_id, product_id, quantity) VALUES (?, ?, ?)", orderReq.ID, prod.ID, prod.Quantity)
			if err != nil {
				c.IndentedJSON(http.StatusInternalServerError, err.Error())
				return
			}
			if err := sendKafkaMessage(producer, c, orderReq.ID, prod.ID, prod.Quantity); err != nil {
				c.IndentedJSON(http.StatusInternalServerError, err.Error())
				return
			}
		}
		c.IndentedJSON(http.StatusCreated, orderReq)
	}
}

func main() {
	producer, err := setupProducer()
	if err != nil {
		log.Fatalf("failed to initialize producer: %v", err)
	}
	defer producer.Close()
	connectDB()
	router := gin.Default()

	router.GET("/", getOrders)
	router.POST("/", addOrder(producer))

	router.Run("0.0.0.0:6666")
}
