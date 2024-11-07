package controller

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"stockz-app/config"
	"stockz-app/middlewares"
	"stockz-app/model"
	"strconv"

	"cloud.google.com/go/storage"
	"github.com/labstack/echo/v4"
	"google.golang.org/api/option"
)

func CreatePost(c echo.Context) error {
	//Get data from token context
	user := c.Get("user").(model.User)
	// Parse the user input
	tittle := c.FormValue("tittle")
	description := c.FormValue("description")
	mediaType := c.FormValue("media_type")
	// image file
	mediaFile, err := c.FormFile("media_file")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Failed to read file")
	}

	src, err := mediaFile.Open()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("failed to open image : %v", err))
	}
	defer src.Close()

	//Initialisation to firestore
	ctx := context.Background()
	serviceAccountKeyPath := "./ServiceAccountKey.json"
	storageClient, err := storage.NewClient(ctx, option.WithCredentialsFile(serviceAccountKeyPath))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("failed to create Firebase storage client: %v", err))
	}
	defer storageClient.Close()

	// Bucket name
	bucketName := os.Getenv("BUCKET_NAME")
	fileName := mediaFile.Filename

	// upload to firestore
	imageUrl, err := middlewares.UploadToFirebase(ctx, storageClient, bucketName, fileName, src)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("failed to upload image: %v", err))
	}

	post := model.Post{
		Title:       tittle,
		Description: description,
		MediaType:   mediaType,
		MediaURL:    imageUrl,
		UserID:      user.ID,
	}

	if err := config.DB.Create(&post).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("failed to create post : %v", err))
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "Post created successfully",
		"status":  true,
	})
}

type postResponse struct {
	Username    string      `json:"username"`
	Tittle      string      `json:"tittle"`
	MediaUrl    string      `json:"media_url"`
	Description string      `json:"description"`
	MediaType   string      `json:"media_type"`
	Like        interface{} `json:"like"`
	Comment     interface{} `json:"comment"`
}
type commentResponse struct {
	Content  string `json:"content"`
	Username string `json:"username"`
}

func GetAllPost(c echo.Context) error {
	posts, err := middlewares.GetCommentInPosts(0)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed")
	}

	var response []postResponse

	for _, post := range posts {
		var comments []commentResponse
		for _, comment := range post.Comments {
			comments = append(comments, commentResponse{
				Username: comment.User.Username,
				Content:  comment.Content,
			})
		}

		response = append(response, postResponse{
			Username:    post.User.Username,
			Tittle:      post.Title,
			MediaUrl:    post.MediaURL,
			Description: post.Description,
			MediaType:   post.MediaType,
			Like:        post.Likes,
			Comment:     comments,
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "All posts retrieved successfully",
		"data":    response,
	})
}

func PostById(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID")
	}

	post, err := middlewares.GetPostId(uint(id))
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Post with Id: "+idParam+" not available")
	}

	if err := config.DB.Preload("User").Preload("Comments").First(&post, id).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Post not found")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "All posts retrieved successfully",
		"data": postResponse{
			Username:    post.User.Username,
			Tittle:      post.Title,
			MediaUrl:    post.MediaURL,
			Description: post.Description,
			MediaType:   post.MediaType,
			Like:        post.Likes,
			Comment:     post.Comments,
		},
	})
}

func UpdatePost(c echo.Context) error {
	user := c.Get("user").(model.User)

	// Update image
	mediaFile, err := c.FormFile("media_file")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Failed to read file")
	}
	src, err := mediaFile.Open()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("failed to open image : %v", err))
	}
	defer src.Close()
	//Initialisation to firestore
	ctx := context.Background()
	serviceAccountKeyPath := "./ServiceAccountKey.json"
	storageClient, err := storage.NewClient(ctx, option.WithCredentialsFile(serviceAccountKeyPath))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("failed to create Firebase storage client: %v", err))
	}
	defer storageClient.Close()
	// Bucket name
	bucketName := os.Getenv("BUCKET_NAME")
	fileName := mediaFile.Filename
	// upload to firestore
	imageUrl, err := middlewares.UploadToFirebase(ctx, storageClient, bucketName, fileName, src)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("failed to upload image: %v", err))
	}

	id := c.Param("id")
	var post model.Post
	if err := config.DB.First(&post, id).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Post not found")
	}

	if user.ID != post.UserID {
		return echo.NewHTTPError(http.StatusBadRequest, "You don't have permission to edit this post")
	}

	if err := c.Bind(&post); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid Input")
	}

	updatePost := model.Post{
		Title:       post.Title,
		Description: post.Description,
		MediaURL:    imageUrl,
		MediaType:   post.MediaType,
	}

	if err := config.DB.Model(&post).Updates(updatePost).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update post")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Post updated",
		"data": postResponse{
			Username:    post.User.Username,
			Tittle:      post.Title,
			Description: post.Description,
			MediaUrl:    post.MediaURL,
			MediaType:   post.MediaType,
			Like:        post.Likes,
			Comment:     post.Comments,
		},
	})
}
