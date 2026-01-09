package server

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"effective-invention/server/pgdb"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func setupTestRouter(db *sql.DB) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	protected := r.Group("/api")
	addOpenRoutes(r, db)
	addProtectedRoutes(protected, db)
	return r
}

func TestUserFullLifecycle(t *testing.T) {

	// setup jwt
	t.Setenv("ACCESS_SECRET", "test_secret_for_lifecycle")
	t.Setenv("REFRESH_SECRET", "test_secret_for_lifecycle")

	// setup db connection
	connStr := "user=postgres password=postgres dbname=postgres sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	assert.NoError(t, err)
	defer db.Close()

	cleanTable := func() {
		_, err := db.Exec("TRUNCATE TABLE users RESTART IDENTITY CASCADE")
		if err != nil {
			t.Logf("Warning: Cleanup failed: %v", err)
		}
	}

	// on start to ensure a fresh state if a previous crash left data
	cleanTable()

	t.Cleanup(func() {
		cleanTable()
	})

	err = pgdb.CreateUsersTable(db)
	assert.NoError(t, err)

	router := setupTestRouter(db)

	var authToken string
	var userID string
	userEmail := "lifecycle@test.com"
	oldPass := "OldPassword123"
	newPass := "NewPassword456"
	userName := "Lifecycle User"

	// create
	t.Run("Register User", func(t *testing.T) {
		regBody := map[string]string{
			"name":     userName,
			"email":    userEmail,
			"password": oldPass,
		}
		body, _ := json.Marshal(regBody)
		req, _ := http.NewRequest("POST", "/auth/register", bytes.NewBuffer(body))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		authToken = resp["token"].(string)
		userData := resp["user"].(map[string]interface{})
		userID = userData["id"].(string)
	})

	// login
	t.Run("Login Success", func(t *testing.T) {
		loginBody := map[string]string{"email": userEmail, "password": oldPass}
		body, _ := json.Marshal(loginBody)
		req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(body))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	// get all users
	t.Run("Get All Users", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/users", nil)
		req.Header.Set("Authorization", "Bearer "+authToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var users []pgdb.User
		json.Unmarshal(w.Body.Bytes(), &users)
		assert.GreaterOrEqual(t, len(users), 1)
	})

	// update password
	t.Run("Update Password", func(t *testing.T) {
		updBody := map[string]string{
			"userId":          userID,
			"currentPassword": oldPass,
			"newPassword":     newPass,
		}
		body, _ := json.Marshal(updBody)

		req, _ := http.NewRequest("PUT", "/api/users/password", bytes.NewBuffer(body))
		req.Header.Set("Authorization", "Bearer "+authToken)
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code, "Password update failed - check your router mapping")
	})

	// get user by id
	t.Run("Get User By ID", func(t *testing.T) {
		url := fmt.Sprintf("/api/users/%s", userID)
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("Authorization", "Bearer "+authToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var user pgdb.User
		json.Unmarshal(w.Body.Bytes(), &user)
		assert.Equal(t, userEmail, user.Email)
	})

	// try old password login
	t.Run("Login Old Password Fail", func(t *testing.T) {
		loginBody := map[string]string{"email": userEmail, "password": oldPass}
		body, _ := json.Marshal(loginBody)
		req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(body))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	var refreshToken string

	// try new password login
	t.Run("Login New Password Success", func(t *testing.T) {
		loginBody := map[string]string{"email": userEmail, "password": newPass}
		body, _ := json.Marshal(loginBody)
		req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(body))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)

		authToken = resp["token"].(string)
		refreshToken = resp["refreshToken"].(string)

		assert.NotEmpty(t, authToken)
		assert.NotEmpty(t, refreshToken)
	})

	// refresh token
	t.Run("Refresh Auth Token", func(t *testing.T) {
		refreshPayload := map[string]string{
			"refresh_token": refreshToken,
		}
		body, _ := json.Marshal(refreshPayload)

		req, _ := http.NewRequest("POST", "/auth/refresh", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code, "Refresh failed: "+w.Body.String())

		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)

		newAccess, ok := resp["access_token"].(string)
		assert.True(t, ok, "access_token not found in refresh response")

		authToken = newAccess
	})

	// delete with new token
	t.Run("Delete User", func(t *testing.T) {
		url := fmt.Sprintf("/api/users/%s", userID)
		req, _ := http.NewRequest("DELETE", url, nil)

		req.Header.Set("Authorization", "Bearer "+authToken)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}
