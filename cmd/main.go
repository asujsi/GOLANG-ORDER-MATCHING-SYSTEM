package main

import (
	"golang-order-matching-system/internal/api"
	"golang-order-matching-system/internal/db"

	"github.com/gin-gonic/gin"
)

func main() {
	db.InitDB()

	r := gin.Default()

	api.RegisterRoutes(r)

	r.Run(":8080")
}
