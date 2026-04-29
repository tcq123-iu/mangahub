package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/tcq123-iu/mangahub/internal/auth"
	"github.com/tcq123-iu/mangahub/internal/manga"
	"github.com/tcq123-iu/mangahub/internal/user"
	"github.com/tcq123-iu/mangahub/pkg/database"
)

func main() {
	// 1. Initialize Database
	db, err := database.InitDB("./data.db")
	if err != nil {
		log.Fatal("Could not connect to database:", err)
	}
	defer db.Close()

	router := gin.Default()

	// 2. Public Auth Routes
	authGroup := router.Group("/auth")
	{
		authGroup.POST("/register", auth.Register(db))
		authGroup.POST("/login", auth.Login(db))
	}

	// 3. Public Manga Routes

	router.GET("/manga", manga.GetMangaList(db))

	// 4. Protected API Routes (Requires Login)
	api := router.Group("/api")
	api.Use(auth.AuthMiddleware()) 
	{
		api.POST("/library", user.AddToLibrary(db))
		api.GET("/library", user.GetLibrary(db))
		api.PATCH("/library/progress", user.UpdateProgress(db))
		api.DELETE("/library", user.RemoveFromLibrary(db))
	}

	// 5. Start the Server 
	log.Println("API Server running on :8080")
	router.Run(":8080")
}