package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go-fundraising/db"
	"go-fundraising/routers"
	"log"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Not found file .env")
	}

	db.InitScylla()
	defer db.CloseScylla()

	r := gin.Default()

	routes.InitRouter(r)

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("🚀 Server is running at port", port)
	r.Run(fmt.Sprintf(":%s", port))
}
