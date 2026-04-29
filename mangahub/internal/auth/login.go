package auth

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tcq123-iu/mangahub/pkg/models"
	"github.com/tcq123-iu/mangahub/pkg/utils"
	"golang.org/x/crypto/bcrypt"
)

func Login(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.LoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		var user models.User
		// We use the query variable here to fetch the user by username
		query := "SELECT id, username, password_hash FROM users WHERE username = ?"
		err := db.QueryRow(query, req.Identifier).Scan(&user.ID, &user.Username, &user.PasswordHash)

		if err == sql.ErrNoRows {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Account not found"})
			return
		}

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}

		// Verify the provided password against the stored bcrypt hash
		if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}

		// Generate the JWT token for the session
		token, expiry, err := utils.GenerateToken(user.ID, user.Username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
			return
		}

		// Success! Return the token and user details
		c.JSON(http.StatusOK, models.LoginResponse{
			Token:     token,
			ExpiresAt: expiry,
			Username:  user.Username,
		})
	}
}