package middlewares

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"stockz-app/config"
	"stockz-app/model"

	"cloud.google.com/go/storage"
)

// UploadToFirebase uploads a file to a specified Firebase bucket and returns the URL of the uploaded file
//
// ctx: The context to use for the operation
// client: The Google Cloud Storage client to use for the operation
// bucketName: The name of the Firebase bucket to upload the file to
// fileName: The name to give the uploaded file
// file: The file to upload
//
// Returns: The URL of the uploaded file and an error, if any
func UploadToFirebase(ctx context.Context, client *storage.Client, bucketName, fileName string, file multipart.File) (string, error) {
	bucket := client.Bucket(bucketName)
	obj := bucket.Object(fileName)

	wc := obj.NewWriter(ctx)
	if _, err := io.Copy(wc, file); err != nil {
		return "", fmt.Errorf("failed to write file to GCP : %v", err)
	}
	if err := wc.Close(); err != nil {
		return "", fmt.Errorf("failed to close GCP writer : %v", err)
	}

	if err := obj.ACL().Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
		return "", fmt.Errorf("failed to set file permission : %v", err)
	}

	return fmt.Sprintf("https://storage.googleapis.com/%s/%s", bucketName, fileName), nil
}

// GetPostId retrieves a single post from the database by its ID
//
// id: The ID of the post to retrieve
//
// Returns: The post with the specified ID and an error, if any
func GetPostId(id uint) (*model.Post, error) {
	var post model.Post
	if err := config.DB.First(&post, id).Error; err != nil {
		return nil, err
	}
	return &post, nil
}

// GetCommentInPosts retrieves a list of posts with their associated comments and users
//
// postId: The ID of the post to retrieve. If 0, all posts will be retrieved
//
// Returns: A list of posts with their associated comments and users, and an error, if any
func GetCommentInPosts(postId uint) ([]model.Post, error) {
	var post []model.Post

	db := config.DB.Preload("User").Preload("Comments.User")

	if postId != 0 {
		if err := db.Where("id = ?", postId).Find(&post).Error; err != nil {
			return nil, err
		}
	} else {
		if err := db.Find(&post).Error; err != nil {
			return nil, err
		}
	}
	return post, nil
}
