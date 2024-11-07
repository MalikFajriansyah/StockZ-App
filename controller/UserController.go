package controller

import (
	"net/http"
	"os"
	"stockz-app/config"
	"stockz-app/model"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

/**
 * Register a new user
 *
 * This function handles the registration of a new user. It expects a JSON payload
 * with the username, email, and password. It returns a JSON response with a
 * success message if the registration is successful, or an error message if
 * there's an issue with the request.
 *
 * Example:
 * ```json
 * {
 *   "username": "johnDoe",
 *   "email": "johndoe@example.com",
 *   "password": "mysecretpassword"
 * }
 * ```
 *
 * @param c echo.Context
 * @return error
 */
func Register(c echo.Context) error {
	var body struct {
		Username string
		Email    string
		Password string
	}

	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Invalid request body",
		})
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Failed generate password",
		})
	}

	newUser := model.User{
		Username: body.Username,
		Email:    body.Email,
		Password: string(hash),
	}
	result := config.DB.Create(&newUser)

	if result.Error != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Failed create user, email already exists",
		})
	}

	return c.JSON(http.StatusCreated, map[string]string{
		"message": "user created",
	})
}

/**
 * Login an existing user
 *
 * This function handles the login of an existing user. It expects a JSON payload
 * with the account (username or email) and password. It returns a JSON response
 * with a success message and a JWT token if the login is successful, or an
 * error message if there's an issue with the request.
 *
 * Example:
 * ```json
 * {
 *   "account": "johnDoe",
 *   "password": "mysecretpassword"
 * }
 * ```
 *
 * @param c echo.Context
 * @return error
 */
func Login(c echo.Context) error {
	var body struct {
		Account  string
		Password string
	}

	if c.Bind(&body) != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Failed to read body",
		})
	}

	var user model.User
	config.DB.First(&user, "username = ? OR email = ?", body.Account, body.Account)

	if user.ID == 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "user not found",
		})
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "password is incorrect",
		})
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":      user.ID,
		"exp":      time.Now().Add(time.Hour * 24 * 30).Unix(),
		"username": user.Username,
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("SECRETKEY")))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Failed create JWT token",
		})
	}

	cookie := new(http.Cookie)
	cookie.Name = "Authorization"
	cookie.Value = tokenString

	cookie.SameSite = http.SameSiteLaxMode
	c.SetCookie(cookie)

	return c.JSON(http.StatusOK, map[string]string{
		"message": "user logged in",
	})
}

/**
 * Home page for logged in users
 *
 * This function returns a JSON response with a personalized message for the
 * logged in user.
 *
 * @param c echo.Context
 * @return error
 */
func Home(c echo.Context) error {
	user := c.Get("user").(model.User)
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Hello " + user.Username,
	})
}

/*
GetProfile returns the profile of a user by username

Example: GET /profile/:username
Response: 200 OK

	{
	  "profile": {
	    "id": 1,
	    "username": "john",
	    "email": "john@example.com",
	    "followers": 10,
	    "followed": 5,
	    "posts": [...]
	  }
	}
*/
func GetProfile(c echo.Context) error {
	var user model.User
	username := c.Param("username")
	if err := config.DB.Preload("Followed").Preload("Follower").Preload("Post").First(&user, username).Error; err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "User doesn't exist")
	}

	type profileResponse struct {
		ID        uint         `json:"id"`
		Username  string       `json:"username"`
		Email     string       `json:"email"`
		Followers int          `json:"followers"`
		Followed  int          `json:"followed"`
		Post      []model.Post `json:"posts"`
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"profile": profileResponse{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			Post:     user.Posts,
		},
	})
}

// SearchUsername searches for users by username
//
// Example: GET /search?username=john
// Response: 202 Accepted
//
//	{
//	  "user": [
//	    {
//	      "id": 1,
//	      "username": "john",
//	      "email": "john@example.com"
//	    },
//	    {
//	      "id": 2,
//	      "username": "johnny",
//	      "email": "johnny@example.com"
//	    }
//	  ]
//	}
func SearchUsername(c echo.Context) error {
	searchName := c.QueryParam("username")
	var users []model.User
	if err := config.DB.Where("username ILIKE ?", "%"+searchName+"%").Find(&users).Error; err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "User with username "+searchName+" not found")
	}

	type userResponse struct {
		ID       uint   `json:"id"`
		Username string `json:"username"`
		Email    string `json:"email"`
	}

	var userResponses []userResponse
	for _, user := range users {
		userResponses = append(userResponses, userResponse{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
		})
	}

	return c.JSON(http.StatusAccepted, map[string]interface{}{
		"user": userResponses,
	})
}
