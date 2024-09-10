package controller

import (
	"net/http"
	"stockz-app/config"
	"stockz-app/model"

	"github.com/labstack/echo/v4"
)

func CreateComment(c echo.Context) error {
	user := c.Get("user").(model.User)

	idPost := c.Param("id")
	var post model.Post
	if err := config.DB.First(&post, idPost).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Post not found")
	}
	if err := c.Bind(&post); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid input")
	}

	comment := model.Comment{
		PostID:  post.ID,
		UserID:  user.ID,
		Content: c.FormValue("content"),
	}

	if err := config.DB.Create(&comment).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create comment")
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "Comment created successfully",
	})
}
