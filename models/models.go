package models

import "github.com/jinzhu/gorm"

// User model
// type User struct {
// 	ID       uint64 `json:"id"`
// 	Username string `json:"username"`
// 	Password string `json:"password"`
// }

// User model
type User struct {
	gorm.Model
	Username string `json:"username"`
	Password string `json:"password"`
}

// Post model
type Post struct {
	gorm.Model
	Title       string `json:"title"`
	Description string `json:"description"`
	User        User   `gorm:"foreignkey:UserRefer"`
	Tags        []Tag  `gorm:"many2many:post_tags;"`
	Parent      *Post  `gorm:"foreignkey:parent_post"`
}

// Tag model
type Tag struct {
	gorm.Model
	Name  string  `json:"name"`
	Posts []*Post `gorm:"many2many:post_tags"`
}
