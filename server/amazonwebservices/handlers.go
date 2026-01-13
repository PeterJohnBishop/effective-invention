package amazonwebservices

import (
	"effective-invention/server/amazonwebservices/auth"
	"effective-invention/server/amazonwebservices/database"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt"
)

func HandleFileUpload(client *s3.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Auth header not found"})
			return
		}
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == authHeader {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "API token error"})
			return
		}
		claims := auth.ParseAccessToken(token)
		if claims == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error parsing auth token"})
			return
		}

		file, header, err := c.Request.FormFile("file") // "file" is the key in the form-data
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get file from request"})
			return
		}
		defer file.Close()

		url, err := UploadFile(client, header.Filename, file)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "File uploaded successfully",
			"url":     url,
		})
	}
}

func HandleFileDOwnloadLink(client *s3.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Auth header not found"})
			return
		}
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == authHeader {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "API token error"})
			return
		}
		claims := auth.ParseAccessToken(token)
		if claims == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error parsing auth token"})
			return
		}

		filename := c.Param("filename")
		if filename == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Filename is required"})
			return
		}

		url, err := DownloadFile(client, filename)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "File download link generated successfully",
			"expires": "5 minutes",
			"url":     url,
		})
	}
}

func HandleFileDownloadStream(client *s3.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Auth header not found"})
			return
		}
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == authHeader {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "API token error"})
			return
		}
		claims := auth.ParseAccessToken(token)
		if claims == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error parsing auth token"})
			return
		}

		filename := c.Param("filename")

		err := StreamDownloadFile(c, client, filename)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
	}
}

func HandleUserCreation(client *dynamodb.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user database.User
		err := json.NewDecoder(c.Request.Body).Decode(&user)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		defer c.Request.Body.Close()

		id, err := uuid.NewV1()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		email := strings.ToLower(user.Email)

		userId := fmt.Sprintf("USER_%s", id)

		hashedPassword, err := auth.HashedPassword(user.Password)
		if err != nil {
			return
		}

		newUser := map[string]types.AttributeValue{
			"id":       &types.AttributeValueMemberS{Value: userId},
			"name":     &types.AttributeValueMemberS{Value: user.Name},
			"email":    &types.AttributeValueMemberS{Value: email},
			"password": &types.AttributeValueMemberS{Value: hashedPassword},
		}

		err = database.CreateUser(client, "users", newUser)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		response := map[string]interface{}{
			"message": "User succesfully created",
			"user":    newUser,
		}

		c.JSON(http.StatusOK, response)
	}
}

func HandleAuthentication(client *dynamodb.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		type LoginRequest struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		var req LoginRequest

		if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		defer c.Request.Body.Close()

		user, err := database.GetUserByEmail(client, "users", req.Email)
		if err != nil || user == nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "User not found or database error",
			})
			return
		}

		pass := auth.CheckPasswordHash(req.Password, user.Password)
		if !pass {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}

		userClaims := auth.UserClaims{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
			StandardClaims: jwt.StandardClaims{
				IssuedAt:  time.Now().Unix(),
				ExpiresAt: time.Now().Add(time.Minute * 15).Unix(),
			},
		}

		token, err := auth.NewAccessToken(userClaims)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Error generating access token",
			})
			return
		}

		refreshToken, err := auth.NewRefreshToken(userClaims.StandardClaims)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Error generating refresh token",
			})
			return
		}

		response := map[string]interface{}{
			"message":       "Login Success",
			"token":         token,
			"refresh_token": refreshToken,
			"user":          user,
		}

		c.JSON(http.StatusOK, response)

	}
}

func HandleGetAllUsers(client *dynamodb.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Auth header not found"})
			return
		}
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == authHeader {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "API token error"})
			return
		}
		claims := auth.ParseAccessToken(token)
		if claims == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error parsing auth token"})
			return
		}
		resp, err := database.GetAllUsers(client, "users")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		var users []database.User

		for _, item := range resp {
			var user database.User
			err = attributevalue.UnmarshalMap(item, &user)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			users = append(users, user)
		}

		response := map[string]interface{}{
			"message": "Users Found!",
			"users":   users,
		}

		c.JSON(http.StatusOK, response)
	}
}

func HandleGetUserById(client *dynamodb.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Auth header not found"})
			return
		}
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == authHeader {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "API token error"})
			return
		}
		claims := auth.ParseAccessToken(token)
		if claims == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error parsing auth token"})
			return
		}

		resp, err := database.GetUserById(client, "users", id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		var user database.User
		err = attributevalue.UnmarshalMap(resp, &user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		response := map[string]interface{}{
			"message": "User Found!",
			"user":    user,
		}

		c.JSON(http.StatusOK, response)

	}
}

func HandleUpdateUser(client *dynamodb.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Auth header not found"})
			return
		}
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == authHeader {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "API token error"})
			return
		}
		claims := auth.ParseAccessToken(token)
		if claims == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error parsing auth token"})
			return
		}

		var user database.User

		if err := json.NewDecoder(c.Request.Body).Decode(&user); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer c.Request.Body.Close()

		err := database.UpdateUser(client, "users", user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		response := map[string]interface{}{
			"message": "User Updated!",
		}

		c.JSON(http.StatusOK, response)

	}
}

func HandleUpdateUserPassword(client *dynamodb.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Auth header not found"})
			return
		}
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == authHeader {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "API token error"})
			return
		}
		claims := auth.ParseAccessToken(token)
		if claims == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error parsing auth token"})
			return
		}

		var user database.User

		if err := json.NewDecoder(c.Request.Body).Decode(&user); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer c.Request.Body.Close()

		hashedPassword, err := auth.HashedPassword(user.Password)
		if err != nil {
			return
		}

		user.Password = hashedPassword

		err = database.UpdatePassword(client, "users", user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		response := map[string]interface{}{
			"message": "User Password Updated!",
		}

		c.JSON(http.StatusOK, response)

	}
}

func HandleDeleteUserById(client *dynamodb.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Auth header not found"})
			return
		}
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == authHeader {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "API token error"})
			return
		}
		claims := auth.ParseAccessToken(token)
		if claims == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error parsing auth token"})
			return
		}

		err := database.DeleteUser(client, "users", id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		response := map[string]interface{}{
			"message": "User Deleted!",
		}

		c.JSON(http.StatusOK, response)

	}
}
