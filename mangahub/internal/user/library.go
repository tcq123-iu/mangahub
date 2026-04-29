package user

import (
	"database/sql"
	"net/http"
	"github.com/gin-gonic/gin"
)

func AddToLibrary(db *sql.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        userID := c.GetString("user_id")
        mangaID := c.Query("manga_id")

        // 1. Basic Validation: Ensure mangaID isn't empty
        if mangaID == "" {
            c.JSON(http.StatusBadRequest, gin.H{"error": "manga_id is required"})
            return
        }

        // 2. Data Integrity Validation: Check if Manga actually exists in our DB
        var exists bool
        checkQuery := "SELECT EXISTS(SELECT 1 FROM manga WHERE id = ?)"
        err := db.QueryRow(checkQuery, mangaID).Scan(&exists)
        
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Database check failed"})
            return
        }

        if !exists {
            c.JSON(http.StatusNotFound, gin.H{"error": "Manga not found in our database"})
            return
        }

        // 3. Logic Execution: If we got here, the manga is real. Now add it.
        _, err = db.Exec("INSERT INTO user_progress (user_id, manga_id, current_chapter, status) VALUES (?, ?, 0, 'Reading')", userID, mangaID)
        if err != nil {
            // Check if it's a duplicate entry (Unique constraint on PRIMARY KEY(user_id, manga_id))
            c.JSON(http.StatusConflict, gin.H{"error": "Manga already in your library"})
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
// internal/user/library.go

func RemoveFromLibrary(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString("user_id")
		mangaID := c.Query("manga_id")

		if mangaID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "manga_id is required"})
			return
		}

		result, err := db.Exec("DELETE FROM user_progress WHERE user_id = ? AND manga_id = ?", userID, mangaID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove manga"})
			return
		}

		// Check if a row was actually deleted
		rowsAffected, _ := result.RowsAffected()
		if rowsAffected == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Manga not found in your library"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Manga removed from library"})
	}
}