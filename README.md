### description

- This simple application will scrape posts from (Secret China)[https://www.secretchina.com/] to csv file, postgres, and download scraped articles.
- The CSV is saved to ./data folder
- Download articles will be compressed and saved to ./downloads/[current-date] subfolder
- App will scrape newly-published posts every hours
- Each fetch from database will be cached in redis for improved performance (cached for 30s for the purpose of demo of the project)
- For the purpose of demo, a simple authentication system is setup using JWT, cookie and golang html/template

### How to run

1. Clone the repo and cd to project folder
2. Run the following command to set up postgres and redis
```docker compose up -d```
3. To start the server, run:
```
go mod download
go run cmd/main.go
```

### Test
1. Go to http://localhost:8000 for main website
2. API:
    - GET   Http://localhost:8000/api/v1/post/ get all posts from database
    - POST  Http://localhost:8000/api/v1/post/:postId update a post 
    - DELETE    Http://localhost:8000/api/v1/post/ delete a post

### Techoology used:
1. golang
2. postgres
3. redis
4. clean architecture
5. [obelisk](https://github.com/go-shiori/obelisk)
6. [colly](https://github.com/gocolly/colly)