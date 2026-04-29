package user

import (
	"database/sql"
	"net/http"
	"github.com/gin-gonic/gin"
)

func AddToLibrary(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString("user_id") // Got this from the Middleware
		mangaID := c.Query("manga_id")

		_, err := db.Exec("INSERT INTO user_progress (user_id, manga_id, current_chapter, status) VALUES (?, ?, 0, 'Reading')", userID, mangaID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not add to library"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Manga added to library"})
	}
}
// internal/user/library.go

// GetLibrary fetches all manga in a user's personal collection
func GetLibrary(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString("user_id")

		query := `
			SELECT m.id, m.title, up.current_chapter, up.status 
			FROM user_progress up 
			JOIN manga m ON up.manga_id = m.id 
			WHERE up.user_id = ?`
		
		rows, err := db.Query(query, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch library"})
			return
		}
		defer rows.Close()

		var library []gin.H
		for rows.Next() {
			var id, title, status string
			var chapter int
			rows.Scan(&id, &title, &chapter, &status)
			library = append(library, gin.H{
				"manga_id": id,
				"title":    title,
				"chapter":  chapter,
				"status":   status,
			})
		}
		c.JSON(http.StatusOK, library)
	}
}

// UpdateProgress updates the chapter count for a specific manga
func UpdateProgress(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString("user_id")
		mangaID := c.Query("manga_id")
		chapter := c.Query("chapter")

		query := "UPDATE user_progress SET current_chapter = ?, updated_at = CURRENT_TIMESTAMP WHERE user_id = ? AND manga_id = ?"
		_, err := db.Exec(query, chapter, userID, mangaID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update progress"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Progress updated successfully"})
	}
}