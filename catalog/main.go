package main

import (
	"bytes"
	"database/sql"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func getProductImage(c *gin.Context) {
	var image []byte
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, err.Error())
		return
	}
	row := DB.QueryRow("SELECT image FROM product WHERE id = ?", id)
	if err := row.Scan(&image); err != nil {
		if err == sql.ErrNoRows {
			c.IndentedJSON(http.StatusNotFound, err.Error())
			return
		}
		c.IndentedJSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.Writer.Write(image)
}

func getProducts(c *gin.Context) {
        qids := c.Query("ids")
        ids := []any{}
        var sb strings.Builder
        sb.WriteString("SELECT id, title, description, category, stock, price FROM product WHERE stock > 0")
        if qids != "" {
                for _, id := range(strings.Split(qids, ",")) {
                        ids = append(ids, id)
                }
                sb.WriteString(" AND id IN (?" + strings.Repeat(",?", len(ids) - 1) + ")")
        }
	products := make([]Product, 0)
	rows, err := DB.Query(sb.String(), ids...)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, "Internal Server Error!")
		return
	}
	defer rows.Close()
	for rows.Next() {
		var product Product
		if err := rows.Scan(&product.ID, &product.Title, &product.Description, &product.Category, &product.Stock, &product.Price); err != nil {
			c.IndentedJSON(http.StatusInternalServerError, "Internal Server Error!")
			return
		}
		products = append(products, product)
	}
	c.IndentedJSON(http.StatusOK, products)
}

func addProduct(c *gin.Context) {
	var newProduct Product

	if err := c.ShouldBind(&newProduct); err != nil {
		c.IndentedJSON(http.StatusBadRequest, err.Error())
		return
	}
	file, _, err := c.Request.FormFile("image")
	if file == nil {
                c.IndentedJSON(http.StatusBadRequest, "Missing 'image' field!")
                return
	}
	defer file.Close()
	if err != nil {
                c.IndentedJSON(http.StatusInternalServerError, err.Error())
		return
	}
	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, file); err != nil {
                c.IndentedJSON(http.StatusInternalServerError, err.Error())
		return
	}
	newProduct.image = buf.Bytes()

	result, err := DB.Exec("INSERT INTO product (title, description, category, stock, price, image) VALUES (?, ?, ?, ?, ?, ?)", newProduct.Title, newProduct.Description, newProduct.Category, newProduct.Stock, newProduct.Price, newProduct.image)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, err.Error())
		return
	}
	id, err := result.LastInsertId()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, err.Error())
		return
	}
	newProduct.ID = id
	c.IndentedJSON(http.StatusCreated, newProduct)
}

func main() {
	connectDB()
	router := gin.Default()

	router.GET("/image/:id", getProductImage)
	router.GET("/", getProducts)
	router.POST("/", addProduct)

	router.Run("0.0.0.0:5555")
}
