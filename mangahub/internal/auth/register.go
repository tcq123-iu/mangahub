package auth

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/tcq123-iu/mangahub/pkg/models"
	"golang.org/x/crypto/bcrypt"
)

func Register(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.RegisterRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			// Syntax error happened here because of citation brackets
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input format"})
			return
		}

		// 1. Hash password using bcrypt
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process password"})
			return
		}

		// Generate a unique ID for the new user
		userID := uuid.New().String()

		// 2. Create user record in SQLite
		query := `INSERT INTO users (id, username, password_hash) VALUES (?, ?, ?)`
		_, err = db.Exec(query, userID, req.Username, string(hashedPassword))
		
		if err != nil {
			// Usually happens if the username already exists in the database
			c.JSON(http.StatusConflict, gin.H{"error": "Username already exists"})
			return
		}

		// 3. Return success confirmation
		c.JSON(http.StatusCreated, gin.H{
			"message": "Account created successfully!",
			"user_id": userID,
		})
	}
}