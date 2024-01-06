package scraper

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"regexp"
	"slices"
	"time"

	"github.com/gocolly/colly"
	"github.com/pcaokhai/scraper/config"
)

const (
	output = "./data/articles.csv"
)

type Post struct {
	Title	string
	Url		string
	PostId	string
	Date 	string
}

// schedule to scrape every hours
func ScheduledScrape(cfg *config.Config) {
	for range time.Tick(cfg.Sync_time * time.Second) {
		Scrape(cfg)
	}
}

func Scrape(cfg *config.Config) {
	posts := importCsv()
	newPosts := make([]Post, 0)

	c := colly.NewCollector()

	// scrape new posts
	c.OnHTML("div#main", func(e *colly.HTMLElement) {
		e.ForEach("div#today ul li", func (i int, h *colly.HTMLElement)  {
			p := saveToPostStruct(h)
			newPosts = append(newPosts, p)
		})

		e.ForEach("div.categ ul li", func (i int, h *colly.HTMLElement)  {
			p := saveToPostStruct(h)
			newPosts = append(newPosts, p)
		})
	})

	c.Visit(cfg.Url)
	// append new published posts to existing scraped posts (sync newly-published posts)
	if len(posts) != 0 {
		newScrapedPosts := newScrapedPosts(posts, newPosts)
		if len(newScrapedPosts) != 0 {
			fmt.Printf("there are new posts, %v\n", len(newScrapedPosts))
			posts = append(posts, newScrapedPosts...)
			go DownloadConcurrently(cfg, newScrapedPosts)
			ImportCSVToPsql(newScrapedPosts)
		} 
	} else {
		// download posts for the first time
		posts = newPosts[:]
		go DownloadConcurrently(cfg, posts)
		ImportCSVToPsql(posts)
	}

	// write to csv
	file, err := os.Create(output)
	if err != nil { 
		log.Fatalln("Failed to create output CSV file", err) 
	} 
	defer file.Close()

	writer := csv.NewWriter(file) 
	header := []string{"title", "date", "postid", "url"}
	writer.Write(header)

	for _, p := range posts {
		record := []string{
			p.Title,
			p.Date,
			p.PostId,
			p.Url,
		}

		writer.Write(record)
	}
	defer writer.Flush()
}

// return only newly-published posts that have not been scraped earlier
func newScrapedPosts(oldPosts []Post, newPosts []Post) []Post {
	if len(oldPosts) == 0 {
		return make([]Post, 0)
	}

	var newScrapedPosts []Post

	for i := range newPosts {
		if !slices.Contains(oldPosts, newPosts[i]) {
			newScrapedPosts = append(newScrapedPosts, newPosts[i])
		}
	}

	return newScrapedPosts
}

// import data from csv
func importCsv() []Post {
	if _, err := os.Stat(output); os.IsNotExist(err) { 
    	return make([]Post, 0)
	}

	f, err := os.Open(output)
    if err != nil {
		return make([]Post, 0)
    }

	defer f.Close()
	csvReader := csv.NewReader(f)
	data, err := csvReader.ReadAll()
	if err != nil {
        log.Fatal(err)
    }

	return readPostsFromCsv(data)
}

func readPostsFromCsv(data [][]string) []Post {
	var postList []Post
	for i, line := range data {
		if i > 0 {
			var record Post
			for j, field := range line {
				if j == 0 {
					record.Title = field
				} else if j == 1 {
					record.Date = field
				} else if j == 2 {
					record.PostId = field
				} else {
					record.Url = field
				}
			}
			postList = append(postList, record)
		}
	}

	return postList
}

func saveToPostStruct(h *colly.HTMLElement) Post {
	p := Post{}
	p.Title = h.ChildText("a")
	url := "https://www.secretchina.com" + h.ChildAttrs("a", "href")[0]
	p.Url = url
	p.PostId = extractPostId(url)
	p.Date = extractDateFromPostURL(url)

	return p
}

// extract post id from url
func extractPostId(url string) string {
	re := regexp.MustCompile(`(\w+).(\w+)$`)
	matches := re.FindStringSubmatch(url)

	return matches[1]
}

// extract date when the article is published
func extractDateFromPostURL(url string) string {
	re := regexp.MustCompile(`https://www.secretchina.com\/news\/gb\/([0-9]*)\/([0-9]*)\/([0-9]*)`)
	matches := re.FindStringSubmatch(url)

	date := fmt.Sprintf("%s-%s-%s", matches[1], matches[2], matches[3])
	return date
}
