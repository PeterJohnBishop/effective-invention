package amazonwebservices

import (
	"bytes"
	"effective-invention/server/amazonwebservices/auth"
	"effective-invention/server/amazonwebservices/database"
	"encoding/json"
	"fmt"
	"image/png"
	"net/http"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/rekognition"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt"
	qrcode "github.com/skip2/go-qrcode"
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

		err = StreamUploadFile(client, header.Filename, file)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "File uploaded successfully",
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

		fileKey := fmt.Sprintf("uploads/%s", filename)

		url, err := GeneratePresignedDownloadURL(client, fileKey)
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

func HandleFileDOwnloadLinkQR(client *s3.Client) gin.HandlerFunc {
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

		fileKey := fmt.Sprintf("uploads/%s", filename)

		url, err := GeneratePresignedDownloadURL(client, fileKey)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		qrImg, err := qrcode.New(url, qrcode.Medium)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate QR code"})
			return
		}
		var buf bytes.Buffer
		if err := png.Encode(&buf, qrImg.Image(256)); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encode QR image"})
			return
		}
		c.Header("Content-Type", "image/png")
		c.Header("Content-Disposition", "inline; filename=\"download_qr.png\"")
		c.Writer.Write(buf.Bytes())

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
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
			return
		}

		u, err := uuid.NewV1()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate ID"})
			return
		}
		userId := fmt.Sprintf("USER_%s", u.String())

		hashedPassword, err := auth.HashedPassword(user.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to secure password"})
			return
		}

		err = database.CreateUser(
			client,
			"users",
			userId,
			user.Name,
			user.Email,
			hashedPassword,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not save user to database"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "User successfully created",
			"user": gin.H{
				"id":    userId,
				"name":  user.Name,
				"email": strings.ToLower(user.Email),
			},
		})
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

func HandleFacialAnalysis(client *rekognition.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		file := c.Param("file")
		fileKey := fmt.Sprintf("uploads/%s", file)

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

		analysis, err := AnalyzeFace(client, fileKey)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		response := map[string]interface{}{
			"message":  "Analysis completed",
			"analysis": analysis,
		}

		c.JSON(http.StatusOK, response)

	}
}

func HandleFacialComparison(client *rekognition.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		sourceKey := c.Param("source")
		targetKey := c.Param("target")

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

		result, err := CompareTwoFaces(client, sourceKey, targetKey)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error comparing faces. No faces found?"})
			return
		}

		response := map[string]interface{}{
			"message":  "Comparison completed",
			"analysis": result,
		}

		c.JSON(http.StatusOK, response)

	}
}

func HandleUploadUserFile(dynamodb_client *dynamodb.Client, s3_client *s3.Client) gin.HandlerFunc {
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

		err = StreamUploadFile(s3_client, header.Filename, file)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		fileKey := fmt.Sprintf("uploads/%s", header.Filename)
		userId := claims.ID

		id, err := uuid.NewV1()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		fileId := fmt.Sprintf("FILE_%s", id)

		saveErr := database.CreateFile(dynamodb_client, "files", fileId, fileKey, userId)
		if saveErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error saving file record."})
			return
		}

		response := map[string]interface{}{
			"message": "file saved",
		}

		c.JSON(http.StatusOK, response)

	}
}

func HandleGetUserFiles(dynamodb_client *dynamodb.Client) gin.HandlerFunc {
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

		userId := claims.ID

		var userFiles []database.UserFile
		var err error
		userFiles, err = database.ListFilesByUserSorted(dynamodb_client, "files", userId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving files."})
			return
		}

		response := map[string]interface{}{
			"message":   "user files found",
			"userFiles": userFiles,
		}

		c.JSON(http.StatusOK, response)
	}
}

func HandleDeleteUserFileById(dynamodb_client *dynamodb.Client, s3_client *s3.Client) gin.HandlerFunc {
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

		var userFile *database.UserFile
		var err error

		userFile, err = database.GetFile(dynamodb_client, "files", id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		err = DeleteS3File(s3_client, userFile.FileKey)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		err = database.DeleteFile(dynamodb_client, "files", id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		response := map[string]interface{}{
			"message": "File Deleted!",
		}

		c.JSON(http.StatusOK, response)

	}
}
