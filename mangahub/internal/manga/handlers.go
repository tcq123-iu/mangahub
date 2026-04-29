package manga

import (
	"database/sql"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/tcq123-iu/mangahub/pkg/models"
	"encoding/json"
    "fmt"
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
func SearchManga(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		title := c.Query("title") // Get ?title=... from URL

		rows, err := db.Query("SELECT id, title, author, genres, status FROM manga WHERE title LIKE ?", "%"+title+"%")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Search failed"})
			return
		}
		defer rows.Close()

		var results []models.Manga
		for rows.Next() {
			var m models.Manga
			rows.Scan(&m.ID, &m.Title, &m.Author, &m.Genres, &m.Status)
			results = append(results, m)
		}

		c.JSON(http.StatusOK, results)
	}
}

// SyncMangaDx fetches 100 manga from the external API and saves them
func SyncMangaDx(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Fetch from MangaDx
		resp, err := http.Get("https://api.mangadex.org/manga?limit=100")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "External API unreachable"})
			return
		}
		defer resp.Body.Close()

		// 2. Parse JSON (Temporary struct to avoid new files)
		var result struct {
			Data []struct {
				ID         string `json:"id"`
				Attributes struct {
					Title  map[string]string `json:"title"`
					Status string            `json:"status"`
				} `json:"attributes"`
			} `json:"data"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse API data"})
			return
		}

		// 3. Store in DB
		count := 0
		for _, d := range result.Data {
			// Using INSERT OR IGNORE so we don't get 409 errors for duplicates
			_, err := db.Exec("INSERT OR IGNORE INTO manga (id, title, author, status) VALUES (?, ?, 'MangaDx API', ?)", 
				d.ID, d.Attributes.Title["en"], d.Attributes.Status)
			if err == nil {
				count++
			}
		}

		c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Successfully synced %d new manga series!", count)})
	}
}