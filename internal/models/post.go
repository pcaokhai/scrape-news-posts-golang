package models

type Post struct {
	Id 			int 		`json:"id" gorm:"primaryKey" redis:"id"`
	Title 		string 		`json:"title" gorm:"not null" validate:"required" redis:"title"`
	PostId 		string 		`json:"postId" gorm:"unique not null" validate:"required" redis:"postid"`
	Url 		string 		`json:"url" gorm:"not null" validate:"required" redis:"url"`
}

func (post *Post) TableName() string {
	return "post"
}