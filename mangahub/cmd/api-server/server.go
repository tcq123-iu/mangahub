package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/tcq123-iu/mangahub/internal/auth"
	"github.com/tcq123-iu/mangahub/pkg/database"
)

func main() {
	// Initialize Database
	db, err := database.InitDB("./data.db")
	if err != nil {
		log.Fatal("Could not connect to database:", err)
	}
	defer db.Close()

	router := gin.Default() 

	authGroup := router.Group("/auth")
	{
		authGroup.POST("/register", auth.Register(db))
        authGroup.POST("/login", auth.Login(db))

	}

	log.Println("API Server running on :8080")
	router.Run(":8080")
}