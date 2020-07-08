package middlewares

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/jsandy26/stackoverflow-golang/models"
)

// Login function for user authentication
func Login(c *gin.Context) {
	var u models.User
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "Invalid json provided")
		return
	}

	var user models.User

	if err := models.DB.Where("Username = ? AND Password = ?", u.Username, u.Password).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Please provide valid login details"})
		return
	}

	token, err := CreateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}
	c.JSON(http.StatusOK, token)
}

// CreateToken generates a jwt token
func CreateToken(userid uint) (string, error) {
	var err error
	os.Setenv("ACCESS_SECRET", "jdnfksdmfksd")
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["user_id"] = userid
	atClaims["exp"] = time.Now().Add(time.Minute * 15).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return "", err
	}
	return token, nil
}


// EnforceAuthenticatedMiddleware to authenticate user
func EnforceAuthenticatedMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, exists := c.Get("currentUser")
		if exists && user.(models.User).ID != 0 {
			return
		}
		// err, _ := c.Get("authErr")
		// _ = c.AbortWithError(http.StatusUnauthorized, err.(error))
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user needs to be signed in to access this service"})
		c.Abort()
		return

	}
}

// UserLoaderMiddleware to load user
func UserLoaderMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		bearer := c.Request.Header["Authorization"]
		if bearer != nil {
			fmt.Println(bearer)
			token, err := jwt.Parse(bearer[0], func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signin method %v", token.Header["alg"])
				}
				secret := []byte(os.Getenv("ACCESS_SECRET"))
				return secret, nil
			})

			if err != nil {
				println(err.Error())
				return
			}
			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				userID := uint(claims["user_id"].(float64))
				fmt.Printf("[+] Authenticated request, authenticated user id is %d\n", userID)

				var user models.User
				if userID != 0 {
					models.DB.First(&user, userID)
				}
				fmt.Println("here we are", user)
				c.Set("currentUser", user)
				c.Set("currentUserId", user.ID)
				c.Next()
			} else {

			}

		}
	}
}
