package scraper

import (
	"compress/gzip"
	"context"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"

	"github.com/go-shiori/obelisk"
	"github.com/pcaokhai/scraper/config"
	"github.com/pcaokhai/scraper/pkg/utils"
)

// data is divided into chunks, num of chunks == num of threads
// download data concurrently in separate threads
func DownloadConcurrently(cfg *config.Config, posts []Post) {
	var wg sync.WaitGroup

	// posts are downloaded concurrently only if num of posts >= num of threads set in the config.json
	if len(posts) < cfg.Num_processes {
		wg.Add(1)
		go downloadPagesAndCompress(posts, &wg)
		return
	} else {
		chunks := chunkSlice(posts, cfg.Num_processes)
		for i := 0; i < cfg.Num_processes; i++ {
			wg.Add(1)
			go downloadPagesAndCompress(chunks[i], &wg)
		}
	}	

	wg.Wait()
	fmt.Println("Downloads finished")
}

func downloadPagesAndCompress(posts []Post, wg *sync.WaitGroup) {
	defer wg.Done()
	requests := convertPostRequest(posts)

	arc := obelisk.Archiver{EnableLog: true}
	arc.Validate()
	currentDate := utils.GetCurrentDate()
	dir := "./downloads/" +currentDate

	if _, err := os.Stat(dir); os.IsNotExist(err) { 
    	os.MkdirAll(dir, os.ModePerm)
	}

	for _, request := range requests {
		result, _, err := arc.Archive(context.Background(), request)
		checkError(err)

		filenameHTML := path.Base(request.URL)
		filename := strings.Replace(filenameHTML, "html", "html.gz", -1)
		
		f, err := os.Create(filepath.Join(dir, filename))
		checkError(err)
		defer f.Close()

		gz := gzip.NewWriter(f)
		gz.Write(result)
		gz.Close()
	}
}

func convertPostRequest(posts []Post) []obelisk.Request {
	requests := make([]obelisk.Request, 0)

	for _, post := range posts {
		request := obelisk.Request{
			URL: post.Url,
		}
		requests = append(requests, request)
	}

	return requests
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

// split a slice into chunks of data
func chunkSlice(posts []Post, numOfChunks int) [][]Post {
	var chunks [][]Post
	for i := 0; i < numOfChunks; i++ {
		min := (i * len(posts) / numOfChunks)
		max := ((i + 1) * len(posts)) /numOfChunks

		chunks = append(chunks, posts[min:max])
	}

	return chunks
}