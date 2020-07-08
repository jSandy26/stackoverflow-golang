package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jsandy26/stackoverflow-golang/models"
)

// CreatePostInput is there to validate the input data
type CreatePostInput struct {
	Title       string   `json:"title" binding:"required"`
	Description string   `json:"description" binding:"required"`
	Tag         []string `json:"tag"`
}

// PostAnswerInput is there to validate the input data
type PostAnswerInput struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
}

// CreatePost creates post
func CreatePost(c *gin.Context) {
	user := c.Keys["currentUser"].(models.User)
	var input CreatePostInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	post := models.Post{
		Title:       input.Title,
		Description: input.Description,
		User:        user,
	}
	models.DB.Save(&post)
	// var tags = make([]models.Tag, 0)
	if len(input.Tag) != 0 {
		// tags = make([]models.Tag, len(input.Tag))
		for _, v := range input.Tag {
			var tag = &models.Tag{Name: v}

			models.DB.FirstOrCreate(&tag, tag)
			models.DB.Model(&post).Association("Tags").Append([]*models.Tag{tag})

		}
	}

	c.JSON(http.StatusOK, gin.H{"data": post})

}

// ListPosts to list all posts related to a tag
func ListPosts(c *gin.Context) {
	tagID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no tag id provided"})
		return
	}
	var posts []models.Post
	var tag models.Tag
	// models.DB.First(&tag, "ID = ?", tagID)
	if err := models.DB.Find(&tag, tagID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "no tag found"})
		return
	}

	// models.DB.Find(&posts)
	// if err := models.DB.Find(&posts).Preload("Tags").Error; err != nil {
	// 	c.JSON(http.StatusNotFound, gin.H{"error": "no posts found"})
	// 	return
	// }
	// models.DB.Model(&posts).Preload("Tags").Find(&posts)
	models.DB.Model(&tag).Preload("Tags").Related(&posts, "Posts")
	c.JSON(http.StatusOK, gin.H{"data": posts})

}

// UpdatePost updates the post description
func UpdatePost(c *gin.Context) {
	user := c.Keys["currentUser"].(models.User)
	postID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no post id provided"})
		return
	}
	var post models.Post
	if err := models.DB.Find(&post, postID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "no posts found"})
		return
	}
	if post.User != user {
		c.JSON(http.StatusNotFound, gin.H{"error": "only author can edit the post"})
		return
	}
	var input CreatePostInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	post.Title = input.Title
	post.Description = input.Description
	post.User = user
	models.DB.Save(&post)
	c.JSON(http.StatusOK, gin.H{"data": post})
}

// CreateAnswer used to post answer to a post
func CreateAnswer(c *gin.Context) {
	user := c.Keys["currentUser"].(models.User)
	postID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no post id provided"})
		return
	}
	var post models.Post
	if err := models.DB.Find(&post, postID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "no posts found"})
		return
	}
	var input PostAnswerInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	answer := models.Post{
		Title:       input.Title,
		Description: input.Description,
		Parent:      &post,
		User:        user,
	}
	models.DB.Save(&answer)
	c.JSON(http.StatusOK, gin.H{"data": answer})
}
