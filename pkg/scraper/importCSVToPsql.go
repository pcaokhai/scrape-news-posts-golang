package scraper

import (
	"github.com/pcaokhai/scraper/internal/models"
	"github.com/pcaokhai/scraper/pkg/db/postgres"
)

// batch import data from csv to postgres
func ImportCSVToPsql(posts []Post) {
	db := postgres.GetDb()

	var postData =make([]models.Post, 0)
	for _, post := range posts {
		p := models.Post{
			Title: post.Title,
			PostId: post.PostId,
			Url: post.Url,
		}
		postData = append(postData, p)
	
	}

	result := db.Create(postData)
	 if result.Error != nil {
        panic(result.Error)
    }
}