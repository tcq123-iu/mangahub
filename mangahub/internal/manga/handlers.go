package manga

import (
	"database/sql"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/tcq123-iu/mangahub/pkg/models"
)

func GetMangaList(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		rows, err := db.Query("SELECT id, title, author, genres, status FROM manga")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}
		defer rows.Close()

		var mangas []models.Manga
		for rows.Next() {
			var m models.Manga
			rows.Scan(&m.ID, &m.Title, &m.Author, &m.Genres, &m.Status)
			mangas = append(mangas, m)
		}
		c.JSON(http.StatusOK, mangas)
	}
}