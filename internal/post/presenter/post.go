package presenter

type PostResponse struct {
	Id 		int
	Title 	string
	PostId 	string
	Url 	string
}

type PostUpdateRequest struct {
	Title	string
}

