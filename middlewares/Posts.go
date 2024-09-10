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

func GetPostId(id uint) (*model.Post, error) {
	var post model.Post
	if err := config.DB.First(&post, id).Error; err != nil {
		return nil, err
	}
	return &post, nil
}

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
